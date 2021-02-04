package video

import "z80/bus"

const busSource uint8 = 1

type ULAConfig struct {
	InitialContendedTstate uint
	TstatesPerScanline     uint
}

type ULA struct {
	Config  ULAConfig
	Tstates uint
	Flash   bool

	PixelRenderer *PixelRenderer

	initialDrawingTstate uint
	io                   *bus.IO
}

func (u *ULA) Step() {
	u.Tstates++

	if u.Tstates < u.initialDrawingTstate {
		// beam returns, ULA has nothing to do
		return
	}

	border := u.io.Read(busSource, 0x00fe)
	u.PixelRenderer.SetBorder(border)

	tstateRef := u.Tstates - u.initialDrawingTstate

	y := tstateRef / u.Config.TstatesPerScanline
	x := (tstateRef % u.Config.TstatesPerScanline) * 2 // every tstate means two pixels

	if x >= fullWidth { // after that beams returns, no drawing is necessary
		return
	}

	u.PixelRenderer.PaintPixel(x, y, u.Flash)
	u.PixelRenderer.PaintPixel(x+1, y, u.Flash)
}

func NewULA(io *bus.IO, pixelRenderer *PixelRenderer, config ULAConfig) *ULA {
	ula := &ULA{io: io, PixelRenderer: pixelRenderer, Config: config}

	ula.initialDrawingTstate = config.InitialContendedTstate - config.TstatesPerScanline*borderTopHeight + borderLeftWidth/2

	return ula
}
