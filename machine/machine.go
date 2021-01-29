package machine

import (
	"z80/cpu"
	"z80/dma"
	"z80/memory"
	"z80/video"
)

type MachineConfig struct {
	FrameLength            uint64
	InitialContendedTstate uint64
	TstatesPerScanline     uint64
	ContentionPattern      [8]uint8
}

type Machine struct {
	Config           MachineConfig
	ContentionDelays []uint8

	CPU   *cpu.CPU
	DMA   *dma.DMA
	Video video.VideoDriver
}

func (m *Machine) buildContentionPattern() {
	m.ContentionDelays = make([]uint8, m.Config.FrameLength)
	for line := uint64(0); line < 192; line++ {
		lineFirstTstate := m.Config.InitialContendedTstate + line*m.Config.TstatesPerScanline
		for x := uint64(0); x < 128; x++ {
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
	m.Video = video.NewSDLVideoDriver(pixelRenderer)

	m.CPU = cpu.NewCPU(m.DMA, cpu.CPUConfig{ContentionDelays: m.ContentionDelays, FrameLength: m.Config.FrameLength})
}

func (m *Machine) Run() {
	panic("RUN")
}

func NewMachine(config MachineConfig) *Machine {
	machine := &Machine{Config: config}
	machine.build()

	return machine
}
