package cpu

import (
	"z80/dma"
)

var parityTable [256]bool = [256]bool{
	/*	      0     1      2     3      4     5      6     7      8     9      A     B      C     D      E     F */
	/* 0 */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* 1 */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* 2 */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* 3 */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* 4 */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* 5 */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* 6 */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* 7 */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* 8 */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* 9 */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* A */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* B */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* C */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* D */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* E */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* F */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
}

type CPUStates struct {
	Halt  bool
	Ports [256]uint8
	IFF1  bool
	IFF2  bool
	IM    uint8
}

type CPU struct {
	PC     uint16
	SP     uint16
	AF     uint16
	AF_    uint16
	BC     uint16
	BC_    uint16
	DE     uint16
	DE_    uint16
	HL     uint16
	HL_    uint16
	I      uint8
	R      uint8
	States CPUStates

	dma *dma.DMA
}

func (c *CPU) getAcc() uint8 {
	return uint8(c.AF >> 8)
}

func (c *CPU) setAcc(value uint8) {
	c.AF = (c.AF & 0x00ff) | (uint16(value) << 8)
}

func (c *CPU) getS() bool {
	return c.AF&0x0080 == 0x0080
}

func (c *CPU) getZ() bool {
	return c.AF&0x0040 == 0x0040
}

func (c *CPU) getH() bool {
	return c.AF&0x0010 == 0x0010
}

func (c *CPU) getPV() bool {
	return c.AF&0x0004 == 0x0004
}

func (c *CPU) getN() bool {
	return c.AF&0x0002 == 0x0002
}

func (c *CPU) getC() bool {
	return c.AF&0x0001 == 0x0001
}

func (c *CPU) getFlags() uint8 {
	return uint8(c.AF)
}

func (c *CPU) setS(value bool) {

	if value {
		c.AF = c.AF | 0x0080
	} else {
		c.AF = c.AF & 0xff7f
	}
}

func (c *CPU) setZ(value bool) {

	if value {
		c.AF = c.AF | 0x0040
	} else {
		c.AF = c.AF & 0xffbf
	}
}

func (c *CPU) setH(value bool) {

	if value {
		c.AF = c.AF | 0x0010
	} else {
		c.AF = c.AF & 0xffef
	}
}

func (c *CPU) setPV(value bool) {

	if value {
		c.AF = c.AF | 0x0004
	} else {
		c.AF = c.AF & 0xfffb
	}
}

func (c *CPU) setN(value bool) {

	if value {
		c.AF = c.AF | 0x0002
	} else {
		c.AF = c.AF & 0xfffd
	}
}

func (c *CPU) setC(value bool) {

	if value {
		c.AF = c.AF | 0x0001
	} else {
		c.AF = c.AF & 0xfffe
	}
}

func (c *CPU) setFlags(value uint8) {
	c.AF = (c.AF & 0xff00) | uint16(value)
}

func (c *CPU) popStack() (value uint16) {
	value = c.readWord(c.SP)
	c.SP += 2

	return
}

func (c *CPU) pushStack(value uint16) {
	c.SP -= 2
	c.writeWord(c.SP, value)
}

func (c *CPU) getPort(addr uint8) uint8 {
	return c.States.Ports[addr]
}

func (c *CPU) setPort(addr uint8, value uint8) {
	c.States.Ports[addr] = value
}

func (c *CPU) disableInterrupts() {
	c.States.IFF1 = false
	c.States.IFF2 = false
}

func (c *CPU) enableInterrupts() {
	c.States.IFF1 = true
	c.States.IFF2 = true
}

func (c *CPU) checkInterrupts() (bool, bool) {
	return c.States.IFF1, c.States.IFF2
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

func (c *CPU) extractRegister(r byte) (rhigh bool, rvalue uint16) {
	switch r {
	case 'A':
		rhigh, rvalue = true, c.AF
	case 'B', 'C':
		rhigh, rvalue = r == 'B', c.BC
	case 'D', 'E':
		rhigh, rvalue = r == 'D', c.DE
	case 'H', 'L':
		rhigh, rvalue = r == 'H', c.HL
	default:
		panic("Invalid `r` part of the mnemonic")
	}

	return
}

func (c *CPU) extractRegisterPair(rr string) (rvalue uint16) {
	switch rr {
	case "AF":
		rvalue = c.AF
	case "BC":
		rvalue = c.BC
	case "DE":
		rvalue = c.DE
	case "HL":
		rvalue = c.HL
	case "SP":
		rvalue = c.SP
	default:
		panic("Invalid `rr` part of the mnemonic")
	}

	return
}

func (c *CPU) increaseRegister(name rune) uint8 {
	var register uint8

	switch name {
	case 'A':
		c.AF += 256
		register = c.getAcc()
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
	case 'L':
		register = uint8(c.HL) + 1
		c.HL = (c.HL & 0xff00) | uint16(register)
	}

	c.setN(false)
	c.setPV(register == 0x80)
	c.setH(register&0x0f == 0)
	c.setZ(register == 0)
	c.setS(register > 127)
	c.PC++

	return 4
}

func (c *CPU) decreaseRegister(name rune) uint8 {
	var register uint8

	switch name {
	case 'A':
		c.AF -= 256
		register = c.getAcc()
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
	case 'L':
		register = uint8(c.HL) - 1
		c.HL = (c.HL & 0xff00) | uint16(register)
	}

	c.setN(true)
	c.setPV(register == 0x7f)
	c.setH(register&0x0f == 0x0f)
	c.setZ(register == 0)
	c.setS(register > 127)

	c.PC++
	return 4
}

// left stores the result
func (c *CPU) addRegisters(left, right *uint16) uint8 {
	sum := *left + *right

	c.setC(sum < *left || sum < *right)
	c.setN(false)
	c.setH((*left^*right^sum)&0x1000 == 0x1000)

	*left = sum
	c.PC++
	return 11
}

func (c *CPU) adcValueToAcc(value uint8) {
	var carryIn, carryOut uint8

	if c.getC() {
		carryIn = 1
	}

	a := c.getAcc()
	result := a + value + carryIn
	c.setAcc(result)

	if c.getC() {
		c.setC(a >= 0xff-value)
	} else {
		c.setC(a > 0xff-value)
	}

	c.setN(false)

	if c.getC() {
		carryOut = 1
	}

	c.setPV((((result ^ a ^ value) >> 7) ^ carryOut) == 1)

	c.setH((a^value^result)&0x10 == 0x10)
	c.setZ(result == 0)
	c.setS(result > 127)
}

func (c *CPU) adc16bit(addendLeft, addendRight uint16) (result uint16) {
	var carryIn, carryOut uint16

	if c.getC() {
		carryIn = 1
	}

	result = addendLeft + addendRight + carryIn

	if c.getC() {
		c.setC(addendLeft >= 0xffff-addendRight)
	} else {
		c.setC(addendLeft > 0xffff-addendRight)
	}

	c.setN(false)

	if c.getC() {
		carryOut = 1
	}

	c.setPV((((result ^ addendLeft ^ addendRight) >> 15) ^ carryOut) == 1)

	c.setH((addendLeft^addendRight^result)&0x1000 == 0x1000)
	c.setZ(result == 0)
	c.setS(result > 0x7fff)

	return
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

func (c *CPU) ld_Bc_A() uint8 {
	c.dma.SetMemoryByte(c.BC, c.getAcc())
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
	c.BC = (c.BC & 0x00ff) | (uint16(c.dma.GetMemory(c.PC+1)) << 8)
	c.PC += 2

	return 7
}

func (c *CPU) rlca() uint8 {
	a := c.getAcc()
	signed := a&128 == 128
	a = a << 1
	c.PC++

	if signed {
		a = a | 0x01
	}
	c.setAcc(a)
	c.setC(signed)
	c.setN(false)
	c.setH(false)

	return 4
}

func (c *CPU) exAfAf_() uint8 {
	c.AF, c.AF_ = c.AF_, c.AF
	c.PC++

	return 4
}

func (c *CPU) addHlBc() uint8 {
	return c.addRegisters(&c.HL, &c.BC)
}

func (c *CPU) ldA_Bc_() uint8 {
	value := c.dma.GetMemory(c.BC)
	c.setAcc(value)
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
	c.BC = (c.BC & 0xff00) | uint16(c.dma.GetMemory(c.PC+1))
	c.PC += 2

	return 7
}

func (c *CPU) rrca() uint8 {
	a := c.getAcc()
	signed := a&1 == 1
	a = a >> 1
	c.PC++

	if signed {
		a = a | 0x80
		c.AF = c.AF | 0x8000
	}
	c.setAcc(a)
	c.setC(signed)
	c.setN(false)
	c.setH(false)

	return 4
}

func (c *CPU) djnzX() uint8 {
	c.BC -= 256
	if c.BC < 256 {
		c.PC += 2
		return 8
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))
	return 13
}

func (c *CPU) ldDeXx() uint8 {
	c.DE = c.readWord(c.PC + 1)
	c.PC += 3

	return 10
}

func (c *CPU) ld_De_A() uint8 {
	c.dma.SetMemoryByte(c.DE, c.getAcc())
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
	c.DE = (c.DE & 0x00ff) | (uint16(c.dma.GetMemory(c.PC+1)) << 8)
	c.PC += 2

	return 7
}

func (c *CPU) rla() uint8 {
	a := c.getAcc()
	signed := a&128 == 128
	a = a << 1

	if c.getC() {
		a = a | 0b00000001
	} else {
		a = a & 0b11111110
	}

	c.setAcc(a)
	c.PC++

	// C (carry) flag
	c.setC(signed)
	c.setN(false)
	c.setH(false)

	return 4
}

func (c *CPU) jrX() uint8 {
	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))

	return 12
}

func (c *CPU) addHlDe() uint8 {
	return c.addRegisters(&c.HL, &c.DE)
}

func (c *CPU) ldA_De_() uint8 {
	value := c.dma.GetMemory(c.DE)
	c.setAcc(value)
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
	c.DE = (c.DE & 0xff00) | uint16(c.dma.GetMemory(c.PC+1))
	c.PC += 2

	return 7
}

func (c *CPU) rra() uint8 {
	a := c.getAcc()
	signed := a&1 == 1
	a = a >> 1

	if c.getC() {
		a = a | 0b10000000
	} else {
		a = a & 0b01111111
	}

	c.setAcc(a)
	c.PC++

	c.setC(signed)
	c.setN(false)
	c.setH(false)

	return 4
}

func (c *CPU) jrNzX() uint8 {
	if c.getZ() {
		c.PC += 2
		return 7
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))
	return 12
}

func (c *CPU) ldHlXx() uint8 {
	c.HL = c.readWord(c.PC + 1)
	c.PC += 3

	return 10
}

func (c *CPU) ld_Xx_Hl() uint8 {
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
	c.HL = (c.HL & 0x00ff) | (uint16(c.dma.GetMemory(c.PC+1)) << 8)
	c.PC += 2

	return 7
}

func (c *CPU) daa() uint8 {
	t := 0
	a := c.getAcc()

	if c.getH() || (a&0x0f) > 9 {
		t++
	}

	if c.getC() || (a > 0x99) {
		t += 2
		c.setC(true)
	}

	if c.getN() && !c.getH() {
		c.setH(false)
	} else {
		if c.getN() && c.getH() {
			c.setH(a&0x0f < 6)
		} else {
			c.setH(a&0x0f > 9)
		}
	}

	switch t {
	case 1:
		if c.getN() {
			a += 0xfa
		} else {
			a += 0x06
		}
	case 2:
		if c.getN() {
			a += 0xa0
		} else {
			a += 0x60
		}
	case 3:
		if c.getN() {
			a += 0x9a
		} else {
			a += 0x66
		}
	}

	c.setS(a > 127)
	c.setZ(a == 0)
	c.setPV(parityTable[a])

	c.setAcc(a)

	c.PC++

	return 4
}

func (c *CPU) jrZX() uint8 {
	if !c.getZ() {
		c.PC += 2
		return 7
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))
	return 12
}

func (c *CPU) addHlHl() uint8 {
	return c.addRegisters(&c.HL, &c.HL)
}

func (c *CPU) ldHl_Xx_() uint8 {
	c.HL = c.readWord(c.readWord(c.PC + 1))
	c.PC += 3

	return 16
}

func (c *CPU) decHl() uint8 {
	c.HL--
	c.PC++

	return 6
}

func (c *CPU) incL() uint8 {
	return c.increaseRegister('L')
}

func (c *CPU) decL() uint8 {
	return c.decreaseRegister('L')
}

func (c *CPU) ldLX() uint8 {
	c.HL = (c.HL & 0xff00) | uint16(c.dma.GetMemory(c.PC+1))
	c.PC += 2

	return 7
}

func (c *CPU) cpl() uint8 {
	c.setAcc(c.getAcc() ^ 0xff)
	c.PC++
	c.setH(true)
	c.setN(true)

	return 4
}

func (c *CPU) jrNcX() uint8 {
	if c.getC() {
		c.PC += 2
		return 7
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))
	return 12
}

func (c *CPU) ldSpXx() uint8 {
	c.SP = c.readWord(c.PC + 1)
	c.PC += 3

	return 10
}

func (c *CPU) ld_Xx_A() uint8 {
	c.dma.SetMemoryByte(c.readWord(c.PC+1), c.getAcc())
	c.PC += 3
	return 13
}

func (c *CPU) incSP() uint8 {
	c.SP++
	c.PC++
	return 6
}

func (c *CPU) inc_Hl_() uint8 {
	result := c.dma.GetMemory(c.HL) + 1
	c.dma.SetMemoryByte(c.HL, result)
	c.PC++

	c.setN(false)
	c.setPV(result == 0x80)
	c.setH(result&0x0f == 0)
	c.setZ(result == 0)
	c.setS(result > 127)

	return 11
}

func (c *CPU) dec_Hl_() uint8 {
	result := c.dma.GetMemory(c.HL) - 1
	c.dma.SetMemoryByte(c.HL, result)
	c.PC++

	c.setN(true)
	c.setPV(result == 0x7f)
	c.setH(result&0x0f == 0x0f)
	c.setZ(result == 0)
	c.setS(result > 127)

	return 11
}

func (c *CPU) ld_Hl_X() uint8 {
	c.dma.SetMemoryByte(c.HL, c.dma.GetMemory(c.PC+1))
	c.PC += 2
	return 10
}

func (c *CPU) scf() uint8 {
	c.PC++

	c.setC(true)
	c.setN(false)
	c.setH(false)

	return 4
}

func (c *CPU) jrCX() uint8 {
	if !c.getC() {
		c.PC += 2
		return 7
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))
	return 12
}

func (c *CPU) addHlSp() uint8 {
	return c.addRegisters(&c.HL, &c.SP)
}

func (c *CPU) ldA_Xx_() uint8 {
	c.setAcc(c.dma.GetMemory(c.readWord(c.PC + 1)))
	c.PC += 3

	return 13
}

func (c *CPU) decSP() uint8 {
	c.SP--
	c.PC++

	return 6
}

func (c *CPU) incA() uint8 {
	return c.increaseRegister('A')
}

func (c *CPU) decA() uint8 {
	return c.decreaseRegister('A')
}

func (c *CPU) ldAX() uint8 {
	c.setAcc(c.dma.GetMemory(c.PC + 1))
	c.PC += 2

	return 7
}

func (c *CPU) ccf() uint8 {
	c.PC++

	c.setH(c.getC())
	c.setN(false)
	c.setC(!c.getC())

	return 4
}

func (c *CPU) ldRR_(r, r_ byte) func() uint8 {
	var lhigh bool
	var lvalue *uint16

	switch r {
	case 'A':
		lhigh, lvalue = true, &c.AF
	case 'B', 'C':
		lhigh, lvalue = r == 'B', &c.BC
	case 'D', 'E':
		lhigh, lvalue = r == 'D', &c.DE
	case 'H', 'L':
		lhigh, lvalue = r == 'H', &c.HL
	default:
		panic("Invalid `r` part of the mnemonic")
	}

	rhigh, rvalue := c.extractRegister(r_)

	return func() uint8 {
		var right uint8

		if rhigh {
			right = uint8(rvalue >> 8)
		} else {
			right = uint8(rvalue)
		}

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(right) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(right)
		}

		c.PC++

		return 4
	}
}

func (c *CPU) ldR_Hl_(r byte) func() uint8 {
	var lhigh bool
	var lvalue *uint16

	switch r {
	case 'A':
		lhigh, lvalue = true, &c.AF
	case 'B', 'C':
		lhigh, lvalue = r == 'B', &c.BC
	case 'D', 'E':
		lhigh, lvalue = r == 'D', &c.DE
	case 'H', 'L':
		lhigh, lvalue = r == 'H', &c.HL
	default:
		panic("Invalid `r` part of the mnemonic")
	}

	return func() uint8 {
		right := c.dma.GetMemory(c.HL)

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(right) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(right)
		}

		c.PC++

		return 7
	}
}

func (c *CPU) ld_Hl_R(r byte) func() uint8 {
	rhigh, rvalue := c.extractRegister(r)

	return func() uint8 {
		var right uint8

		if rhigh {
			right = uint8(rvalue >> 8)
		} else {
			right = uint8(rvalue)
		}

		c.dma.SetMemoryByte(c.HL, right)

		c.PC++

		return 7
	}
}

func (c *CPU) halt() uint8 {
	c.PC++
	c.States.Halt = true

	return 4
}

func (c *CPU) addAR(r byte) func() uint8 {
	rhigh, rvalue := c.extractRegister(r)

	return func() uint8 {
		c.setC(false)
		if rhigh {
			c.adcValueToAcc(uint8(rvalue >> 8))
		} else {
			c.adcValueToAcc(uint8(rvalue))
		}

		c.PC++

		return 4
	}
}

func (c *CPU) addA_Hl_() uint8 {
	c.setC(false)
	c.adcValueToAcc(c.dma.GetMemory(c.HL))

	c.PC++

	return 7
}

func (c *CPU) adcAR(r byte) func() uint8 {
	rhigh, rvalue := c.extractRegister(r)

	return func() uint8 {
		if rhigh {
			c.adcValueToAcc(uint8(rvalue >> 8))
		} else {
			c.adcValueToAcc(uint8(rvalue))
		}

		c.PC++

		return 4
	}
}

func (c *CPU) adcA_Hl_() uint8 {
	c.adcValueToAcc(c.dma.GetMemory(c.HL))

	c.PC++

	return 7
}

func (c *CPU) subR(r byte) func() uint8 {
	rhigh, rvalue := c.extractRegister(r)

	return func() uint8 {
		c.setC(true)
		if rhigh {
			c.adcValueToAcc(uint8(rvalue>>8) ^ 0xff)
		} else {
			c.adcValueToAcc(uint8(rvalue) ^ 0xff)
		}

		c.PC++
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 4
	}
}

func (c *CPU) sub_Hl_() uint8 {
	c.setC(true)
	c.adcValueToAcc(c.dma.GetMemory(c.HL) ^ 0xff)

	c.PC++
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 7
}

func (c *CPU) sbcR(r byte) func() uint8 {
	rhigh, rvalue := c.extractRegister(r)

	return func() uint8 {
		c.setC(!c.getC())
		if rhigh {
			c.adcValueToAcc(uint8(rvalue>>8) ^ 0xff)
		} else {
			c.adcValueToAcc(uint8(rvalue) ^ 0xff)
		}

		c.PC++
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 4
	}
}

func (c *CPU) sbc_Hl_() uint8 {
	c.setC(!c.getC())
	c.adcValueToAcc(c.dma.GetMemory(c.HL) ^ 0xff)

	c.PC++
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 7
}

func (c *CPU) andR(r byte) func() uint8 {
	rhigh, rvalue := c.extractRegister(r)

	return func() uint8 {
		var result uint8
		if rhigh {
			result = c.getAcc() & uint8(rvalue>>8)
		} else {
			result = c.getAcc() & uint8(rvalue)
		}

		c.PC++
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(true)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)

		return 4
	}
}

func (c *CPU) and_Hl_() uint8 {
	result := c.getAcc() & c.dma.GetMemory(c.HL)

	c.PC++
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(true)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)

	return 7
}

func (c *CPU) xorR(r byte) func() uint8 {
	rhigh, rvalue := c.extractRegister(r)

	return func() uint8 {
		var result uint8
		if rhigh {
			result = c.getAcc() ^ uint8(rvalue>>8)
		} else {
			result = c.getAcc() ^ uint8(rvalue)
		}

		c.PC++
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)

		return 4
	}
}

func (c *CPU) xor_Hl_() uint8 {
	result := c.getAcc() ^ c.dma.GetMemory(c.HL)

	c.PC++
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(false)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)

	return 7
}

func (c *CPU) orR(r byte) func() uint8 {
	rhigh, rvalue := c.extractRegister(r)

	return func() uint8 {
		var result uint8
		if rhigh {
			result = c.getAcc() | uint8(rvalue>>8)
		} else {
			result = c.getAcc() | uint8(rvalue)
		}

		c.PC++
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)

		return 4
	}
}

func (c *CPU) or_Hl_() uint8 {
	result := c.getAcc() | c.dma.GetMemory(c.HL)

	c.PC++
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(false)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)

	return 7
}

func (c *CPU) cpR(r byte) func() uint8 {
	rhigh, rvalue := c.extractRegister(r)

	return func() uint8 {
		acc := c.getAcc()
		c.setC(true)
		if rhigh {
			c.adcValueToAcc(uint8(rvalue>>8) ^ 0xff)
		} else {
			c.adcValueToAcc(uint8(rvalue) ^ 0xff)
		}

		c.PC++
		c.setAcc(acc)
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 4
	}
}

func (c *CPU) cp_Hl_() uint8 {
	acc := c.getAcc()
	c.setC(true)
	c.adcValueToAcc(c.dma.GetMemory(c.HL) ^ 0xff)

	c.PC++
	c.setAcc(acc)
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 7
}

func (c *CPU) retNz() uint8 {
	if c.getZ() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) popBc() uint8 {
	c.BC = c.popStack()
	c.PC++

	return 10
}

func (c *CPU) jpNzXx() uint8 {
	if c.getZ() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) jpXx() uint8 {
	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) callNzXx() uint8 {
	if c.getZ() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) pushBc() uint8 {
	c.pushStack(c.BC)
	c.PC++

	return 11
}

func (c *CPU) addAX() uint8 {
	c.setC(false)
	c.adcValueToAcc(c.dma.GetMemory(c.PC + 1))

	c.PC += 2

	return 7
}

func (c *CPU) rst(p uint8) func() uint8 {
	if p != 0x00 && p != 0x08 && p != 0x10 && p != 0x18 && p != 0x20 && p != 0x28 && p != 0x30 && p != 0x38 {
		panic("Invalid `p` value for RST")
	}

	return func() uint8 {
		c.pushStack(c.PC)
		c.PC = uint16(p)

		return 11
	}
}

func (c *CPU) retZ() uint8 {
	if !c.getZ() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) ret() uint8 {
	c.PC = c.popStack()

	return 10
}

func (c *CPU) jpZXx() uint8 {
	if !c.getZ() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) callZXx() uint8 {
	if !c.getZ() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) callXx() uint8 {
	c.pushStack(c.PC)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) adcAX() uint8 {
	c.adcValueToAcc(c.dma.GetMemory(c.PC + 1))

	c.PC += 2
	return 7
}

func (c *CPU) retNc() uint8 {
	if c.getC() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) popDe() uint8 {
	c.DE = c.popStack()
	c.PC++

	return 10
}

func (c *CPU) jpNcXx() uint8 {
	if c.getC() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) out_X_A() uint8 {
	c.setPort(c.dma.GetMemory(c.PC+1), c.getAcc())

	c.PC += 2
	return 11
}

func (c *CPU) callNcXx() uint8 {
	if c.getC() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) pushDe() uint8 {
	c.pushStack(c.DE)
	c.PC++

	return 11
}

func (c *CPU) subX() uint8 {
	c.setC(true)
	c.adcValueToAcc(c.dma.GetMemory(c.PC+1) ^ 0xff)

	c.PC += 2
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 7
}

func (c *CPU) retC() uint8 {
	if !c.getC() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) exx() uint8 {
	c.BC, c.BC_ = c.BC_, c.BC
	c.DE, c.DE_ = c.DE_, c.DE
	c.HL, c.HL_ = c.HL_, c.HL

	c.PC++
	return 4
}

func (c *CPU) jpCXx() uint8 {
	if !c.getC() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) inA_X_() uint8 {
	c.setAcc(c.getPort(c.dma.GetMemory(c.PC + 1)))

	c.PC += 2
	return 11
}

func (c *CPU) callCXx() uint8 {
	if !c.getC() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) sbcAX() uint8 {
	c.setC(!c.getC())
	c.adcValueToAcc(c.dma.GetMemory(c.PC+1) ^ 0xff)

	c.PC += 2
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 7
}

func (c *CPU) retPo() uint8 {
	if c.getPV() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) popHl() uint8 {
	c.HL = c.popStack()
	c.PC++

	return 10
}

func (c *CPU) jpPoXx() uint8 {
	if c.getPV() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) ex_Sp_Hl() uint8 {
	value := c.readWord(c.SP)
	c.writeWord(c.SP, c.HL)
	c.HL = value

	c.PC++
	return 19
}

func (c *CPU) callPoXx() uint8 {
	if c.getPV() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) pushHl() uint8 {
	c.pushStack(c.HL)
	c.PC++

	return 11
}

func (c *CPU) andX() uint8 {
	result := c.getAcc() & c.dma.GetMemory(c.PC+1)

	c.PC++
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(true)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)

	return 7
}

func (c *CPU) retPe() uint8 {
	if !c.getPV() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) jp_Hl_() uint8 {
	c.PC = c.readWord(c.HL)
	return 4
}

func (c *CPU) jpPeXx() uint8 {
	if !c.getPV() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) exDeHl() uint8 {
	c.DE, c.HL = c.HL, c.DE

	c.PC++
	return 4
}

func (c *CPU) callPeXx() uint8 {
	if !c.getPV() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) xorX() uint8 {
	result := c.getAcc() ^ c.dma.GetMemory(c.PC+1)

	c.PC += 2
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(false)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)

	return 7
}

func (c *CPU) retP() uint8 {
	if c.getS() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) popAf() uint8 {
	c.AF = c.popStack()
	c.PC++

	return 10
}

func (c *CPU) jpPXx() uint8 {
	if c.getS() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) di() uint8 {
	c.disableInterrupts()

	c.PC++
	return 4
}

func (c *CPU) callPXx() uint8 {
	if c.getS() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) pushAf() uint8 {
	c.pushStack(c.AF)
	c.PC++

	return 11
}

func (c *CPU) orX() uint8 {
	result := c.getAcc() | c.dma.GetMemory(c.PC+1)

	c.PC += 2
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(false)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)

	return 7
}

func (c *CPU) retM() uint8 {
	if !c.getS() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) ldSpHl() uint8 {
	c.SP = c.HL

	c.PC++
	return 6
}

func (c *CPU) jpMXx() uint8 {
	if !c.getS() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) ei() uint8 {
	c.enableInterrupts()

	c.PC++
	return 4
}

func (c *CPU) callMXx() uint8 {
	if !c.getS() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) cpX() uint8 {
	acc := c.getAcc()
	c.setC(true)
	c.adcValueToAcc(c.dma.GetMemory(c.PC+1) ^ 0xff)

	c.PC += 2
	c.setAcc(acc)
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 7
}

func (c *CPU) inR_C_(r byte) func() uint8 {
	var lhigh bool
	var lvalue *uint16

	switch r {
	case 'A':
		lhigh, lvalue = true, &c.AF
	case 'B', 'C':
		lhigh, lvalue = r == 'B', &c.BC
	case 'D', 'E':
		lhigh, lvalue = r == 'D', &c.DE
	case 'H', 'L':
		lhigh, lvalue = r == 'H', &c.HL
	case ' ':
	default:
		panic("Invalid `r` part of the mnemonic")
	}

	return func() uint8 {
		result := c.getPort(uint8(c.BC))

		if r != ' ' {
			if lhigh {
				*lvalue = (*lvalue & 0x00ff) | (uint16(result) << 8)
			} else {
				*lvalue = (*lvalue & 0xff00) | uint16(result)
			}
		}

		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)

		c.PC += 2

		return 12
	}

}

func (c *CPU) out_C_R(r byte) func() uint8 {
	var rhigh bool
	var rvalue uint16

	if r == ' ' {
		rvalue = 0
	} else {
		rhigh, rvalue = c.extractRegister(r)
	}

	return func() uint8 {
		if rhigh {
			c.setPort(uint8(c.BC), uint8(rvalue>>8))
		} else {
			c.setPort(uint8(c.BC), uint8(rvalue))
		}

		c.PC += 2
		return 12
	}
}

func (c *CPU) sbcHlRr(rr string) func() uint8 {
	rvalue := c.extractRegisterPair(rr)

	return func() uint8 {
		c.setC(!c.getC())
		c.HL = c.adc16bit(c.HL, rvalue^0xffff)

		c.PC += 2
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 15
	}
}

func (c *CPU) adcHlRr(rr string) func() uint8 {
	rvalue := c.extractRegisterPair(rr)

	return func() uint8 {
		c.HL = c.adc16bit(c.HL, rvalue)

		c.PC += 2
		return 15
	}
}

func (c *CPU) ld_Xx_Rr(rr string) func() uint8 {
	rvalue := c.extractRegisterPair(rr)

	return func() uint8 {
		c.writeWord(c.readWord(c.PC+2), rvalue)

		c.PC += 4
		return 20
	}
}

func (c *CPU) neg() uint8 {
	a := c.getAcc()
	c.setAcc(0)

	c.setC(false)
	c.adcValueToAcc(a ^ 0xff)

	c.PC += 2
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 8
}

func (c *CPU) retn() uint8 {
	c.PC = c.popStack()
	c.States.IFF1 = c.States.IFF2

	return 14
}

func (c *CPU) im(mode uint8) func() uint8 {
	return func() uint8 {
		c.States.IM = mode
		c.PC += 2

		return 8
	}
}

func (c *CPU) ldIA() uint8 {
	c.I = c.getAcc()

	c.PC += 2
	return 9
}

func (c *CPU) ldAI() uint8 {
	c.setAcc(c.I)

	c.setS(c.I > 127)
	c.setZ(c.I == 0)
	c.setH(false)
	c.setPV(c.States.IFF2)
	c.setN(false)

	c.PC += 2
	return 9
}

func (c *CPU) ldRr_Xx_(rr string) func() uint8 {
	var lvalue *uint16

	switch rr {
	case "AF":
		lvalue = &c.AF
	case "BC":
		lvalue = &c.BC
	case "DE":
		lvalue = &c.DE
	case "HL":
		lvalue = &c.HL
	case "SP":
		lvalue = &c.SP
	default:
		panic("Invalid `rr` part of the mnemonic")
	}

	return func() uint8 {
		*lvalue = c.readWord(c.readWord(c.PC + 2))

		c.PC += 4
		return 20
	}
}

func (c *CPU) ldRA() uint8 {
	c.R = c.getAcc()

	c.PC += 2
	return 9
}

func (c *CPU) ldAR() uint8 {
	c.setAcc(c.R)

	c.setS(c.R > 127)
	c.setZ(c.R == 0)
	c.setH(false)
	c.setPV(c.States.IFF2)
	c.setN(false)

	c.PC += 2
	return 9
}

func (c *CPU) rrd() uint8 {
	value := c.dma.GetMemory(c.HL)
	a := c.getAcc()

	c.setAcc((a & 0xf0) | (value & 0x0f))
	value = value >> 4
	value = (a << 4) | value

	c.dma.SetMemoryByte(c.HL, value)
	a = c.getAcc()

	c.setS(a > 127)
	c.setZ(a == 0)
	c.setH(false)
	c.setPV(parityTable[a])
	c.setN(false)

	c.PC += 2
	return 18
}

func (c *CPU) rld() uint8 {
	value := c.dma.GetMemory(c.HL)
	a := c.getAcc()

	c.setAcc((a & 0xf0) | ((value >> 4) & 0x0f))
	value = value << 4
	value = (a & 0x0f) | value

	c.dma.SetMemoryByte(c.HL, value)
	a = c.getAcc()

	c.setS(a > 127)
	c.setZ(a == 0)
	c.setH(false)
	c.setPV(parityTable[a])
	c.setN(false)

	c.PC += 2
	return 18
}

func (c *CPU) Reset() {
	c.PC = 0
	c.SP = 0
	c.AF = 0
	c.AF_ = 0
	c.BC = 0
	c.BC_ = 0
	c.DE = 0
	c.DE_ = 0
	c.HL = 0
	c.HL_ = 0
	c.I = 0
	c.R = 0
	c.States = CPUStates{IFF1: true, IFF2: true}
}

func CPUNew(dma *dma.DMA) *CPU {
	cpu := new(CPU)
	cpu.dma = dma
	cpu.Reset()
	return cpu
}
