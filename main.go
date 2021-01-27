package main

import (
	"bufio"
	"os"
	"z80/cpu"
	"z80/dma"
	"z80/machine"
	"z80/memory"
	"z80/video"
)

func loadFileToMemory(dma *dma.DMA, address uint16, filePath string) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic("error loading file!")
	}

	stat, _ := file.Stat()

	bytes := make([]byte, stat.Size())
	buf := bufio.NewReader(file)
	buf.Read(bytes)

	dma.LoadData(address, bytes)
}

func main() {
	mem := memory.NewWritableMemory()
	videoMemoryHandler := video.VideoMemoryHandlerNew()
	dma := dma.DMANew(mem, videoMemoryHandler)
	//video := video.VideoNew(dma)
	loadFileToMemory(dma, 0x0000, "./roms/48.rom")

	cpu := cpu.CPUNew(dma, machine.Spectrum48k)
	//reader := bufio.NewReader(os.Stdin)

	for {
		cpu.DebugStep()
		//reader.ReadString('\n')
	}
}
