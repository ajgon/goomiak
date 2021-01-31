package cpu

import (
	"fmt"
	"z80/dma"
	"z80/loader"
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

type CPUConfig struct {
	ContentionDelays []uint8
	FrameLength      uint
}

type CPU struct {
	PC  uint16
	SP  uint16
	WZ  uint16
	AF  uint16
	AF_ uint16
	BC  uint16
	BC_ uint16
	DE  uint16
	DE_ uint16
	HL  uint16
	HL_ uint16
	I   uint8
	R   uint8
	IX  uint16
	IY  uint16

	Halt bool
	IFF1 bool
	IFF2 bool
	IM   uint8
	IRQ  bool

	Ports   [65536]uint8
	Tstates uint

	config    CPUConfig
	dma       *dma.DMA
	mnemonics CPUMnemonics
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

func (c *CPU) setF5(value bool) {
	if value {
		c.AF = c.AF | 0x0020
	} else {
		c.AF = c.AF & 0xffdf
	}
}

func (c *CPU) setH(value bool) {
	if value {
		c.AF = c.AF | 0x0010
	} else {
		c.AF = c.AF & 0xffef
	}
}

func (c *CPU) setF3(value bool) {
	if value {
		c.AF = c.AF | 0x0008
	} else {
		c.AF = c.AF & 0xfff7
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

func (c *CPU) pushStack(value uint16) {
	c.SP -= 2
	c.writeWord(c.SP, value, 3, 3)
}

func (c *CPU) popStack() (value uint16) {
	value = c.readWord(c.SP, 3, 3)
	c.SP += 2

	return
}

func (c *CPU) GetPort(addressHigh, addressLow uint8, tstates uint) uint8 {
	var value uint8

	c.Tstates += uint(tstates)

	if addressLow&0x01 == 0x00 { // ULA
		if addressHigh == 0x00 {
			value = c.Ports[0x00fe]
		} else {
			addressLeft := (uint16(addressHigh&0x0f) << 8) | 0xf0fe
			addressRight := (uint16(addressHigh&0xf0) << 8) | 0x0ffe

			valueLeft := c.Ports[addressLeft]
			valueRight := c.Ports[addressRight]

			value = valueLeft & valueRight
		}
	} else {
		address := (uint16(addressHigh) << 8) | uint16(addressLow)

		value = c.Ports[address]
	}

	return value
}

func (c *CPU) SetPort(addressHigh, addressLow, value uint8, tstates uint) {
	c.Tstates += tstates

	if addressLow&0x01 == 0x00 && tstates != 0 { // ULA
		c.Ports[0x00fe] = value
	} else {
		c.Ports[(uint16(addressHigh)<<8)|uint16(addressLow)] = value
	}
}

func (c *CPU) disableInterrupts() {
	c.IFF1 = false
	c.IFF2 = false
}

func (c *CPU) enableInterrupts() {
	c.IFF1 = true
	c.IFF2 = true
}

func (c *CPU) checkInterrupts() (bool, bool) {
	return c.IFF1, c.IFF2
}

func (c *CPU) shiftedAddress(base uint16, shift uint8) uint16 {
	if shift > 127 {
		c.WZ = base + uint16(shift) - 256
	} else {
		c.WZ = base + uint16(shift)
	}

	return c.WZ
}

func (c *CPU) readByte(address uint16, usedTstates uint) uint8 {
	value, contended := c.dma.GetMemoryByte(address)

	if contended {
		c.Tstates += uint(c.config.ContentionDelays[c.Tstates%c.config.FrameLength])
	}

	c.Tstates += usedTstates

	return value
}

func (c *CPU) writeByte(address uint16, value uint8, usedTstates uint) {
	contended := c.dma.SetMemoryByte(address, value)

	if contended && usedTstates > 0 {
		c.Tstates += uint(c.config.ContentionDelays[c.Tstates%c.config.FrameLength])
	}

	c.Tstates += usedTstates
}

// reads word and maintains endianess
// example:
// 0040 34 21
// readWord(0x0040) => 0x1234
func (c *CPU) readWord(address uint16, usedTstates1, usedTstates2 uint) uint16 {
	return uint16(c.readByte(address+1, usedTstates1))<<8 | uint16(c.readByte(address, usedTstates2))
}

// writes word to given address and address+1 and maintains endianess
// example:
// writeWord(0x1234, 0x5678)
// 1234  78 56
func (c *CPU) writeWord(address uint16, value uint16, usedTstates1, usedTstates2 uint) {
	c.writeByte(address, uint8(value), usedTstates1)
	c.writeByte(address+1, uint8(value>>8), usedTstates2)
}

func (c *CPU) extractRegister(r byte) uint8 {
	switch r {
	case 'A':
		return c.getAcc()
	case 'B':
		return uint8(c.BC >> 8)
	case 'C':
		return uint8(c.BC)
	case 'D':
		return uint8(c.DE >> 8)
	case 'E':
		return uint8(c.DE)
	case 'H':
		return uint8(c.HL >> 8)
	case 'L':
		return uint8(c.HL)
	case 'X':
		return uint8(c.IX >> 8)
	case 'x':
		return uint8(c.IX)
	case 'Y':
		return uint8(c.IY >> 8)
	case 'y':
		return uint8(c.IY)
	}

	panic("Invalid `r` part of the mnemonic")
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
	case "IX":
		rvalue = c.IX
	case "IY":
		rvalue = c.IY
	default:
		panic("Invalid `rr` part of the mnemonic")
	}

	return
}

// left stores the result
// @todo replace with adc16bit?
func (c *CPU) addRegisters(left *uint16, right uint16) {
	sum := *left + right

	c.setC(sum < *left || sum < right)
	c.setN(false)
	c.setH((*left^right^sum)&0x1000 == 0x1000)
	c.setF5(sum&0x2000 == 0x2000)
	c.setF3(sum&0x0800 == 0x0800)

	*left = sum
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
	c.setF5(result&0x20 == 0x20)
	c.setF3(result&0x08 == 0x08)
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
	c.setF5(result&0x2000 == 0x2000)
	c.setF3(result&0x0800 == 0x0800)

	return
}

func (c *CPU) DebugStep() (tstates uint8) {
	debugT := c.Tstates % c.config.FrameLength

	opcode := c.readByte(c.PC, 4)
	dbOpcode := opcode

	if dbOpcode == 0xcb || dbOpcode == 0xdd || dbOpcode == 0xed || dbOpcode == 0xfd {
		dbOpcode = c.readByte(c.PC+1, 0)
		if dbOpcode == 0xcb {
			dbOpcode = c.readByte(c.PC+3, 0)
		}
	}

	fmt.Printf(
		"%x: AF=%d BC=%d DE=%d HL=%d AF_=%d BC_=%d DE_=%d HL_=%d IX=%d IY=%d SP=%d PC=%d (HL)=%d t=%d\n",
		dbOpcode, c.AF, c.BC, c.DE, c.HL, c.AF_, c.BC_, c.DE_, c.HL_, c.IX, c.IY, c.SP, c.PC, c.readByte(c.HL, 0), debugT,
	)

	switch opcode {
	case 0xcb:
		opcode = c.readByte(c.PC+1, 4)
		c.mnemonics.xxBITxx[opcode]()
	case 0xdd:
		opcode = c.readByte(c.PC+1, 4)
		switch opcode {
		case 0xcb:
			opcode = c.readByte(c.PC+3, 3)
			c.mnemonics.xxIXBITxx[opcode]()
		default:
			c.mnemonics.xxIXxx[opcode]()
		}
	case 0xed:
		opcode = c.readByte(c.PC+1, 4)
		c.mnemonics.xx80xx[opcode]()
	case 0xfd:
		opcode = c.readByte(c.PC+1, 4)
		switch opcode {
		case 0xcb:
			opcode = c.readByte(c.PC+3, 3)
			c.mnemonics.xxIYBITxx[opcode]()
		default:
			c.mnemonics.xxIYxx[opcode]()
		}
	default:
		c.mnemonics.base[opcode]()
	}

	return
}

func (c *CPU) Step() {
	opcode := c.readByte(c.PC, 4)
	switch opcode {
	case 0xcb:
		opcode = c.readByte(c.PC+1, 4)
		c.mnemonics.xxBITxx[opcode]()
	case 0xdd:
		opcode = c.readByte(c.PC+1, 4)
		switch opcode {
		case 0xcb:
			opcode = c.readByte(c.PC+3, 3)
			c.mnemonics.xxIXBITxx[opcode]()
		default:
			c.mnemonics.xxIXxx[opcode]()
		}
	case 0xed:
		opcode = c.readByte(c.PC+1, 4)
		c.mnemonics.xx80xx[opcode]()
	case 0xfd:
		opcode = c.readByte(c.PC+1, 4)
		switch opcode {
		case 0xcb:
			opcode = c.readByte(c.PC+3, 3)
			c.mnemonics.xxIYBITxx[opcode]()
		default:
			c.mnemonics.xxIYxx[opcode]()
		}
	default:
		c.mnemonics.base[opcode]()
	}

	return
}

func (c *CPU) SetIRQ(state bool) {
	c.IRQ = state
}

func (c *CPU) HandleInterrupt() {
	if !c.IRQ || !c.IFF1 {
		return
	}

	if c.Halt {
		c.Halt = false
		c.PC++
	}

	c.IFF1, c.IFF2 = false, false
	c.pushStack(c.PC)

	switch c.IM {
	case 0:
		panic("IM 0")
	case 1:
		c.PC = 0x0038
		c.Tstates += 7
	case 2:
		inttemp := uint16((uint16(c.I) << 8) | 0x00ff)
		c.PC = c.readWord(inttemp, 3, 3)
		c.Tstates += 7
	}
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
	c.IX = 0
	c.IY = 0

	c.Tstates = 0

	c.Halt = false
	c.IFF1 = true
	c.IFF2 = true
	c.IM = 0
	c.IRQ = false

	for i := 0; i < 65536; i++ {
		c.Ports[i] = 0xff
	}
}

func (c *CPU) LoadSnapshot(snapshot loader.Snapshot) {
	c.Reset()
	c.PC = snapshot.PC
	c.SP = snapshot.SP
	c.AF = snapshot.AF
	c.AF_ = snapshot.AF_
	c.BC = snapshot.BC
	c.BC_ = snapshot.BC_
	c.DE = snapshot.DE
	c.DE_ = snapshot.DE_
	c.HL = snapshot.HL
	c.HL_ = snapshot.HL_
	c.I = snapshot.I
	c.R = snapshot.R
	c.IX = snapshot.IX
	c.IY = snapshot.IY

	c.IM = snapshot.IM
	c.IFF1 = snapshot.IFF1
	c.IFF2 = snapshot.IFF2
}

func NewCPU(dma *dma.DMA, config CPUConfig) *CPU {
	cpu := new(CPU)
	cpu.dma = dma
	cpu.config = config

	cpu.initializeMnemonics()
	cpu.Reset()

	return cpu
}
