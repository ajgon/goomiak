package machine

func NewSpectrum48k() *Machine {
	machineConfig := MachineConfig{
		FrameLength:            69888,
		InitialContendedTstate: 14335,
		TstatesPerScanline:     224,
		ContentionPattern:      [8]uint8{6, 5, 4, 3, 2, 1, 0, 0},
	}

	return NewMachine(machineConfig)
}
