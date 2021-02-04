package machine

import (
	"bufio"
	"os"
	"time"
	"z80/bus"
	"z80/cpu"
	"z80/dma"
	"z80/loader"
	"z80/memory"
	"z80/video"
)

type MachineConfig struct {
	FrameLength            uint
	InitialContendedTstate uint
	TstatesPerScanline     uint
	ContentionPattern      [8]uint8
}

type Machine struct {
	Config           MachineConfig
	ContentionDelays []uint8

	IO          *bus.IO // @todo temporary
	CPU         *cpu.CPU
	DMA         *dma.DMA
	ULA         *video.ULA
	VideoDriver video.VideoDriver

	fullSpeed bool
}

func (m *Machine) buildContentionPattern() {
	m.ContentionDelays = make([]uint8, m.Config.FrameLength)
	for line := uint(0); line < 192; line++ {
		lineFirstTstate := m.Config.InitialContendedTstate + line*m.Config.TstatesPerScanline
		for x := uint(0); x < 128; x++ {
			m.ContentionDelays[lineFirstTstate+x] = m.Config.ContentionPattern[x%8]
		}
	}
}

func (m *Machine) build() {
	m.buildContentionPattern()

	memory := memory.NewMemory()
	m.IO = bus.NewIO()
	videoMemoryHandler := video.VideoMemoryHandlerNew()
	m.DMA = dma.NewDMA(memory, videoMemoryHandler)

	pixelRenderer := video.NewPixelRenderer(m.DMA)
	m.VideoDriver = video.NewSDLVideoDriver(pixelRenderer)
	m.ULA = video.NewULA(
		m.IO,
		pixelRenderer,
		video.ULAConfig{
			InitialContendedTstate: m.Config.InitialContendedTstate,
			TstatesPerScanline:     m.Config.TstatesPerScanline,
		},
	)

	m.CPU = cpu.NewCPU(m.IO, m.DMA, cpu.CPUConfig{ContentionDelays: m.ContentionDelays, FrameLength: m.Config.FrameLength})
}

func (m *Machine) FullSpeed(value bool) {
	m.fullSpeed = value
}

func (m *Machine) Run() {
	running := true
	frames := 0
	for running {
		startTime := time.Now()
		for m.ULA.Tstates < m.Config.FrameLength {
			if m.ULA.Tstates == 32 {
				m.CPU.SetIRQ(false)
			}

			m.ULA.Step()
			if m.CPU.HandleInterrupt() {
				continue
			}

			if m.CPU.Tstates%m.Config.FrameLength <= m.ULA.Tstates {
				m.CPU.Step()
			}
		}

		m.ULA.Tstates = 0
		m.ULA.Flash = (frames/32)%2 == 1
		frames++
		m.VideoDriver.DrawScreen() // @todo this goes to ULA as it needs to handle SDL events as well
		keyPressedMasks := m.VideoDriver.KeyPressedOut()
		for kpAddr, kpValue := range keyPressedMasks {
			// @todo this goes to ULA
			m.IO.Write(1, (uint16(kpAddr)<<8)|0xfe, kpValue)
		}
		m.CPU.SetIRQ(true)
		passedTime := time.Since(startTime)

		if !m.fullSpeed && passedTime < 20*time.Millisecond {
			time.Sleep(20*time.Millisecond - passedTime)
		}
	}
}

func (m *Machine) LoadSnapshot(snapshot loader.Snapshot) {
	m.DMA.LoadSnapshot(snapshot)
	m.CPU.LoadSnapshot(snapshot)
}

func (m *Machine) InsertTape(tape *loader.TapFile) {
	m.CPU.InsertTape(tape)
}

func NewMachine(config MachineConfig) *Machine {
	machine := &Machine{Config: config}
	machine.build()

	return machine
}

func (m *Machine) LoadDataToMemory(address uint16, bytes []byte) {
	m.DMA.LoadData(address, bytes)
}

func (m *Machine) LoadFileToMemory(address uint16, filePath string) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic("error loading file!")
	}

	stat, _ := file.Stat()

	bytes := make([]byte, stat.Size())
	buf := bufio.NewReader(file)
	buf.Read(bytes)

	m.DMA.LoadData(address, bytes)
}
