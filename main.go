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

	stat, _ := file.Stat()

	bytes := make([]byte, stat.Size())
	buf := bufio.NewReader(file)
	buf.Read(bytes)

	dma.SetMemoryBulk(address, bytes)
}

func main() {
	mem := memory.MemoryNew()
	videoMemoryHandler := video.VideoMemoryHandlerNew()
	dma := dma.DMANew(mem, videoMemoryHandler)
	//video := video.VideoNew(dma)
	//loadFileToMemory(dma, 0x0000, "./roms/48.rom")
	//loadFileToMemory(dma, 0x8000, "./roms/zexdoc.rom")
	loadFileToMemory(dma, 0x0100, "./roms/zexall.cpm")

	cpu := cpu.CPUNew(dma)
	cpu.PC = 0x0100
	cpu.SP = 0x0000
	dma.SetMemoryByte(0x05, 0xc9) // RET

	for {
		cpu.Step()
		if cpu.PC == 0 {
			break
		}

		if cpu.PC == 5 {
			if uint8(cpu.BC) == 2 {
				fmt.Printf("%c", uint8(cpu.DE))
			}

			if uint8(cpu.BC) == 9 {
				i := cpu.DE
				for {
					char := dma.GetMemory(i)
					if char != 36 {
						fmt.Printf("%c", char)
					} else {
						break
					}
					i++
				}
			}
		}
	}
}
