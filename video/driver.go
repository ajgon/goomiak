package video

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type VideoDriver interface {
	DrawScreen()
	Close()
}

type SDLVideoDriver struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Texture  *sdl.Texture

	PixelRenderer PixelRenderer
}

func (svd *SDLVideoDriver) Close() {
	svd.Window.Destroy()
	sdl.Quit()
}

func (svd *SDLVideoDriver) DrawScreen() {
	pixels := svd.PixelRenderer.Pixels()

	svd.Texture.Update(nil, pixels, int(fullWidth*4))
	svd.Renderer.Copy(svd.Texture, nil, nil)
	svd.Renderer.Present()
}

func NewSDLVideoDriver(pixelRenderer *PixelRenderer) *SDLVideoDriver {
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		"test",
		//sdl.WINDOWPOS_UNDEFINED,
		//sdl.WINDOWPOS_UNDEFINED,
		10, 10,
		int32(fullWidth),
		int32(fullHeight),
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.WINDOW_OPENGL|sdl.WINDOW_ALLOW_HIGHDPI,
	)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	renderer.Clear()

	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		-1,
		int32(fullWidth),
		int32(fullHeight),
	)

	renderer.Copy(texture, nil, nil)
	renderer.Present()
	fmt.Println("SLEEP")
	time.Sleep(5 * time.Second)
	fmt.Println("/SLEEP")

	return &SDLVideoDriver{
		Window:   window,
		Renderer: renderer,
		Texture:  texture,
	}
}
