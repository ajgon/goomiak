package main

import (
	"bufio"
	"fmt"
	"os"
	"z80/cpu"
	"z80/dma"
	"z80/machine"
	"z80/memory"
	"z80/video"

	"github.com/veandco/go-sdl2/sdl"
)

func drawScreen(renderer *sdl.Renderer, texture *sdl.Texture, video *video.Video) {
	pixels := video.AllPixels()
	texture.Update(nil, pixels, (256+48+48)*4)
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

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
	video := video.VideoNew(dma)
	loadFileToMemory(dma, 0x0000, "./roms/48.rom")

	cpu := cpu.CPUNew(dma, machine.Spectrum48k)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 256+48+48, 192+48+56, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.WINDOW_OPENGL|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	renderer.Clear()
	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, -1, 256+48+48, 192+48+56)

	// 14336 - 10776 =
	tstatesUla := uint64(0)
	var screenWidthPixels, screenHeightPixels uint64 = 256, 192
	var borderTopPixels, borderRightPixels, borderLeftPixels uint64 = 48, 48, 48
	initialDrawingTstate := machine.Spectrum48k.InitialContendedTstate - machine.Spectrum48k.TstatesPerScanline*borderTopPixels + borderLeftPixels/2
	//reader := bufio.NewReader(os.Stdin)
	running := true
	for running {
		for tstatesUla = 0; tstatesUla < machine.Spectrum48k.FrameLength; tstatesUla++ {
			if tstatesUla >= initialDrawingTstate { // beam returned, ULA starts drawing border
				y := (tstatesUla - initialDrawingTstate) / machine.Spectrum48k.TstatesPerScanline
				x := ((tstatesUla - initialDrawingTstate) % machine.Spectrum48k.TstatesPerScanline) * 2 // every tstate is 2 pixels
				if x < borderLeftPixels+screenWidthPixels+borderRightPixels {                           // for x >= 352 it effectively means beam return
					//fmt.Printf("y = %d, x = [%d, %d] (ulaT = %d, cpuT = %d)", y, x, x+1, tstatesUla, cpu.Tstates())
					video.PaintPixel(x, y)
					video.PaintPixel(x+1, y)
					if y < borderTopPixels || y >= borderTopPixels+screenHeightPixels || x < borderLeftPixels || x >= borderLeftPixels+screenWidthPixels {
						//fmt.Printf(" [BORDER]")
					}
					//fmt.Printf("\n")
				}
			}

			if cpu.Tstates()%machine.Spectrum48k.FrameLength <= tstatesUla {
				//cpu.DebugStep()
				cpu.Step()
				//if tstatesUla > 3500 {
				//}
			}
		}
		drawScreen(renderer, texture, video)
		//reader.ReadString('\n')

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				for i := 0x4000; i < 0x4000+6912; i++ {
					a, _ := dma.GetMemoryByte(uint16(i))
					fmt.Printf("%02x", a)
				}
				running = false
				break
			}
		}
	}
	//for {
	//if c.tstates > 3584 { // beam returned, ULA starts drawing border

	//}
	//cpu.DebugStep()
	//}
}
