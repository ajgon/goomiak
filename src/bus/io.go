package bus

const sourceCPU uint8 = 0
const sourceULA uint8 = 1

type IO struct {
	border      uint8
	keysPressed [256]uint8
}

// source - who wants to read
func (io *IO) Read(source uint8, address uint16) uint8 {
	if address&0x0001 == 0x0000 {
		// ULA reads BUS, wants to get border
		if source == sourceULA {
			return io.border
		}

		// CPU reads BUS, wants keyboard code
		// @todo include EAR
		if source == sourceCPU {
			return io.keysPressed[address>>8]
		}
	}

	return 0xff
}

// source - who writes the bus
func (io *IO) Write(source uint8, address uint16, value uint8) {
	if address&0x0001 == 0x0000 {
		if source == sourceCPU {
			// CPU writes BUS, wants to set border
			io.border = value
			return
		}

		if source == sourceULA {
			// ULA writes BUS, wants to set keystroke
			io.keysPressed[uint8(address>>8)] = value & 0xbf
		}
	}
}

func NewIO() *IO {
	io := &IO{}

	for i := 0; i < 0xff; i++ {
		io.keysPressed[i] = 0xbf
	}
	io.keysPressed[0xff] = 0xbf
	return io
}
