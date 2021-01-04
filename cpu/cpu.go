package cpu

import (
	"z80/dma"
)

type CPUFlags struct {
	S  bool
	Z  bool
	H  bool
	PV bool
	N  bool
	C  bool
}

func (cf *CPUFlags) toRegister() uint8 {
	var register uint8 = 0x00
	if cf.S {
		register = register | 0x80
	}
	if cf.Z {
		register = register | 0x40
	}
	if cf.H {
		register = register | 0x10
	}
	if cf.PV {
		register = register | 0x04
	}
	if cf.N {
		register = register | 0x02
	}
	if cf.C {
		register = register | 0x01
	}

	return register
}

func (cf *CPUFlags) fromRegister(register uint8) {
	cf.S = register&0x80 == 0x80
	cf.Z = register&0x40 == 0x40
	cf.H = register&0x10 == 0x10
	cf.PV = register&0x04 == 0x04
	cf.N = register&0x02 == 0x02
	cf.C = register&0x01 == 0x01
}

type CPU struct {
	PC    uint16
	AF    uint16
	AF_   uint16
	BC    uint16
	DE    uint16
	HL    uint16
	Flags CPUFlags

	dma *dma.DMA
}

func (c *CPU) readByte(address uint16) uint8 {
	return c.dma.GetMemory(address)
}

// reads word and maintains endianess
// example:
// 0040 34 21
// readWord(0x0040) => 0x1234
func (c *CPU) readWord(address uint16) uint16 {
	return uint16(c.dma.GetMemory(address+1))<<8 | uint16(c.dma.GetMemory(address))
}

// writes word to given address and address+1 and maintains endianess
// example:
// writeWord(0x1234, 0x5678)
// 1234  78 56
func (c *CPU) writeWord(address uint16, value uint16) {
	c.dma.SetMemoryBulk(address, []uint8{uint8(value), uint8(value >> 8)})
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
	case 'E':
		register = uint8(c.DE) + 1
		c.DE = (c.DE & 0xff00) | uint16(register)
	case 'H':
		c.HL += 256
		register = uint8(c.HL >> 8)
	}

	c.Flags.N = false
	c.Flags.PV = register == 0x80
	c.Flags.H = register&0x0f == 0
	c.Flags.Z = register == 0
	c.Flags.S = register > 127
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
	case 'E':
		register = uint8(c.DE) - 1
		c.DE = (c.DE & 0xff00) | uint16(register)
	case 'H':
		c.HL -= 256
		register = uint8(c.HL >> 8)
	}

	c.Flags.N = true
	c.Flags.PV = register == 0x7f
	c.Flags.H = register&0x0f == 0x0f
	c.Flags.Z = register == 0
	c.Flags.S = register > 127

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
	c.BC = (c.BC & 0x00ff) | (uint16(c.readByte(c.PC+1)) << 8)
	c.PC += 2

	return 7
}

func (c *CPU) rlca() uint8 {
	a := uint8(c.AF >> 8)
	signed := a&128 == 128
	a = a << 1
	c.AF = (c.AF & 0x00ff) | (uint16(a) << 8)
	c.PC++

	if signed {
		c.AF = c.AF | 0x0100
	}
	c.Flags.C = signed
	c.Flags.H = false
	c.Flags.N = false

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

	c.Flags.C = sum < *left || sum < *right
	c.Flags.N = false
	c.Flags.H = (*left^*right^sum)&0x1000 == 0x1000

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
	c.BC = (c.BC & 0xff00) | uint16(c.readByte(c.PC+1))
	c.PC += 2

	return 7
}

func (c *CPU) rrca() uint8 {
	a := uint8(c.AF >> 8)
	signed := a&1 == 1
	a = a >> 1
	c.AF = (c.AF & 0x00ff) | (uint16(a) << 8)
	c.PC++

	if signed {
		c.AF = c.AF | 0x8000
	}
	c.Flags.C = signed
	c.Flags.H = false
	c.Flags.N = false

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
	c.DE = (c.DE & 0x00ff) | (uint16(c.readByte(c.PC+1)) << 8)
	c.PC += 2

	return 7
}

func (c *CPU) rla() uint8 {
	a := uint8(c.AF >> 8)
	signed := a&128 == 128
	a = a << 1

	if c.Flags.C {
		a = a | 0b00000001
	} else {
		a = a & 0b11111110
	}

	c.AF = (c.AF & 0x00ff) | (uint16(a) << 8)
	c.PC++

	// C (carry) flag
	c.Flags.C = signed
	c.Flags.H = false
	c.Flags.N = false

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

func (c *CPU) incE() uint8 {
	return c.increaseRegister('E')
}

func (c *CPU) decE() uint8 {
	return c.decreaseRegister('E')
}

func (c *CPU) ldEX() uint8 {
	c.DE = (c.DE & 0xff00) | uint16(c.readByte(c.PC+1))
	c.PC += 2

	return 7
}

func (c *CPU) rra() uint8 {
	a := uint8(c.AF >> 8)
	signed := a&1 == 1
	a = a >> 1

	if c.Flags.C {
		a = a | 0b10000000
	} else {
		a = a & 0b01111111
	}

	c.AF = (c.AF & 0x00ff) | (uint16(a) << 8)
	c.PC++

	c.Flags.C = signed
	c.Flags.H = false
	c.Flags.N = false

	return 4
}

func (c *CPU) jrNzX() uint8 {
	if c.Flags.Z {
		c.PC += 2
		return 7
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.readByte(c.PC+1))))
	return 12
}

func (c *CPU) ldHlXx() uint8 {
	c.HL = c.readWord(c.PC + 1)
	c.PC += 3

	return 10
}

func (c *CPU) ldXxHl() uint8 {
	c.writeWord(c.readWord(c.PC+1), c.HL)
	c.PC += 3
	return 5
}

func (c *CPU) incHl() uint8 {
	c.HL++
	c.PC++
	return 6
}

func (c *CPU) incH() uint8 {
	return c.increaseRegister('H')
}

func (c *CPU) decH() uint8 {
	return c.decreaseRegister('H')
}

func (c *CPU) ldHX() uint8 {
	c.HL = (c.HL & 0x00ff) | (uint16(c.readByte(c.PC+1)) << 8)
	c.PC += 2

	return 7
}

func (c *CPU) Reset() {
	c.PC = 0
	c.AF = 0
	c.AF_ = 0
	c.BC = 0
	c.DE = 0
	c.HL = 0
	c.Flags = CPUFlags{}
}

func CPUNew(dma *dma.DMA) *CPU {
	cpu := new(CPU)
	cpu.dma = dma
	return cpu
}
