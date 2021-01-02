package memory

import (
	"fmt"
)

type Memory struct {
	bytes [65536]uint8
}

func (m *Memory) Get(address uint16) uint8 {
	return m.bytes[address]
}

func (m *Memory) SetByte(address uint16, value uint8) {
	m.bytes[address] = value
}

func (m *Memory) SetBulk(address uint16, bytes []uint8) {
	if len(bytes) == 0 {
		return
	}

	for index, value := range bytes {
		m.bytes[address+uint16(index)] = value
	}
}

func (m *Memory) Dump(addressRange ...uint16) {
	var from, to uint16

	if len(addressRange) != 2 {
		from = 0
		to = 65535
	} else {
		from = addressRange[0]
		to = addressRange[1]
	}

	for i := from / 16; i < to/16; i++ {
		fmt.Printf("%04x  ", i*16)
		for b := uint16(0); b < 16; b++ {
			fmt.Printf("%02x", m.bytes[i*16+b])
			if b == 7 {
				fmt.Printf("  ")
			} else if b != 15 {
				fmt.Printf(" ")
			}
		}
		fmt.Print("\n")
	}
}

func (m *Memory) Clear() {
	m.bytes = [65536]uint8{}
}

func MemoryNew() *Memory {
	return new(Memory)
}
