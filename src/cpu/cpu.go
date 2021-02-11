package cpu

import (
	"z80/bus"
	"z80/dma"
	"z80/loader"
)

const busSource uint8 = 0

type CPUTrap struct {
	address  uint16
	opcode   uint8
	trapFunc func() bool
}

func (c *CPU) attachTraps() {
	c.traps = []CPUTrap{
		CPUTrap{address: 0x056b, opcode: 0xc0, trapFunc: c.tapeLoad},
		CPUTrap{address: 0x0111, opcode: 0xc0, trapFunc: c.tapeLoad},
	}
}

func (c *CPU) handleTrap(opcode uint8) bool {
	for _, trap := range c.traps {
		if c.PC == trap.address && trap.opcode == opcode {
			return trap.trapFunc()
		}
	}

	return true
}

func (c *CPU) tapeLoad() bool {
	if !c.currentTape.Loaded() {
		return true
	}

	data := c.currentTape.NextBlock()
	length := uint16(len(data))

	if length == 0 {
		// allow tape to be loaded multiple times
		c.currentTape.Rewind()

		c.HL = (c.HL & 0xff00) | 0x0001
		c.AF_ = (c.AF_ & 0xff00) | 0x0001
		c.setC(false)
		c.PC = 0x05e2
		return true
	}

	read := length - 1
	if read > c.DE {
		read = c.DE
	}

	i := c.AF_ >> 8
	c.AF_ = 0x0145
	c.setAcc(0)

	c.HL = (c.HL & 0xff00) | uint16(data[0])
	parity := data[0]
	data = data[1:]

	if parity != uint8(i) {
		// error
		c.setC(false)
		c.BC = (c.BC & 0xff00) | 0x01
		c.HL = (c.HL & 0x00ff) | (uint16(parity) << 8)
		c.DE -= i
		c.IX += i
		c.PC = 0x05e2
		return false
	}

	c.HL = (c.HL & 0xff00) | uint16(data[read-1])

	if c.AF_&0x01 == 0x01 {
		// LOAD
		for i = 0; i < read; i++ {
			parity ^= data[i]
			c.writeByte(c.IX+i, data[i])
		}
	} else {
		// VERIFY
		for i = 0; i < read; i++ {
			parity ^= data[i]
			if data[i] != c.readByte(c.IX+i) {
				c.HL = (c.HL & 0xff00) | uint16(data[i])
				// this is error routine, it repeats few times here, refactor it @todo
				c.setC(false)
				c.BC = (c.BC & 0xff00) | 0x01
				c.HL = (c.HL & 0x00ff) | (uint16(parity) << 8)
				c.DE -= i
				c.IX += i
				c.PC = 0x05e2
				return false
			}
		}
	}

	if c.DE == i && read+1 < length {
		parity ^= data[read]
		c.setAcc(parity)
		// CP 1 start @todo refactor this
		c.setC(true)
		c.adcValueToAcc(1 ^ 0xff)

		c.setAcc(parity)
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())
		c.setF5(false)
		c.setF3(false)
		// CP 1 end
		//fmt.Println("CHECK C", c.getC())
		c.BC = (c.BC & 0x00ff) | 0x01
	} else {
		//fmt.Println("NOK")
		c.BC = (c.BC & 0x00ff) | 0xff
		c.HL = (c.HL & 0xff00) | 0x01
		c.incR('B')
		c.setC(false)
	}

	c.BC = (c.BC & 0xff00) | 0x01
	c.HL = (c.HL & 0x00ff) | (uint16(parity) << 8)
	c.DE -= i
	c.IX += i

	c.PC = 0x05e2
	return false
}

func (c *CPU) InsertTape(tape *loader.TapFile) {
	c.currentTape = tape
}

type CPUConfig struct {
	ContentionDelays []uint8
	FrameLength      uint
}

type CPU struct {
	PC  uint16
	SP  uint16
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

	Q  uint8
	WZ uint16

	Halt bool
	IFF1 bool
	IFF2 bool
	IM   uint8
	IRQ  bool

	Tstates uint

	config      CPUConfig
	dma         *dma.DMA
	io          *bus.IO
	mnemonics   CPUMnemonics
	currentTape *loader.TapFile
	traps       []CPUTrap
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

func (c *CPU) increaseR() {
	c.R = ((c.R + 1) & 0x7f) | (c.R & 0x80)
}

func (c *CPU) pushStack(value uint16) {
	c.SP -= 2
	c.writeWord(c.SP, value)
}

func (c *CPU) popStack() (value uint16) {
	value = c.readWord(c.SP)
	c.SP += 2

	return
}

func (c *CPU) getPort(addressHigh, addressLow uint8, tstates uint) uint8 {
	c.Tstates += tstates

	return c.io.Read(busSource, (uint16(addressHigh)<<8)|uint16(addressLow))
}

func (c *CPU) setPort(addressHigh, addressLow, value uint8, tstates uint) {
	c.Tstates += tstates

	c.io.Write(busSource, (uint16(addressHigh)<<8)|uint16(addressLow), value)
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

func (c *CPU) contendMemory(address uint16, tstates uint) {
	_, contended := c.dma.GetMemoryByte(address)
	if contended {
		c.Tstates += uint(c.config.ContentionDelays[c.Tstates%c.config.FrameLength])
	}

	c.Tstates += tstates
}

func (c *CPU) readOpcode(address uint16) uint8 {
	value, contended := c.dma.GetMemoryByte(address)

	if contended {
		c.Tstates += uint(c.config.ContentionDelays[c.Tstates%c.config.FrameLength])
	}

	c.Tstates += 4

	return value
}

func (c *CPU) readByte(address uint16) uint8 {
	value, contended := c.dma.GetMemoryByte(address)

	if contended {
		c.Tstates += uint(c.config.ContentionDelays[c.Tstates%c.config.FrameLength])
	}

	c.Tstates += 3

	return value
}

func (c *CPU) writeByte(address uint16, value uint8) {
	contended := c.dma.SetMemoryByte(address, value)

	if contended {
		c.Tstates += uint(c.config.ContentionDelays[c.Tstates%c.config.FrameLength])
	}

	c.Tstates += 3
}

// reads word and maintains endianess
// example:
// 0040 34 21
// readWord(0x0040) => 0x1234
func (c *CPU) readWord(address uint16) uint16 {
	return uint16(c.readByte(address)) | uint16(c.readByte(address+1))<<8
}

// writes word to given address and address+1 and maintains endianess
// example:
// writeWord(0x1234, 0x5678)
// 1234  78 56
func (c *CPU) writeWord(address uint16, value uint16) {
	c.writeByte(address+1, uint8(value>>8))
	c.writeByte(address, uint8(value))
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
	c.Q = uint8(c.AF)

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
	c.Q = uint8(c.AF)
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
	c.Q = uint8(c.AF)

	return
}

func (c *CPU) Step() {
	c.increaseR()
	opcode := c.readOpcode(c.PC)
	if !c.handleTrap(opcode) {
		return
	}

	switch opcode {
	case 0xcb:
		c.increaseR()
		opcode = c.readOpcode(c.PC + 1)
		c.mnemonics.xxBITxx[opcode]()
	case 0xdd:
		c.increaseR()
		opcode = c.readOpcode(c.PC + 1)
		switch opcode {
		case 0xcb:
			// read mnemonic directly without worrying about contention, mnemonic itself would handle that
			opcode, _ = c.dma.GetMemoryByte(c.PC + 3)
			c.mnemonics.xxIXBITxx[opcode]()
		default:
			c.mnemonics.xxIXxx[opcode]()
		}
	case 0xed:
		c.increaseR()
		opcode = c.readOpcode(c.PC + 1)
		c.mnemonics.xx80xx[opcode]()
	case 0xfd:
		c.increaseR()
		opcode = c.readOpcode(c.PC + 1)
		switch opcode {
		case 0xcb:
			// read mnemonic directly without worrying about contention, mnemonic itself would handle that
			opcode, _ = c.dma.GetMemoryByte(c.PC + 3)
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

func (c *CPU) HandleInterrupt() bool {
	if !c.IRQ || !c.IFF1 {
		return false
	}

	if c.Halt {
		c.Halt = false
		c.PC++
	}

	c.increaseR()
	c.IFF1, c.IFF2 = false, false
	c.pushStack(c.PC)

	switch c.IM {
	case 0:
		c.PC = 0x0038
		c.Tstates += 13
	case 1:
		c.PC = 0x0038
		c.Tstates += 13
	case 2:
		inttemp := uint16((uint16(c.I) << 8) | 0x00ff)
		c.PC = c.readWord(inttemp)
		c.Tstates += 13
	}

	return true
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
	c.setPort(0x00, 0xfe, snapshot.Border, 0)
}

func NewCPU(io *bus.IO, dma *dma.DMA, config CPUConfig) *CPU {
	cpu := new(CPU)
	cpu.dma = dma
	cpu.io = io
	cpu.config = config

	cpu.initializeMnemonics()
	cpu.attachTraps()
	cpu.Reset()

	return cpu
}
