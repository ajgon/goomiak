package machine

var Spectrum48k Machine = Machine{
	FrameLength:            69888,
	InitialContendedTstate: 14335,
	TstatesPerScanline:     224,
	ContentionPattern:      [8]uint8{6, 5, 4, 3, 2, 1, 0, 0},
}

func init() {
	Spectrum48k.buildContentionPattern()
}
