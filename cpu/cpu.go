package cpu

import (
	"z80/dma"
)

type CPU struct {
	PC    uint16
	AF    uint16
	AF_   uint16
	BC    uint16
	DE    uint16
	HL    uint16
	Flags uint8 // [  S  ][  Z  ][     ][  H  ][     ][ P/V ][  N  ][  C  ]

	dma *dma.DMA
}

func (c *CPU) readByte(address uint16) uint8 {
	return c.dma.GetMemory(address)
}

func (c *CPU) readWord(address uint16) uint16 {
	return uint16(c.dma.GetMemory(address+1))<<8 | uint16(c.dma.GetMemory(address))
}

func (c *CPU) increaseRegister(name rune) uint8 {
	var register uint8

	switch name {
	case 'B':
		c.BC += 256
		register = uint8(c.BC >> 8)
	case 'C':
		register = uint8(c.BC) + 1
		c.BC = (c.BC & 0xff00) | uint16(register)
	case 'D':
		c.DE += 256
		register = uint8(c.DE >> 8)
	}

	// C (carry) is not set
	// N (sub/add) flag
	c.Flags = c.Flags & 0b11111101
	// P/V (overflow) flag
	if register == 0x80 {
		c.Flags = c.Flags | 0b00000100
	} else {
		c.Flags = c.Flags & 0b11111011
	}
	// H (half carry) flag
	if register&0b00001111 == 0 {
		c.Flags = c.Flags | 0b00010000
	} else {
		c.Flags = c.Flags & 0b11101111
	}
	// Z (zero) flag
	if register == 0 {
		c.Flags = c.Flags | 0b01000000
	} else {
		c.Flags = c.Flags & 0b10111111
	}
	// S (sign) flag
	if register > 127 {
		c.Flags = c.Flags | 0b10000000
	} else {
		c.Flags = c.Flags & 0b01111111
	}
	c.PC++

	return 4
}

func (c *CPU) decreaseRegister(name rune) uint8 {
	var register uint8

	switch name {
	case 'B':
		c.BC -= 256
		register = uint8(c.BC >> 8)
	case 'C':
		register = uint8(c.BC) - 1
		c.BC = (c.BC & 0xff00) | uint16(register)
	case 'D':
		c.DE -= 256
		register = uint8(c.DE >> 8)
	}

	// C (carry) is not set
	// N (sub/add flag)
	c.Flags = c.Flags | 0b00000010
	// P/V (overflow) flag
	if register == 0x7f {
		c.Flags = c.Flags | 0b00000100
	} else {
		c.Flags = c.Flags & 0b11111011
	}
	// H (half carry) flag
	if register&0b00001111 == 15 {
		c.Flags = c.Flags | 0b00010000
	} else {
		c.Flags = c.Flags & 0b11101111
	}
	// Z (zero) flag
	if register == 0 {
		c.Flags = c.Flags | 0b01000000
	} else {
		c.Flags = c.Flags & 0b10111111
	}
	// S (sign) flag
	if register > 127 {
		c.Flags = c.Flags | 0b10000000
	} else {
		c.Flags = c.Flags & 0b01111111
	}

	c.PC++
	return 4
}

func (c *CPU) nop() uint8 {
	c.PC++

	return 4
}

func (c *CPU) ldBcXx() uint8 {
	c.BC = c.readWord(c.PC + 1)
	c.PC += 3

	return 10
}

func (c *CPU) ldBcA() uint8 {
	c.dma.SetMemoryByte(c.BC, uint8(c.AF>>8))
	c.PC++
	return 7
}

func (c *CPU) incBc() uint8 {
	c.BC++
	c.PC++
	return 6
}

func (c *CPU) incB() uint8 {
	return c.increaseRegister('B')
}

func (c *CPU) decB() uint8 {
	return c.decreaseRegister('B')
}

func (c *CPU) ldBX() uint8 {
	c.PC++
	c.BC = (c.BC & 0x00ff) | (uint16(c.readByte(c.PC)) << 8)
	c.PC++

	return 7
}

func (c *CPU) rlca() uint8 {
	a := uint8(c.AF >> 8)
	signed := a&128 == 128
	a = a << 1
	c.AF = (c.AF & 0x00ff) | (uint16(a) << 8)
	c.PC++

	// C (carry) flag
	if signed {
		c.Flags = c.Flags | 0b00000001
		c.AF = c.AF | 0x0100
	} else {
		c.Flags = c.Flags & 0b11111110
	}
	// H (half carry) flag
	c.Flags = c.Flags & 0b11101111
	// N (sub/add) flag
	c.Flags = c.Flags & 0b11111101

	return 4
}

func (c *CPU) exAfAf_() uint8 {
	c.AF, c.AF_ = c.AF_, c.AF
	c.PC++

	return 4
}

// left stores the result
func (c *CPU) addRegisters(left, right *uint16) uint8 {
	sum := *left + *right

	// C (carry) flag
	if sum < *left || sum < *right {
		c.Flags = c.Flags | 0b00000001
	} else {
		c.Flags = c.Flags & 0b11111110
	}

	// N (sub/add) flag
	c.Flags = c.Flags & 0b11111101

	// H (half carry) flag
	h := (*left ^ *right ^ sum) & 0x1000

	if h == 0x1000 {
		c.Flags = c.Flags | 0b00010000
	} else {
		c.Flags = c.Flags & 0b11101111
	}

	*left = sum
	c.PC++
	return 11

}

func (c *CPU) addHlBc() uint8 {
	return c.addRegisters(&c.HL, &c.BC)
}

func (c *CPU) ldABc() uint8 {
	value := c.dma.GetMemory(c.BC)
	c.AF = (c.AF & 0x00ff) | uint16(value)<<8
	c.PC++

	return 7
}

func (c *CPU) decBc() uint8 {
	c.BC--
	c.PC++

	return 6
}

func (c *CPU) incC() uint8 {
	return c.increaseRegister('C')
}

func (c *CPU) decC() uint8 {
	return c.decreaseRegister('C')
}

func (c *CPU) ldCX() uint8 {
	c.PC++
	c.BC = (c.BC & 0xff00) | uint16(c.readByte(c.PC))
	c.PC++

	return 7
}

func (c *CPU) rrca() uint8 {
	a := uint8(c.AF >> 8)
	signed := a&1 == 1
	a = a >> 1
	c.AF = (c.AF & 0x00ff) | (uint16(a) << 8)
	c.PC++

	// C (carry) flag
	if signed {
		c.Flags = c.Flags | 0b00000001
		c.AF = c.AF | 0x8000
	} else {
		c.Flags = c.Flags & 0b11111110
	}
	// H (half carry) flag
	c.Flags = c.Flags & 0b11101111
	// N (sub/add) flag
	c.Flags = c.Flags & 0b11111101

	return 4
}

func (c *CPU) djnzX() uint8 {
	c.BC -= 256
	if c.BC < 256 {
		c.PC += 2
		return 8
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.readByte(c.PC+1))))
	return 13
}

func (c *CPU) ldDeXx() uint8 {
	c.DE = c.readWord(c.PC + 1)
	c.PC += 3

	return 10
}

func (c *CPU) ldDeA() uint8 {
	c.dma.SetMemoryByte(c.DE, uint8(c.AF>>8))
	c.PC++
	return 7
}

func (c *CPU) incDe() uint8 {
	c.DE++
	c.PC++
	return 6
}

func (c *CPU) incD() uint8 {
	return c.increaseRegister('D')
}

func (c *CPU) decD() uint8 {
	return c.decreaseRegister('D')
}

func (c *CPU) ldDX() uint8 {
	c.PC++
	c.DE = (c.DE & 0x00ff) | (uint16(c.readByte(c.PC)) << 8)
	c.PC++

	return 7
}

func (c *CPU) rla() uint8 {
	a := uint8(c.AF >> 8)
	signed := a&128 == 128
	a = a << 1

	if c.Flags&1 == 1 {
		a = a | 0b00000001
	} else {
		a = a & 0b11111110
	}

	c.AF = (c.AF & 0x00ff) | (uint16(a) << 8)
	c.PC++

	// C (carry) flag
	if signed {
		c.Flags = c.Flags | 0b00000001
	} else {
		c.Flags = c.Flags & 0b11111110
	}
	// H (half carry) flag
	c.Flags = c.Flags & 0b11101111
	// N (sub/add) flag
	c.Flags = c.Flags & 0b11111101

	return 4
}

func (c *CPU) jrX() uint8 {
	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.readByte(c.PC+1))))

	return 12
}

func (c *CPU) addHlDe() uint8 {
	return c.addRegisters(&c.HL, &c.DE)
}

func (c *CPU) ldADe() uint8 {
	value := c.dma.GetMemory(c.DE)
	c.AF = (c.AF & 0x00ff) | uint16(value)<<8
	c.PC++

	return 7
}

func (c *CPU) decDe() uint8 {
	c.DE--
	c.PC++

	return 6
}

func (c *CPU) Reset() {
	c.PC = 0
	c.AF = 0
	c.AF_ = 0
	c.BC = 0
	c.DE = 0
	c.HL = 0
	c.Flags = 0
}

func CPUNew(dma *dma.DMA) *CPU {
	cpu := new(CPU)
	cpu.dma = dma
	return cpu
}
