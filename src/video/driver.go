package video

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type VideoDriver interface {
	DrawScreen()
	KeyPressedOut() map[uint8]uint8
	Close()
}

type SDLVideoDriver struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Texture  *sdl.Texture

	PixelRenderer   *PixelRenderer
	keyMap          map[sdl.Keycode]uint16
	keyPressedMasks map[uint8]uint8
}

func (svd *SDLVideoDriver) Close() {
	svd.Window.Destroy()
	sdl.Quit()
}

func (svd *SDLVideoDriver) KeyPressedOut() map[uint8]uint8 {
	return svd.keyPressedMasks
}

func (svd *SDLVideoDriver) DrawScreen() {
	pixels := svd.PixelRenderer.Pixels()

	for _, mask := range [8]uint8{0xfe, 0xfd, 0xfb, 0xf7, 0xef, 0xdf, 0xbf, 0x7f} {
		svd.keyPressedMasks[mask] = 0x1f
	}

	svd.Texture.Update(nil, pixels, int(fullWidth*4))
	svd.Renderer.Copy(svd.Texture, nil, nil)
	svd.Renderer.Present()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			os.Exit(0)
			break
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym

			if keyCode < 10000 {
				if t.Repeat > 0 || t.State == sdl.PRESSED {
					mask := svd.keyMap[keyCode]
					if mask != 0 {
						svd.keyPressedMasks[uint8(mask>>8)] = uint8(mask)
					}
				}
			}

			if t.Keysym.Mod&sdl.KMOD_LCTRL == sdl.KMOD_LCTRL {
				// Symbol
				mask := svd.keyMap[1]
				svd.keyPressedMasks[uint8(mask>>8)] = svd.keyPressedMasks[uint8(mask>>8)] & uint8(mask)
			}

			if t.Keysym.Mod&sdl.KMOD_LSHIFT == sdl.KMOD_LSHIFT {
				// Caps-shift
				mask := svd.keyMap[0]
				svd.keyPressedMasks[uint8(mask>>8)] = svd.keyPressedMasks[uint8(mask>>8)] & uint8(mask)
			}
		}
	}

}

func NewSDLVideoDriver(pixelRenderer *PixelRenderer) *SDLVideoDriver {
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		"GOomiak",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
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

	// keyCode => 0xRRDD (rr - row address (higher byte), dd column mask)
	// https://sinclair.wiki.zxnet.co.uk/wiki/ZX_Spectrum_ULA#Keyboard_Half-rows
	keyMap := map[sdl.Keycode]uint16{
		//     1         2           3           4           5           6           7           8           9           0
		49: 0xf71e, 50: 0xf71d, 51: 0xf71b, 52: 0xf717, 53: 0xf70f, 54: 0xef0f, 55: 0xef17, 56: 0xef1b, 57: 0xef1d, 48: 0xef1e,
		//      Q          W            E            R            T            Y           U             I            O            P
		113: 0xfb1e, 119: 0xfb1d, 101: 0xfb1b, 114: 0xfb17, 116: 0xfb0f, 121: 0xdf0f, 117: 0xdf17, 105: 0xdf1b, 111: 0xdf1d, 112: 0xdf1e,
		//     A          S            D            F            G            H            J            K            L          ENTER
		97: 0xfd1e, 115: 0xfd1d, 100: 0xfd1b, 102: 0xfd17, 103: 0xfd0f, 104: 0xbf0f, 106: 0xbf17, 107: 0xbf1b, 108: 0xbf1d, 13: 0xbf1e,
		//  CAPS         Z            X           C            V           B            N            M         symbol       SPACE
		0: 0xfe1e, 122: 0xfe1d, 120: 0xfe1b, 99: 0xfe17, 118: 0xfe0f, 98: 0x7f0f, 110: 0x7f17, 109: 0x7f1b, 1: 0x7f1d, 32: 0x7f1e,
	}

	return &SDLVideoDriver{
		Window:          window,
		Renderer:        renderer,
		Texture:         texture,
		PixelRenderer:   pixelRenderer,
		keyMap:          keyMap,
		keyPressedMasks: make(map[uint8]uint8),
	}
}
