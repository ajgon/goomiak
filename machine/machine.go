package machine

type Machine struct {
	FrameLength            uint64
	InitialContendedTstate uint64
	TstatesPerScanline     uint64
	ContentionPattern      [8]uint8
	ContentionDelays       []uint8
}

func (m *Machine) buildContentionPattern() {
	m.ContentionDelays = make([]uint8, m.FrameLength)
	for line := uint64(0); line < 192; line++ {
		lineFirstTstate := m.InitialContendedTstate + line*224
		for x := uint64(0); x < 128; x++ {
			m.ContentionDelays[lineFirstTstate+x] = m.ContentionPattern[x%8]
		}
	}
}
