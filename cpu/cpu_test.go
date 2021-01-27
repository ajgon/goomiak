package cpu

import (
	"math/rand"
	"testing"
	"z80/dma"
	"z80/machine"
	"z80/memory"
	"z80/video"
)

var mem = memory.NewWritableMemory()
var videoMemoryHandler = video.VideoMemoryHandlerNew()
var dmaX = dma.DMANew(mem, videoMemoryHandler)
var cpu = CPUNew(dmaX, machine.Spectrum48k)

func getMemoryByte(address uint16) (value uint8) {
	value, _ = dmaX.GetMemoryByte(address)

	return
}

func checkCpu(t *testing.T, opcodeSize uint8, expectedTstates uint64, expected map[string]uint16, instructionCall func()) {
	t.Helper()
	var expectedSP, expectedBC, expectedDE, expectedHL uint16
	var expectedAF_, expectedBC_, expectedDE_, expectedHL_ uint16
	var expectedIX, expectedIY, expectedWZ uint16
	var expectedA, expectedFlags, expectedI, expectedR uint8

	switch opcodeSize {
	case 1:
		cpu.tstates += 4
	case 2:
		cpu.tstates += 8
	case 3:
		cpu.tstates += 11
	default:
		panic("Invalid opcode size")
	}

	if sp, ok := expected["SP"]; ok {
		expectedSP = sp
	} else {
		cpu.SP = uint16(rand.Uint32())
		expectedSP = cpu.SP
	}

	if a, ok := expected["A"]; ok {
		expectedA = uint8(a)
	} else {
		cpu.setAcc(uint8(rand.Uint32()))
		expectedA = cpu.getAcc()
	}

	if af_, ok := expected["AF_"]; ok {
		expectedAF_ = af_
	} else {
		cpu.AF_ = uint16(rand.Uint32())
		expectedAF_ = cpu.AF_
	}

	if bc, ok := expected["BC"]; ok {
		expectedBC = bc
	} else {
		cpu.BC = uint16(rand.Uint32())
		expectedBC = cpu.BC
	}

	if bc_, ok := expected["BC_"]; ok {
		expectedBC_ = bc_
	} else {
		cpu.BC_ = uint16(rand.Uint32())
		expectedBC_ = cpu.BC_
	}

	if de, ok := expected["DE"]; ok {
		expectedDE = de
	} else {
		cpu.DE = uint16(rand.Uint32())
		expectedDE = cpu.DE
	}

	if de_, ok := expected["DE_"]; ok {
		expectedDE_ = de_
	} else {
		cpu.DE_ = uint16(rand.Uint32())
		expectedDE_ = cpu.DE_
	}

	if hl, ok := expected["HL"]; ok {
		expectedHL = hl
	} else {
		cpu.HL = uint16(rand.Uint32())
		expectedHL = cpu.HL
	}

	if hl_, ok := expected["HL_"]; ok {
		expectedHL_ = hl_
	} else {
		cpu.HL_ = uint16(rand.Uint32())
		expectedHL_ = cpu.HL_
	}

	if flags, ok := expected["Flags"]; ok {
		expectedFlags = uint8(flags)
	} else {
		expectedFlags = uint8(rand.Uint32() & 0b11010111)
		cpu.setFlags(expectedFlags)
	}

	if i, ok := expected["I"]; ok {
		expectedI = uint8(i)
	} else {
		cpu.I = uint8(rand.Uint32())
		expectedI = cpu.I
	}

	if r, ok := expected["R"]; ok {
		expectedR = uint8(r)
	} else {
		cpu.R = uint8(rand.Uint32())
		expectedR = cpu.R
	}

	if ix, ok := expected["IX"]; ok {
		expectedIX = ix
	} else {
		cpu.IX = uint16(rand.Uint32())
		expectedIX = cpu.IX
	}

	if iy, ok := expected["IY"]; ok {
		expectedIY = iy
	} else {
		cpu.IY = uint16(rand.Uint32())
		expectedIY = cpu.IY
	}

	if wz, ok := expected["WZ"]; ok {
		expectedWZ = wz
	} else {
		cpu.WZ = uint16(rand.Uint32())
		expectedWZ = cpu.WZ
	}

	instructionCall()

	if pc, ok := expected["PC"]; ok {
		if cpu.PC != pc {
			t.Errorf("PC: got %d, want %d", cpu.PC, pc)
		}
	} else {
		panic("Every mnemonic test should validate PC!")
	}

	if cpu.SP != expectedSP {
		t.Errorf("SP: got %x, want %x", cpu.SP, expectedSP)
	}

	if cpu.getAcc() != expectedA {
		t.Errorf("A: got %x, want %x", cpu.getAcc(), expectedA)
	}

	if cpu.AF_ != expectedAF_ {
		t.Errorf("AF': got %x, want %x", cpu.AF_, expectedAF_)
	}

	if cpu.BC != expectedBC {
		t.Errorf("BC: got %x, want %x", cpu.BC, expectedBC)
	}

	if cpu.BC_ != expectedBC_ {
		t.Errorf("BC_: got %x, want %x", cpu.BC_, expectedBC_)
	}

	if cpu.DE != expectedDE {
		t.Errorf("DE: got %x, want %x", cpu.DE, expectedDE)
	}

	if cpu.DE_ != expectedDE_ {
		t.Errorf("DE_: got %x, want %x", cpu.DE_, expectedDE_)
	}

	if cpu.HL != expectedHL {
		t.Errorf("HL: got %x, want %x", cpu.HL, expectedHL)
	}

	if cpu.HL_ != expectedHL_ {
		t.Errorf("HL_: got %x, want %x", cpu.HL_, expectedHL_)
	}

	if cpu.getFlags() != expectedFlags {
		t.Errorf("Flags: got %08b, want %08b", cpu.getFlags(), expectedFlags)
	}

	if cpu.I != expectedI {
		t.Errorf("I: got %x, want %x", cpu.I, expectedI)
	}

	if cpu.R != expectedR {
		t.Errorf("R: got %x, want %x", cpu.R, expectedR)
	}

	if cpu.IX != expectedIX {
		t.Errorf("IX: got %x, want %x", cpu.IX, expectedIX)
	}

	if cpu.IY != expectedIY {
		t.Errorf("IY: got %x, want %x", cpu.IY, expectedIY)
	}

	if cpu.WZ != expectedWZ {
		t.Errorf("WZ: got %x, want %x", cpu.WZ, expectedWZ)
	}

	if cpu.tstates != expectedTstates {
		t.Errorf("Tstates: got %d, want %d", cpu.tstates, expectedTstates)
	}
}

func resetAll() {
	cpu.Reset()
	mem.Clear()
}

func TestReadByte(t *testing.T) {
	resetAll()

	dmaX.SetMemoryByte(0x1234, 0x44)

	gotValue, gotTstates := cpu.readByte(0x1234, 3), cpu.tstates
	wantValue, wantTstates := uint8(0x44), uint64(3)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}

	// contented memory, normal tstate
	dmaX.SetMemoryByte(0x4004, 0x44)
	cpu.tstates = 10000

	gotValue, gotTstates = cpu.readByte(0x4004, 4), cpu.tstates
	wantValue, wantTstates = uint8(0x44), uint64(10004)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}

	// contented memory, contented tstate
	dmaX.SetMemoryByte(0x4004, 0x44)
	cpu.tstates = 44191

	gotValue, gotTstates = cpu.readByte(0x4004, 4), cpu.tstates
	wantValue, wantTstates = uint8(0x44), uint64(44201)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}

	// contented memory, contented tstate wrapped out
	dmaX.SetMemoryByte(0x4004, 0x44)
	cpu.tstates = 114079

	gotValue, gotTstates = cpu.readByte(0x4004, 4), cpu.tstates
	wantValue, wantTstates = uint8(0x44), uint64(114089)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}
}

func TestWriteByte(t *testing.T) {
	resetAll()

	cpu.writeByte(0x1234, 0x55, 3)

	gotValue, gotTstates := getMemoryByte(0x1234), cpu.tstates
	wantValue, wantTstates := uint8(0x55), uint64(3)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}

	// contented memory, normal tstate
	cpu.tstates = 10000
	cpu.writeByte(0x4004, 0x55, 3)

	gotValue, gotTstates = getMemoryByte(0x4004), cpu.tstates
	wantValue, wantTstates = uint8(0x55), uint64(10003)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}

	// contented memory, contented tstate
	cpu.tstates = 44191
	cpu.writeByte(0x4004, 0x55, 4)

	gotValue, gotTstates = getMemoryByte(0x4004), cpu.tstates
	wantValue, wantTstates = uint8(0x55), uint64(44201)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}

	// contented memory, contented tstate wrapped out
	cpu.tstates = 114079
	cpu.writeByte(0x4004, 0x55, 4)

	gotValue, gotTstates = getMemoryByte(0x4004), cpu.tstates
	wantValue, wantTstates = uint8(0x55), uint64(114089)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}
}

func TestReadWord(t *testing.T) {
	resetAll()

	dmaX.SetMemoryBulk(0x1234, []uint8{0x78, 0x56})

	gotValue, gotTstates := cpu.readWord(0x1234, 3, 3), cpu.tstates
	wantValue, wantTstates := uint16(0x5678), uint64(6)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}

	// contended memory, normal cycle
	dmaX.SetMemoryBulk(0x4004, []uint8{0x78, 0x56})
	cpu.tstates = 10000

	gotValue, gotTstates = cpu.readWord(0x4004, 4, 3), cpu.tstates
	wantValue, wantTstates = uint16(0x5678), uint64(10007)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}

	// contended memory, contented cycle
	dmaX.SetMemoryBulk(0x4004, []uint8{0x78, 0x56})
	cpu.tstates = 44191

	gotValue, gotTstates = cpu.readWord(0x4004, 4, 3), cpu.tstates
	wantValue, wantTstates = uint16(0x5678), uint64(44208)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}

	// contented memory, contended cycle wrapped out
	dmaX.SetMemoryBulk(0x4004, []uint8{0x78, 0x56})
	cpu.tstates = 114079

	gotValue, gotTstates = cpu.readWord(0x4004, 4, 3), cpu.tstates
	wantValue, wantTstates = uint16(0x5678), uint64(114096)

	if gotValue != wantValue || gotTstates != wantTstates {
		t.Errorf("got 0x%x/%d, want 0x%x/%d", gotValue, gotTstates, wantValue, wantTstates)
	}
}

func TestWriteWord(t *testing.T) {
	resetAll()

	cpu.writeWord(0x1234, 0x5678, 3, 3)

	gotH, gotL, gotT := getMemoryByte(0x1234), getMemoryByte(0x1235), cpu.tstates
	wantH, wantL, wantT := uint8(0x78), uint8(0x56), uint64(6)

	if gotH != wantH || gotL != wantL || gotT != wantT {
		t.Errorf("got 0x%x%x/%d, want 0x%x%x/%d", gotH, gotL, gotT, wantH, wantL, wantT)
	}

	// contented memory, normal cycle
	cpu.tstates = 10000
	cpu.writeWord(0x4004, 0x5678, 4, 3)

	gotH, gotL, gotT = getMemoryByte(0x4004), getMemoryByte(0x4005), cpu.tstates
	wantH, wantL, wantT = uint8(0x78), uint8(0x56), uint64(10007)

	if gotH != wantH || gotL != wantL || gotT != wantT {
		t.Errorf("got 0x%x%x/%d, want 0x%x%x/%d", gotH, gotL, gotT, wantH, wantL, wantT)
	}

	// contented memory, contented cycle
	cpu.tstates = 44191
	cpu.writeWord(0x4004, 0x5678, 4, 3)

	gotH, gotL, gotT = getMemoryByte(0x4004), getMemoryByte(0x4005), cpu.tstates
	wantH, wantL, wantT = uint8(0x78), uint8(0x56), uint64(44208)

	if gotH != wantH || gotL != wantL || gotT != wantT {
		t.Errorf("got 0x%x%x/%d, want 0x%x%x/%d", gotH, gotL, gotT, wantH, wantL, wantT)
	}

	// contented memory, contented cycle wrapped out
	cpu.tstates = 114079
	cpu.writeWord(0x4004, 0x5678, 4, 3)

	gotH, gotL, gotT = getMemoryByte(0x4004), getMemoryByte(0x4005), cpu.tstates
	wantH, wantL, wantT = uint8(0x78), uint8(0x56), uint64(114096)

	if gotH != wantH || gotL != wantL || gotT != wantT {
		t.Errorf("got 0x%x%x/%d, want 0x%x%x/%d", gotH, gotL, gotT, wantH, wantL, wantT)
	}
}

func TestWriteReadWord(t *testing.T) {
	resetAll()
	cpu.writeWord(0x1234, 0x5678, 3, 3)

	got := cpu.readWord(0x1234, 3, 3)
	want := uint16(0x5678)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	gotH, gotL := getMemoryByte(0x1234), getMemoryByte(0x1235)
	wantH, wantL := uint8(0x78), uint8(0x56)

	if gotH != wantH || gotL != wantL {
		t.Errorf("got 0x%x%x, want 0x%x%x", gotH, gotL, wantH, wantL)
	}
}
