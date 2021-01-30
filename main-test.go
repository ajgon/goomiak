package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
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
	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, -1, 256, 192)
	running := true
	for running {
		texture.Update(nil, make([]byte, 256*192*4), 256*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

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
