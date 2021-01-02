package cpu

import "z80/dma"

type CPU struct {
	PC    uint16
	BC    uint16
	AF    uint16
	Flags uint8 // [  S  ][  Z  ][     ][  H  ][     ][ P/V ][  N  ][  C  ]
	AF_   uint16

	dma *dma.DMA
}

func (c *CPU) readByte(address uint16) uint8 {
	return c.dma.GetMemory(address)
}

func (c *CPU) readWord(address uint16) uint16 {
	return uint16(c.dma.GetMemory(address+1))<<8 | uint16(c.dma.GetMemory(address))
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
	c.BC += 256
	c.PC++

	b := uint8(c.BC >> 8)

	// C (carry) is not set
	// N (sub/add) flag
	c.Flags = c.Flags & 0b11111101
	// P/V (overflow) flag
	if b == 0x80 {
		c.Flags = c.Flags | 0b00000100
	} else {
		c.Flags = c.Flags & 0b11111011
	}
	// H (half carry) flag
	if b&0b00001111 == 0 {
		c.Flags = c.Flags | 0b00010000
	} else {
		c.Flags = c.Flags & 0b11101111
	}
	// Z (zero) flag
	if b == 0 {
		c.Flags = c.Flags | 0b01000000
	} else {
		c.Flags = c.Flags & 0b10111111
	}
	// S (sign) flag
	if b > 127 {
		c.Flags = c.Flags | 0b10000000
	} else {
		c.Flags = c.Flags & 0b01111111
	}
	return 4
}

func (c *CPU) decB() uint8 {
	c.BC -= 256
	c.PC++

	b := uint8(c.BC >> 8)

	// C (carry) is not set
	// N (sub/add flag)
	c.Flags = c.Flags | 0b00000010
	// P/V (overflow) flag
	if b == 0x7f {
		c.Flags = c.Flags | 0b00000100
	} else {
		c.Flags = c.Flags & 0b11111011
	}
	// H (half carry) flag
	if b&0b00001111 == 15 {
		c.Flags = c.Flags | 0b00010000
	} else {
		c.Flags = c.Flags & 0b11101111
	}
	// Z (zero) flag
	if b == 0 {
		c.Flags = c.Flags | 0b01000000
	} else {
		c.Flags = c.Flags & 0b10111111
	}
	// S (sign) flag
	if b > 127 {
		c.Flags = c.Flags | 0b10000000
	} else {
		c.Flags = c.Flags & 0b01111111
	}
	return 4
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

func (c *CPU) Reset() {
	c.AF = 0
	c.PC = 0
	c.AF_ = 0
	c.Flags = 0
	c.BC = 0
}

func CPUNew(dma *dma.DMA) *CPU {
	cpu := new(CPU)
	cpu.dma = dma
	return cpu
}
