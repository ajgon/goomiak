package video

type ULAConfig struct {
	InitialContendedTstate uint
	TstatesPerScanline     uint
}

type ULA struct {
	Config  ULAConfig
	Tstates uint

	PixelRenderer        *PixelRenderer
	initialDrawingTstate uint
}

func (u *ULA) Step() {
	u.Tstates++

	if u.Tstates < u.initialDrawingTstate {
		// beam returns, ULA has nothing to do
		return
	}

	tstateRef := u.Tstates - u.initialDrawingTstate

	y := tstateRef / u.Config.TstatesPerScanline
	x := (tstateRef / u.Config.TstatesPerScanline) * 2 // every tstate means two pixels

	if x > fullWidth { // after that beams returns, no drawing is necessary
		return
	}

	u.PixelRenderer.PaintPixel(x, y)
	u.PixelRenderer.PaintPixel(x+1, y)
}

func NewULA(pixelRenderer *PixelRenderer, config ULAConfig) *ULA {
	ula := &ULA{Config: config}
	ula.initialDrawingTstate = config.InitialContendedTstate - config.TstatesPerScanline*borderTopHeight + borderLeftWidth/2
	ula.PixelRenderer = pixelRenderer

	return ula
}
