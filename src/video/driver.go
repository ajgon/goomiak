package video

import (
	"math/bits"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type VideoDriver interface {
	DrawScreen()
	KeyPressedOut() [256]uint8
	Close()
}

type SDLVideoDriver struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Texture  *sdl.Texture

	PixelRenderer   *PixelRenderer
	keyMap          map[sdl.Keycode]uint16
	keyPressedMasks [256]uint8
}

func (svd *SDLVideoDriver) Close() {
	svd.Window.Destroy()
	sdl.Quit()
}

func (svd *SDLVideoDriver) KeyPressedOut() [256]uint8 {
	return svd.keyPressedMasks
}

func (svd *SDLVideoDriver) DrawScreen() {
	pixels := svd.PixelRenderer.Pixels()

	svd.keyPressedMasks[0xff] = 0xbd
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
			mergedValue := uint8(0xbf)
			for i := 0; i <= 0xff; i++ {
				svd.keyPressedMasks[i] = 0xbf
			}

			if keyCode < 10000 {
				if t.Repeat > 0 || t.State == sdl.PRESSED {
					mask := svd.keyMap[keyCode]
					if mask != 0 {
						mergedValue = uint8(mask)
						ulaAddess := uint8(mask >> 8)
						if ulaAddess < 0xf0 {
							base := ulaAddess & 0xf0
							svd.keyPressedMasks[base+0x07] = mergedValue
							svd.keyPressedMasks[base+0x0b] = mergedValue
							svd.keyPressedMasks[base+0x0d] = mergedValue
							svd.keyPressedMasks[base+0x0e] = mergedValue
							svd.keyPressedMasks[base+0x0f] = mergedValue
						} else {
							base := ulaAddess & 0x0f
							svd.keyPressedMasks[base+0x70] = mergedValue
							svd.keyPressedMasks[base+0xb0] = mergedValue
							svd.keyPressedMasks[base+0xd0] = mergedValue
							svd.keyPressedMasks[base+0xe0] = mergedValue
							svd.keyPressedMasks[base+0xf0] = mergedValue
						}
					}
				}
			}

			if t.Keysym.Mod&sdl.KMOD_LCTRL == sdl.KMOD_LCTRL {
				// Symbol
				mask := svd.keyMap[1]
				mergedValue = mergedValue & uint8(mask)
				ulaAddess := uint8(mask>>8) & 0xf0
				svd.keyPressedMasks[ulaAddess+0x07] = svd.keyPressedMasks[ulaAddess+0x07] & uint8(mask)
				svd.keyPressedMasks[ulaAddess+0x0b] = svd.keyPressedMasks[ulaAddess+0x0b] & uint8(mask)
				svd.keyPressedMasks[ulaAddess+0x0d] = svd.keyPressedMasks[ulaAddess+0x0d] & uint8(mask)
				svd.keyPressedMasks[ulaAddess+0x0e] = svd.keyPressedMasks[ulaAddess+0x0e] & uint8(mask)
				svd.keyPressedMasks[ulaAddess+0x0f] = svd.keyPressedMasks[ulaAddess+0x0f] & uint8(mask)
			}

			if t.Keysym.Mod&sdl.KMOD_LSHIFT == sdl.KMOD_LSHIFT {
				// Caps-shift
				mask := svd.keyMap[0]
				mergedValue = mergedValue & uint8(mask)
				ulaAddess := uint8(mask>>8) & 0x0f
				svd.keyPressedMasks[ulaAddess+0x70] = svd.keyPressedMasks[ulaAddess+0x70] & uint8(mask)
				svd.keyPressedMasks[ulaAddess+0xb0] = svd.keyPressedMasks[ulaAddess+0xb0] & uint8(mask)
				svd.keyPressedMasks[ulaAddess+0xd0] = svd.keyPressedMasks[ulaAddess+0xd0] & uint8(mask)
				svd.keyPressedMasks[ulaAddess+0xe0] = svd.keyPressedMasks[ulaAddess+0xe0] & uint8(mask)
				svd.keyPressedMasks[ulaAddess+0xf0] = svd.keyPressedMasks[ulaAddess+0xf0] & uint8(mask)
			}

			for i := uint8(0); i < 0xff; i++ {
				if bits.OnesCount8(i) >= 6 {
					continue
				}

				if svd.keyPressedMasks[i] == 0xbf {
					svd.keyPressedMasks[i] = mergedValue
				}
			}

			if keyCode == 1073741886 {
				svd.keyPressedMasks[0xff] = 0x33
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
		49: 0xf7be, 50: 0xf7bd, 51: 0xf7bb, 52: 0xf7b7, 53: 0xf7af, 54: 0xefaf, 55: 0xefb7, 56: 0xefbb, 57: 0xefbd, 48: 0xefbe,
		//      Q          W            E            R            T            Y           U             I            O            P
		113: 0xfbbe, 119: 0xfbbd, 101: 0xfbbb, 114: 0xfbb7, 116: 0xfbaf, 121: 0xdfaf, 117: 0xdfb7, 105: 0xdfbb, 111: 0xdfbd, 112: 0xdfbe,
		//     A          S            D            F            G            H            J            K            L          ENTER
		97: 0xfdbe, 115: 0xfdbd, 100: 0xfdbb, 102: 0xfdb7, 103: 0xfdaf, 104: 0xbfaf, 106: 0xbfb7, 107: 0xbfbb, 108: 0xbfbd, 13: 0xbfbe,
		//  CAPS         Z            X           C            V           B            N            M         symbol       SPACE
		0: 0xfebe, 122: 0xfebd, 120: 0xfebb, 99: 0xfeb7, 118: 0xfeaf, 98: 0x7faf, 110: 0x7fb7, 109: 0x7fbb, 1: 0x7fbd, 32: 0x7fbe,
	}

	keyPressedMasks := [256]uint8{}
	for i := 0; i < 256; i++ {
		keyPressedMasks[i] = 0xbf
	}

	return &SDLVideoDriver{
		Window:          window,
		Renderer:        renderer,
		Texture:         texture,
		PixelRenderer:   pixelRenderer,
		keyMap:          keyMap,
		keyPressedMasks: keyPressedMasks,
	}
}
