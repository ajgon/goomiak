package machine

import (
	"bufio"
	"os"
	"z80/cpu"
	"z80/dma"
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

	CPU         *cpu.CPU
	DMA         *dma.DMA
	ULA         *video.ULA
	VideoDriver video.VideoDriver
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
	videoMemoryHandler := video.VideoMemoryHandlerNew()
	m.DMA = dma.NewDMA(memory, videoMemoryHandler)

	pixelRenderer := video.NewPixelRenderer(m.DMA)
	m.VideoDriver = video.NewSDLVideoDriver(pixelRenderer)
	m.ULA = video.NewULA(
		pixelRenderer,
		video.ULAConfig{
			InitialContendedTstate: m.Config.InitialContendedTstate,
			TstatesPerScanline:     m.Config.TstatesPerScanline,
		},
	)

	m.CPU = cpu.NewCPU(m.DMA, cpu.CPUConfig{ContentionDelays: m.ContentionDelays, FrameLength: m.Config.FrameLength})
}

func (m *Machine) Run() {
	running := true
	for running {
		for m.ULA.Tstates < m.Config.FrameLength {
			if m.CPU.States.IRQ && m.CPU.States.IFF1 {
				m.CPU.HandleInterrupt()
			}

			m.ULA.Step()

			if m.CPU.Tstates%m.Config.FrameLength <= m.ULA.Tstates {
				m.CPU.Step()
			}
		}

		m.ULA.Tstates = 0
		m.VideoDriver.DrawScreen() // @todo this goes to ULA as it needs to handle SDL events as well
		m.CPU.States.IRQ = true
	}
}

func NewMachine(config MachineConfig) *Machine {
	machine := &Machine{Config: config}
	machine.build()

	return machine
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
