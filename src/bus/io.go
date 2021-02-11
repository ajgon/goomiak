package bus

const sourceCPU uint8 = 0
const sourceULA uint8 = 1

type IO struct {
	border      uint8
	keysPressed map[uint16]uint8
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
			addressLeft := (address & 0x0f00) | 0xf0fe
			addressRight := (address & 0xf000) | 0x0ffe

			valueLeft := io.keysPressed[addressLeft]
			valueRight := io.keysPressed[addressRight]

			value := valueLeft & valueRight
			// if no key is pressed, then return 0xbf as ULA doesn't expose anything
			// (assuming EAR is not used - @todo)
			if value < 0x20 {
				return valueLeft & valueRight
			}

			return 0xbf
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
			io.keysPressed[address] = value & 0xbf
		}
	}
}

func NewIO() *IO {
	io := &IO{}
	io.keysPressed = make(map[uint16]uint8)
	for i := uint16(0); i <= 0xfffe; i++ {
		io.keysPressed[i] = 0xff
	}
	io.keysPressed[0xffff] = 0xff

	return io
}
