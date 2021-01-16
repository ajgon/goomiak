package main

import (
	"bufio"
	"fmt"
	"os"
	"z80/cpu"
	"z80/dma"
	"z80/memory"
	"z80/video"
)

func loadFileToMemory(dma *dma.DMA, address uint16, filePath string) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic("error loading file!")
	}

	bytes := make([]byte, 6912)
	buf := bufio.NewReader(file)
	buf.Read(bytes)

	dma.SetMemoryBulk(address, bytes)
}

func main() {
	mem := memory.MemoryNew()
	videoMemoryHandler := video.VideoMemoryHandlerNew()
	dma := dma.DMANew(mem, videoMemoryHandler)
	//video := video.VideoNew(dma)
	//loadFileToMemory(dma, 0x8000, "./roms/zexdoc.rom")
	loadFileToMemory(dma, 0x0000, "./roms/48.rom")

	cpu := cpu.CPUNew(dma)
	cpu.PC = 0x0000
	tstates := uint64(0)
	reader := bufio.NewReader(os.Stdin)

	for {
		//for i := 0; i < 32; i++ {
		fmt.Printf("T: %d => ", tstates)
		tstates += uint64(cpu.Step())
		reader.ReadString('\n')
	}
}
