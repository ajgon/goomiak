package main

import (
	"bufio"
	"fmt"
	"os"
	"z80/cpu"
	"z80/dma"
	"z80/memory"
	"z80/video"

	"github.com/veandco/go-sdl2/sdl"
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

func drawScreen(renderer *sdl.Renderer, texture *sdl.Texture, video *video.Video) {
	if video.NeedsRefresh() {
		pixels := video.Pixels()
		texture.Update(nil, pixels, 256*4)
	}
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

func main() {
	mem := memory.MemoryNew()
	videoMemoryHandler := video.VideoMemoryHandlerNew()

	dma := dma.DMANew(mem, videoMemoryHandler)
	video := video.VideoNew(dma)
	//loadFileToMemory(dma, 0x4000, "./video/example.scr")
	loadFileToMemory(dma, 0x0000, "./roms/48.rom")
	//loadFileToMemory(dma, 0x8000, "./roms/zexdoc.rom")
	cpu := cpu.CPUNew(dma)
	cpu.PC = 0x0000
	tstates := uint64(0)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 256, 192, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.WINDOW_OPENGL|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	renderer.Clear()

	frames := uint64(0)
	start := sdl.GetPerformanceCounter()

	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, -1, 256, 192)

	running := true
	for running {
		for {
			//fmt.Printf("T: %d => ", tstates)
			tstates += uint64(cpu.Step())
			if tstates > 70908 {
				tstates = 0
				break
			}
		}
		drawScreen(renderer, texture, video)
		frames += 1
		end := sdl.GetPerformanceCounter()
		freq := sdl.GetPerformanceFrequency()
		seconds := float64(end-start) / float64(freq)

		if seconds > 1 {
			fmt.Printf("%d frames in %.1f seconds = %.1f FPS (%.3f ms/frame)\n", frames, seconds, float64(frames)/seconds, (seconds*1000)/float64(frames))
			start = end
			frames = 0
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}

}
