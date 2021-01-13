package cpu

import (
	"math/rand"
	"testing"
	"z80/dma"
	"z80/memory"
	"z80/video"
)

var mem = memory.MemoryNew()
var videoMemoryHandler = video.VideoMemoryHandlerNew()
var dmaX = dma.DMANew(mem, videoMemoryHandler)
var cpu = CPUNew(dmaX)

func checkCpu(t *testing.T, expectedCycles uint8, expected map[string]uint16, instructionCall func() uint8) {
	t.Helper()
	var expectedSP, expectedBC, expectedDE, expectedHL uint16
	var expectedAF_, expectedBC_, expectedDE_, expectedHL_ uint16
	var expectedIX, expectedIY uint16
	var expectedA, expectedFlags, expectedI, expectedR uint8

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

	cycles := instructionCall()

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

	if cycles != expectedCycles {
		t.Errorf("cycles: got %d, want %d", cycles, expectedCycles)
	}
}

func resetAll() {
	cpu.Reset()
	mem.Clear()
}

func TestReadWord(t *testing.T) {
	resetAll()

	dmaX.SetMemoryBulk(0x1234, []uint8{0x78, 0x56})

	got := cpu.readWord(0x1234)
	want := uint16(0x5678)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestWriteWord(t *testing.T) {
	resetAll()

	cpu.writeWord(0x1234, 0x5678)

	gotH, gotL := dmaX.GetMemory(0x1234), dmaX.GetMemory(0x1235)
	wantH, wantL := uint8(0x78), uint8(0x56)

	if gotH != wantH || gotL != wantL {
		t.Errorf("got 0x%x%x, want 0x%x%x", gotH, gotL, wantH, wantL)
	}
}

func TestWriteReadWord(t *testing.T) {
	resetAll()
	cpu.writeWord(0x1234, 0x5678)

	got := cpu.readWord(0x1234)
	want := uint16(0x5678)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	gotH, gotL := dmaX.GetMemory(0x1234), dmaX.GetMemory(0x1235)
	wantH, wantL := uint8(0x78), uint8(0x56)

	if gotH != wantH || gotL != wantL {
		t.Errorf("got 0x%x%x, want 0x%x%x", gotH, gotL, wantH, wantL)
	}
}

func TestNop(t *testing.T) {
	resetAll()
	checkCpu(t, 4, map[string]uint16{"PC": 1}, cpu.nop)
}

func TestLdBcXx(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x01, 0x64, 0x32})

	checkCpu(t, 10, map[string]uint16{"PC": 3, "BC": 0x3264}, cpu.ldBcXx)
}

func TestLd_Bc_A(t *testing.T) {
	resetAll()
	cpu.setAcc(0x7a)
	cpu.BC = 0x1015

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x7a, "BC": 0x1015}, cpu.ld_Bc_A)

	got := dmaX.GetMemory(0x1015)
	want := uint8(0x7A)
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestIncBc(t *testing.T) {
	resetAll()
	cpu.BC = 0x1020

	checkCpu(t, 6, map[string]uint16{"PC": 1, "BC": 0x1021}, cpu.incBc)
}

func TestIncB(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010111)
	cpu.BC = 0x1002

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1102, "Flags": 0b00000001}, cpu.incB)

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.BC = 0xff02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0002, "Flags": 0b01010000}, cpu.incB)

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.BC = 0x7f02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x8002, "Flags": 0b10010100}, cpu.incB)
}

func TestDecB(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010101)
	cpu.BC = 0x0102

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0002, "Flags": 0b01000011}, cpu.decB)

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.BC = 0x0002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0xff02, "Flags": 0b10010010}, cpu.decB)

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.BC = 0x8002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x7f02, "Flags": 0b00010110}, cpu.decB)
}

func TestLdBX(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "BC": 0x6400}, cpu.ldBX)
}

func TestRlca(t *testing.T) {
	resetAll()
	cpu.setAcc(0x8c)
	cpu.setFlags(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x19, "Flags": 0b11000101}, cpu.rlca)

	resetAll()
	cpu.setAcc(0x4d)
	cpu.setFlags(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x9a, "Flags": 0b11000100}, cpu.rlca)
}

func TestExAfAf_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x12)
	cpu.setFlags(0xd7)
	cpu.AF_ = 0x5653
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x56, "AF_": 0x12d7, "Flags": 0x53}, cpu.exAfAf_)
}

func TestAddHlBc(t *testing.T) {
	resetAll()
	cpu.BC = 0xa76c //  1010 0111 0110 1100
	cpu.HL = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "BC": 0xa76c, "HL": 0x009f, "Flags": 0b00010001}, cpu.addHlBc)

	resetAll()
	cpu.BC = 0x7fff
	cpu.HL = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "BC": 0x7fff, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addHlBc)
}

func TestLdA_Bc_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryByte(0x1257, 0x64)
	cpu.setAcc(0xff)
	cpu.BC = 0x1257

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x64, "BC": 0x1257}, cpu.ldA_Bc_)
}

func TestDecBc(t *testing.T) {
	resetAll()
	cpu.BC = 0x1000

	checkCpu(t, 6, map[string]uint16{"PC": 1, "BC": 0x0fff}, cpu.decBc)
}

func TestIncC(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010111)
	cpu.BC = 0x0210

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0211, "Flags": 0b00000001}, cpu.incC)

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.BC = 0x02ff
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0200, "Flags": 0b01010000}, cpu.incC)

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.BC = 0x027f
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0280, "Flags": 0b10010100}, cpu.incC)
}

func TestDecC(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010101)
	cpu.BC = 0x0201

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0200, "Flags": 0b01000011}, cpu.decC)

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.BC = 0x0200
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x02ff, "Flags": 0b10010010}, cpu.decC)

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.BC = 0x0280
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x027f, "Flags": 0b00010110}, cpu.decC)
}

func TestLdCX(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "BC": 0x0064}, cpu.ldCX)
}

func TestRrca(t *testing.T) {
	resetAll()
	cpu.setAcc(0x8d)
	cpu.setFlags(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0xc6, "Flags": 0b11000101}, cpu.rrca)

	resetAll()
	cpu.setAcc(0x4c)
	cpu.setFlags(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x26, "Flags": 0b11000100}, cpu.rrca)
}

func TestDjnzX(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.BC = 0x1234
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 13, map[string]uint16{"PC": 0x37, "BC": 0x1134}, cpu.djnzX)

	resetAll()
	cpu.PC = 3
	cpu.BC = 0x0134
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 8, map[string]uint16{"PC": 0x05, "BC": 0x0034}, cpu.djnzX)

	resetAll()
	cpu.PC = 3
	cpu.BC = 0x0034
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 13, map[string]uint16{"PC": 0x37, "BC": 0xff34}, cpu.djnzX)

	resetAll()
	cpu.PC = 3
	cpu.BC = 0x0534
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0xfb})

	checkCpu(t, 13, map[string]uint16{"PC": 0x00, "BC": 0x0434}, cpu.djnzX)
}

func TestLdDeXx(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x01, 0x64, 0x32})

	checkCpu(t, 10, map[string]uint16{"PC": 3, "DE": 0x3264}, cpu.ldDeXx)
}

func TestLd_De_A(t *testing.T) {
	resetAll()
	cpu.setAcc(0x7a)
	cpu.DE = 0x1015

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x7a, "DE": 0x1015}, cpu.ld_De_A)

	got := dmaX.GetMemory(0x1015)
	want := uint8(0x7A)
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestIncDe(t *testing.T) {
	resetAll()
	cpu.DE = 0x1020

	checkCpu(t, 6, map[string]uint16{"PC": 1, "DE": 0x1021}, cpu.incDe)
}

func TestIncD(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010111)
	cpu.DE = 0x1002

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1102, "Flags": 0b00000001}, cpu.incD)

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.DE = 0xff02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0002, "Flags": 0b01010000}, cpu.incD)

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.DE = 0x7f02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x8002, "Flags": 0b10010100}, cpu.incD)
}

func TestDecD(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010101)
	cpu.DE = 0x0102

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0002, "Flags": 0b01000011}, cpu.decD)

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.DE = 0x0002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0xff02, "Flags": 0b10010010}, cpu.decD)

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.DE = 0x8002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x7f02, "Flags": 0b00010110}, cpu.decD)
}

func TestLdDX(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "DE": 0x6400}, cpu.ldDX)
}

func TestRla(t *testing.T) {
	resetAll()
	cpu.setAcc(0x8c)
	cpu.setFlags(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x18, "Flags": 0b11000101}, cpu.rla)

	resetAll()
	cpu.setAcc(0x4d)
	cpu.setFlags(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x9b, "Flags": 0b11000100}, cpu.rla)
}

func TestJrX(t *testing.T) {
	resetAll()
	cpu.PC = 3
	dmaX.SetMemoryBulk(0x0003, []uint8{0x18, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37}, cpu.jrX)

	resetAll()
	cpu.PC = 3
	dmaX.SetMemoryBulk(0x0003, []uint8{0x18, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37}, cpu.jrX)

	resetAll()
	cpu.PC = 3
	dmaX.SetMemoryBulk(0x0003, []uint8{0x18, 0xfb})

	checkCpu(t, 12, map[string]uint16{"PC": 0x00}, cpu.jrX)
}

func TestAddHlDe(t *testing.T) {
	resetAll()
	cpu.DE = 0xa76c //  1010 0111 0110 1100
	cpu.HL = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "DE": 0xa76c, "HL": 0x009f, "Flags": 0b00010001}, cpu.addHlDe)

	resetAll()
	cpu.DE = 0x7fff
	cpu.HL = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "DE": 0x7fff, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addHlDe)
}

func TestLdA_De_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryByte(0x1257, 0x64)
	cpu.setAcc(0xff)
	cpu.DE = 0x1257

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x64, "DE": 0x1257}, cpu.ldA_De_)
}

func TestDecDe(t *testing.T) {
	resetAll()
	cpu.DE = 0x1000

	checkCpu(t, 6, map[string]uint16{"PC": 1, "DE": 0x0fff}, cpu.decDe)
}

func TestIncE(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010111)
	cpu.DE = 0x0210

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0211, "Flags": 0b00000001}, cpu.incE)

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.DE = 0x02ff
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0200, "Flags": 0b01010000}, cpu.incE)

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.DE = 0x027f
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0280, "Flags": 0b10010100}, cpu.incE)
}

func TestDecE(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010101)
	cpu.DE = 0x0201

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0200, "Flags": 0b01000011}, cpu.decE)

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.DE = 0x0200
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x02ff, "Flags": 0b10010010}, cpu.decE)

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.DE = 0x0280
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x027f, "Flags": 0b00010110}, cpu.decE)
}

func TestLdEX(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "DE": 0x0064}, cpu.ldEX)
}

func TestRra(t *testing.T) {
	resetAll()
	cpu.setAcc(0x8d)
	cpu.setFlags(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x46, "Flags": 0b11000101}, cpu.rra)

	resetAll()
	cpu.setAcc(0x4c)
	cpu.setFlags(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0xa6, "Flags": 0b11000100}, cpu.rra)
}

func TestJrNzX(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b10010111}, cpu.jrNzX)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 7, map[string]uint16{"PC": 0x05, "Flags": 0b11010111}, cpu.jrNzX)
}

func TestLdHlXx(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x01, 0x64, 0x32})

	checkCpu(t, 10, map[string]uint16{"PC": 3, "HL": 0x3264}, cpu.ldHlXx)
}

func TestLd_Xx_Hl(t *testing.T) {
	resetAll()
	cpu.HL = 0x483a
	dmaX.SetMemoryBulk(0x0000, []uint8{0x22, 0x29, 0xb2})

	checkCpu(t, 5, map[string]uint16{"PC": 3, "HL": 0x483a}, cpu.ld_Xx_Hl)

	gotH, gotL := dmaX.GetMemory(0xb229), dmaX.GetMemory(0xb22a)
	wantH, wantL := uint8(0x3a), uint8(0x48)

	if gotH != wantH || gotL != wantL {
		t.Errorf("got 0x%x%x, want 0x%x%x", gotH, gotL, wantH, wantL)
	}
}

func TestIncHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1020

	checkCpu(t, 6, map[string]uint16{"PC": 1, "HL": 0x1021}, cpu.incHl)
}

func TestIncH(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010111)
	cpu.HL = 0x1002

	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x1102, "Flags": 0b00000001}, cpu.incH)

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.HL = 0xff02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0002, "Flags": 0b01010000}, cpu.incH)

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.HL = 0x7f02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x8002, "Flags": 0b10010100}, cpu.incH)
}

func TestDecH(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010101)
	cpu.HL = 0x0102

	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0002, "Flags": 0b01000011}, cpu.decH)

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.HL = 0x0002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0xff02, "Flags": 0b10010010}, cpu.decH)

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.HL = 0x8002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x7f02, "Flags": 0b00010110}, cpu.decH)
}

func TestLdHX(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "HL": 0x6400}, cpu.ldHX)
}

func TestJrZX(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b11010111}, cpu.jrZX)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 7, map[string]uint16{"PC": 0x05, "Flags": 0b10010111}, cpu.jrZX)
}

func TestAddHlHl(t *testing.T) {
	resetAll()
	cpu.HL = 0xae6c
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x5cd8, "Flags": 0b00010001}, cpu.addHlHl)

	resetAll()
	cpu.HL = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addHlHl)
}

func TestLdHl_Xx_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x2a, 0x29, 0xb2})
	dmaX.SetMemoryBulk(0xb229, []uint8{0x37, 0xa1})

	checkCpu(t, 16, map[string]uint16{"PC": 3, "HL": 0xa137}, cpu.ldHl_Xx_)
}

func TestDecHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1000

	checkCpu(t, 6, map[string]uint16{"PC": 1, "HL": 0x0fff}, cpu.decHl)
}

func TestIncL(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010111)
	cpu.HL = 0x0210

	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0211, "Flags": 0b00000001}, cpu.incL)

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.HL = 0x02ff
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0200, "Flags": 0b01010000}, cpu.incL)

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.HL = 0x027f
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0280, "Flags": 0b10010100}, cpu.incL)
}

func TestDecL(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010101)
	cpu.HL = 0x0201

	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0200, "Flags": 0b01000011}, cpu.decL)

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.HL = 0x0200
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x02ff, "Flags": 0b10010010}, cpu.decL)

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.HL = 0x0280
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x027f, "Flags": 0b00010110}, cpu.decL)
}

func TestLdLX(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "HL": 0x0064}, cpu.ldLX)
}

func TestCpl(t *testing.T) {
	resetAll()
	cpu.setFlags(0b00000000)
	cpu.setAcc(0xe7)

	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x18, "Flags": 0b00010010}, cpu.cpl)
}

func TestJrNcX(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b11010110}, cpu.jrNcX)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 7, map[string]uint16{"PC": 0x05, "Flags": 0b11010111}, cpu.jrNcX)
}

func TestLdSpXx(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x01, 0x64, 0x32})

	checkCpu(t, 10, map[string]uint16{"PC": 3, "SP": 0x3264}, cpu.ldSpXx)
}

func TestLd_Xx_A(t *testing.T) {
	resetAll()
	cpu.setAcc(0xd7)
	dmaX.SetMemoryBulk(0x0000, []uint8{0x32, 0x41, 0x31})

	checkCpu(t, 13, map[string]uint16{"PC": 3, "A": 0xd7}, cpu.ld_Xx_A)

	got := dmaX.GetMemory(0x3141)
	want := uint8(0xd7)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestIncSp(t *testing.T) {
	resetAll()
	cpu.SP = 0x1020

	checkCpu(t, 6, map[string]uint16{"PC": 1, "SP": 0x1021}, cpu.incSp)
}

func TestInc_Hl_(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010111)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x25)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b00000001}, cpu.inc_Hl_)

	got := dmaX.GetMemory(0x3572)
	want := uint8(0x26)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0xff)
	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b01010000}, cpu.inc_Hl_)

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x7f)
	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b10010100}, cpu.inc_Hl_)

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x80)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestDec_Hl_(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010101)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x01)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b01000011}, cpu.dec_Hl_)

	got := dmaX.GetMemory(0x3572)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x00)
	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b10010010}, cpu.dec_Hl_)

	got = dmaX.GetMemory(0x3572)
	want = uint8(0xff)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x80)
	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b00010110}, cpu.dec_Hl_)

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x7f)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestLd_Hl_X(t *testing.T) {
	resetAll()
	cpu.HL = 0x1015
	dmaX.SetMemoryBulk(0x0000, []uint8{0x36, 0x28})

	checkCpu(t, 10, map[string]uint16{"PC": 2, "HL": 0x1015}, cpu.ld_Hl_X)

	got := dmaX.GetMemory(0x1015)
	want := uint8(0x28)
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestScf(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010110)

	checkCpu(t, 4, map[string]uint16{"PC": 1, "Flags": 0b11000101}, cpu.scf)

}

func TestJrCX(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x38, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b11010111}, cpu.jrCX)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x38, 0x32})

	checkCpu(t, 7, map[string]uint16{"PC": 0x05, "Flags": 0b11010110}, cpu.jrCX)
}

func TestAddHlSp(t *testing.T) {
	resetAll()
	cpu.SP = 0xa76c //  1010 0111 0110 1100
	cpu.HL = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0xa76c, "HL": 0x009f, "Flags": 0b00010001}, cpu.addHlSp)

	resetAll()
	cpu.SP = 0x7fff
	cpu.HL = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0x7fff, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addHlSp)
}

func TestLdA_Xx_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x3a, 0x57, 0x12})
	dmaX.SetMemoryByte(0x1257, 0x64)
	cpu.setAcc(0xff)

	checkCpu(t, 13, map[string]uint16{"PC": 3, "A": 0x64}, cpu.ldA_Xx_)
}

func TestDecSp(t *testing.T) {
	resetAll()
	cpu.SP = 0x1000

	checkCpu(t, 6, map[string]uint16{"PC": 1, "SP": 0x0fff}, cpu.decSp)
}

func TestIncA(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010111)
	cpu.setAcc(0x10)

	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x11, "Flags": 0b00000001}, cpu.incA)

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.setAcc(0xff)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x00, "Flags": 0b01010000}, cpu.incA)

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.setAcc(0x7f)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x80, "Flags": 0b10010100}, cpu.incA)
}

func TestDecA(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010101)
	cpu.setAcc(0x01)

	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x00, "Flags": 0b01000011}, cpu.decA)

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.setAcc(0x00)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0xff, "Flags": 0b10010010}, cpu.decA)

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.setAcc(0x80)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x7f, "Flags": 0b00010110}, cpu.decA)
}

func TestLdAX(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "A": 0x64}, cpu.ldAX)
}

func TestCcf(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010110)

	checkCpu(t, 4, map[string]uint16{"PC": 1, "Flags": 0b11000101}, cpu.ccf)

	resetAll()
	cpu.setFlags(0b11000111)

	checkCpu(t, 4, map[string]uint16{"PC": 1, "Flags": 0b11010100}, cpu.ccf)
}

func TestLdRR_(t *testing.T) {
	resetAll()
	cpu.BC = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1234}, cpu.ldRR_('B', 'B'))

	resetAll()
	cpu.BC = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x3434}, cpu.ldRR_('B', 'C'))

	resetAll()
	cpu.BC = 0x1234
	cpu.DE = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x5634, "DE": 0x5678}, cpu.ldRR_('B', 'D'))

	resetAll()
	cpu.BC = 0x1234
	cpu.DE = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x7834, "DE": 0x5678}, cpu.ldRR_('B', 'E'))

	resetAll()
	cpu.BC = 0x1234
	cpu.HL = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x5634, "HL": 0x5678}, cpu.ldRR_('B', 'H'))

	resetAll()
	cpu.BC = 0x1234
	cpu.HL = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x7834, "HL": 0x5678}, cpu.ldRR_('B', 'L'))

	resetAll()
	cpu.BC = 0x1234
	cpu.setAcc(0x56)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x5634, "A": 0x56}, cpu.ldRR_('B', 'A'))

	resetAll()
	cpu.BC = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1212}, cpu.ldRR_('C', 'B'))

	resetAll()
	cpu.BC = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1234}, cpu.ldRR_('C', 'C'))

	resetAll()
	cpu.BC = 0x1234
	cpu.DE = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1256, "DE": 0x5678}, cpu.ldRR_('C', 'D'))

	resetAll()
	cpu.BC = 0x1234
	cpu.DE = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1278, "DE": 0x5678}, cpu.ldRR_('C', 'E'))

	resetAll()
	cpu.BC = 0x1234
	cpu.HL = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1256, "HL": 0x5678}, cpu.ldRR_('C', 'H'))

	resetAll()
	cpu.BC = 0x1234
	cpu.HL = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1278, "HL": 0x5678}, cpu.ldRR_('C', 'L'))

	resetAll()
	cpu.BC = 0x1234
	cpu.setAcc(0x56)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1256, "A": 0x56}, cpu.ldRR_('C', 'A'))

	resetAll()
	cpu.BC = 0x5678
	cpu.DE = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x5634, "BC": 0x5678}, cpu.ldRR_('D', 'B'))

	resetAll()
	cpu.BC = 0x5678
	cpu.DE = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x7834, "BC": 0x5678}, cpu.ldRR_('D', 'C'))

	resetAll()
	cpu.DE = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1234}, cpu.ldRR_('D', 'D'))

	resetAll()
	cpu.DE = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x3434}, cpu.ldRR_('D', 'E'))

	resetAll()
	cpu.DE = 0x1234
	cpu.HL = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x5634, "HL": 0x5678}, cpu.ldRR_('D', 'H'))

	resetAll()
	cpu.DE = 0x1234
	cpu.HL = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x7834, "HL": 0x5678}, cpu.ldRR_('D', 'L'))

	resetAll()
	cpu.DE = 0x1234
	cpu.setAcc(0x56)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x5634, "A": 0x56}, cpu.ldRR_('D', 'A'))

	resetAll()
	cpu.BC = 0x5678
	cpu.DE = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1256, "BC": 0x5678}, cpu.ldRR_('E', 'B'))

	resetAll()
	cpu.BC = 0x5678
	cpu.DE = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1278, "BC": 0x5678}, cpu.ldRR_('E', 'C'))

	resetAll()
	cpu.DE = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1212}, cpu.ldRR_('E', 'D'))

	resetAll()
	cpu.DE = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1234}, cpu.ldRR_('E', 'E'))

	resetAll()
	cpu.DE = 0x1234
	cpu.HL = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1256, "HL": 0x5678}, cpu.ldRR_('E', 'H'))

	resetAll()
	cpu.DE = 0x1234
	cpu.HL = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1278, "HL": 0x5678}, cpu.ldRR_('E', 'L'))

	resetAll()
	cpu.DE = 0x1234
	cpu.setAcc(0x56)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1256, "A": 0x56}, cpu.ldRR_('E', 'A'))

	resetAll()
	cpu.BC = 0x5678
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x5634, "BC": 0x5678}, cpu.ldRR_('H', 'B'))

	resetAll()
	cpu.BC = 0x5678
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x7834, "BC": 0x5678}, cpu.ldRR_('H', 'C'))

	resetAll()
	cpu.DE = 0x5678
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x5634, "DE": 0x5678}, cpu.ldRR_('H', 'D'))

	resetAll()
	cpu.DE = 0x5678
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x7834, "DE": 0x5678}, cpu.ldRR_('H', 'E'))

	resetAll()
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x1234}, cpu.ldRR_('H', 'H'))

	resetAll()
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x3434}, cpu.ldRR_('H', 'L'))

	resetAll()
	cpu.HL = 0x1234
	cpu.setAcc(0x56)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x5634, "A": 0x56}, cpu.ldRR_('H', 'A'))

	resetAll()
	cpu.BC = 0x5678
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x1256, "BC": 0x5678}, cpu.ldRR_('L', 'B'))

	resetAll()
	cpu.BC = 0x5678
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x1278, "BC": 0x5678}, cpu.ldRR_('L', 'C'))

	resetAll()
	cpu.DE = 0x5678
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x1256, "DE": 0x5678}, cpu.ldRR_('L', 'D'))

	resetAll()
	cpu.DE = 0x5678
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x1278, "DE": 0x5678}, cpu.ldRR_('L', 'E'))

	resetAll()
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x1212}, cpu.ldRR_('L', 'H'))

	resetAll()
	cpu.HL = 0x1234
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x1234}, cpu.ldRR_('L', 'L'))

	resetAll()
	cpu.HL = 0x1234
	cpu.setAcc(0x56)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x1256, "A": 0x56}, cpu.ldRR_('L', 'A'))

	resetAll()
	cpu.BC = 0x5678
	cpu.setAcc(0x56)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x56, "BC": 0x5678}, cpu.ldRR_('A', 'B'))

	resetAll()
	cpu.BC = 0x5678
	cpu.setAcc(0x56)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x78, "BC": 0x5678}, cpu.ldRR_('A', 'C'))

	resetAll()
	cpu.DE = 0x5678
	cpu.setAcc(0x12)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x56, "DE": 0x5678}, cpu.ldRR_('A', 'D'))

	resetAll()
	cpu.DE = 0x5678
	cpu.setAcc(0x12)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x78, "DE": 0x5678}, cpu.ldRR_('A', 'E'))

	resetAll()
	cpu.HL = 0x5678
	cpu.setAcc(0x12)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x56, "HL": 0x5678}, cpu.ldRR_('A', 'H'))

	resetAll()
	cpu.HL = 0x5678
	cpu.setAcc(0x12)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x78, "HL": 0x5678}, cpu.ldRR_('A', 'L'))

	resetAll()
	cpu.setAcc(0x12)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x12}, cpu.ldRR_('A', 'A'))
}

func TestLdR_Hl_(t *testing.T) {
	resetAll()
	cpu.BC = 0x1234
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "BC": 0xab34, "HL": 0x5678}, cpu.ldR_Ss_('B', "HL"))

	resetAll()
	cpu.BC = 0x1234
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "BC": 0x12ab, "HL": 0x5678}, cpu.ldR_Ss_('C', "HL"))

	resetAll()
	cpu.DE = 0x1234
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "DE": 0xab34, "HL": 0x5678}, cpu.ldR_Ss_('D', "HL"))

	resetAll()
	cpu.DE = 0x1234
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "DE": 0x12ab, "HL": 0x5678}, cpu.ldR_Ss_('E', "HL"))

	resetAll()
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "HL": 0xab78}, cpu.ldR_Ss_('H', "HL"))

	resetAll()
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "HL": 0x56ab}, cpu.ldR_Ss_('L', "HL"))

	resetAll()
	cpu.setAcc(0x12)
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0xab, "HL": 0x5678}, cpu.ldR_Ss_('A', "HL"))
}

func TestLd_Hl_R(t *testing.T) {
	resetAll()
	cpu.BC = 0x1234
	cpu.HL = 0x5678
	checkCpu(t, 7, map[string]uint16{"PC": 1, "BC": 0x1234, "HL": 0x5678}, cpu.ld_Hl_R('B'))

	got := dmaX.GetMemory(0x5678)
	want := uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.BC = 0x1234
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "BC": 0x1234, "HL": 0x5678}, cpu.ld_Hl_R('C'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x34)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.DE = 0x1234
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "DE": 0x1234, "HL": 0x5678}, cpu.ld_Hl_R('D'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.DE = 0x1234
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "DE": 0x1234, "HL": 0x5678}, cpu.ld_Hl_R('E'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x34)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "HL": 0x5678}, cpu.ld_Hl_R('H'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x56)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "HL": 0x5678}, cpu.ld_Hl_R('L'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x78)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.setAcc(0x12)
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x12, "HL": 0x5678}, cpu.ld_Hl_R('A'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}
}

func TestHalt(t *testing.T) {
	resetAll()
	checkCpu(t, 4, map[string]uint16{"PC": 1}, cpu.halt)

	got := cpu.States.Halt
	want := true

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}

func TestAndR(t *testing.T) {
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		resetAll()
		if register == 'A' {
			cpu.setAcc(0x00)
		} else {
			cpu.setAcc(0x56)

		}
		cpu.BC = 0xa9a9
		cpu.DE = 0xa9a9
		cpu.HL = 0xa9a9

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x00, "Flags": 0b01010100}, cpu.andR(register))

		resetAll()
		if register == 'A' {
			cpu.setAcc(0x97)
		} else {
			cpu.setAcc(0xdf)
		}
		cpu.BC = 0xb7b7
		cpu.DE = 0xb7b7
		cpu.HL = 0xb7b7

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x97, "Flags": 0b10010000}, cpu.andR(register))
	}
}

func TestAnd_Hl_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0xa9)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x00, "HL": 0x1234, "Flags": 0b01010100}, cpu.and_Hl_)

	resetAll()
	cpu.setAcc(0xdf)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0xb7)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x97, "HL": 0x1234, "Flags": 0b10010000}, cpu.and_Hl_)
}

func TestXorR(t *testing.T) {
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		resetAll()
		if register == 'A' {
			cpu.setAcc(0x00)
		} else {
			cpu.setAcc(0x56)

		}
		cpu.BC = 0x5656
		cpu.DE = 0x5656
		cpu.HL = 0x5656

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x00, "Flags": 0b01000100}, cpu.xorR(register))

		if register == 'A' {
			continue
		}

		resetAll()
		cpu.setAcc(0x20)
		cpu.BC = 0xb7b7
		cpu.DE = 0xb7b7
		cpu.HL = 0xb7b7

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x97, "Flags": 0b10000000}, cpu.xorR(register))
	}
}

func TestXor_Hl_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0x56)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x00, "HL": 0x1234, "Flags": 0b01000100}, cpu.xor_Hl_)

	resetAll()
	cpu.setAcc(0x20)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0xb7)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x97, "HL": 0x1234, "Flags": 0b10000000}, cpu.xor_Hl_)
}

func TestOrR(t *testing.T) {
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		resetAll()
		cpu.setAcc(0x00)
		cpu.BC = 0x0000
		cpu.DE = 0x0000
		cpu.HL = 0x0000

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x00, "Flags": 0b01000100}, cpu.orR(register))

		resetAll()
		if register == 'A' {
			cpu.setAcc(0x97)
		} else {
			cpu.setAcc(0x84)
		}

		cpu.BC = 0x1313
		cpu.DE = 0x1313
		cpu.HL = 0x1313

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x97, "Flags": 0b10000000}, cpu.orR(register))
	}
}

func TestOr_Hl_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x00)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0x00)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x00, "HL": 0x1234, "Flags": 0b01000100}, cpu.or_Hl_)

	resetAll()
	cpu.setAcc(0x84)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0x13)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x97, "HL": 0x1234, "Flags": 0b10000000}, cpu.or_Hl_)
}

func TestRetNz(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 11, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b10010111}, cpu.retNz)

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 5, map[string]uint16{"PC": 0x1235, "SP": 0xfffc, "Flags": 0b11010111}, cpu.retNz)
}

func TestPopBc(t *testing.T) {
	resetAll()
	cpu.SP = 0xfffe
	cpu.BC = 0x1234
	dmaX.SetMemoryBulk(0xfffe, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 1, "SP": 0x0000, "BC": 0x5678}, cpu.popBc)
}

func TestJpNzXx(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b10010111}, cpu.jpNzXx)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010111}, cpu.jpNzXx)
}

func TestJpXx(t *testing.T) {
	resetAll()
	cpu.PC = 3
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678}, cpu.jpXx)
}

func TestCallNzXx(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b10010111}, cpu.callNzXx)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010111}, cpu.callNzXx)

	gotL, gotH = dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH = uint8(0x00), uint8(0x00)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestPushBc(t *testing.T) {
	resetAll()
	cpu.BC = 0x1234
	cpu.SP = 0x0000
	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0xfffe, "BC": 0x1234}, cpu.pushBc)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestRst(t *testing.T) {
	for _, addr := range []uint8{0x00, 0x08, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38} {
		resetAll()
		cpu.PC = 0x1234
		cpu.SP = 0x0000

		checkCpu(t, 11, map[string]uint16{"PC": uint16(addr), "SP": 0xfffe}, cpu.rst(addr))

		gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
		wantL, wantH := uint8(0x34), uint8(0x12)

		if gotL != wantL || gotH != wantH {
			t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
		}
	}
}

func TestRetZ(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 11, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.retZ)

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 5, map[string]uint16{"PC": 0x1235, "SP": 0xfffc, "Flags": 0b10010111}, cpu.retZ)
}

func TestRet(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "SP": 0xfffe}, cpu.ret)
}

func TestJpZXx(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010111}, cpu.jpZXx)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b10010111}, cpu.jpZXx)
}

func TestCallZXx(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.callZXx)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b10010111}, cpu.callZXx)

	gotL, gotH = dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH = uint8(0x00), uint8(0x00)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestCallXx(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe}, cpu.callXx)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestRetNc(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 11, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010110}, cpu.retNc)

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 5, map[string]uint16{"PC": 0x1235, "SP": 0xfffc, "Flags": 0b11010111}, cpu.retNc)
}

func TestPopDe(t *testing.T) {
	resetAll()
	cpu.SP = 0xfffe
	cpu.DE = 0x1234
	dmaX.SetMemoryBulk(0xfffe, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 1, "SP": 0x0000, "DE": 0x5678}, cpu.popDe)
}

func TestJpNcXx(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010110}, cpu.jpNcXx)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010111}, cpu.jpNcXx)
}

func TestOut_X_A(t *testing.T) {
	resetAll()
	cpu.setAcc(0xaf)
	dmaX.SetMemoryByte(0x0001, 0x45)

	checkCpu(t, 11, map[string]uint16{"PC": 2, "A": 0xaf}, cpu.out_X_A)

	got := cpu.getPort(0x45)
	want := uint8(0xaf)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestCallNcXx(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010110}, cpu.callNcXx)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010111}, cpu.callNcXx)

	gotL, gotH = dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH = uint8(0x00), uint8(0x00)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestPushDe(t *testing.T) {
	resetAll()
	cpu.DE = 0x1234
	cpu.SP = 0x0000
	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0xfffe, "DE": 0x1234}, cpu.pushDe)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestRetC(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 11, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.retC)

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 5, map[string]uint16{"PC": 0x1235, "SP": 0xfffc, "Flags": 0b11010110}, cpu.retC)
}

func TestExx(t *testing.T) {
	resetAll()
	cpu.BC = 0x1234
	cpu.BC_ = 0x4321
	cpu.DE = 0x5678
	cpu.DE_ = 0x8765
	cpu.HL = 0x9abc
	cpu.HL_ = 0xcba9

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x4321, "BC_": 0x1234, "DE": 0x8765, "DE_": 0x5678, "HL": 0xcba9, "HL_": 0x9abc}, cpu.exx)
}

func TestJpCXx(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010111}, cpu.jpCXx)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010110}, cpu.jpCXx)
}

func TestInA_X_(t *testing.T) {
	resetAll()
	cpu.setPort(0x45, 0xaf)
	dmaX.SetMemoryByte(0x0001, 0x45)

	checkCpu(t, 11, map[string]uint16{"PC": 2, "A": 0xaf}, cpu.inA_X_)
}

func TestCallCXx(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.callCXx)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010110}, cpu.callCXx)

	gotL, gotH = dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH = uint8(0x00), uint8(0x00)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestRetPo(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 11, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010011}, cpu.retPo)

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 5, map[string]uint16{"PC": 0x1235, "SP": 0xfffc, "Flags": 0b11010111}, cpu.retPo)
}

func TestPopHl(t *testing.T) {
	resetAll()
	cpu.SP = 0xfffe
	cpu.HL = 0x1234
	dmaX.SetMemoryBulk(0xfffe, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 1, "SP": 0x0000, "HL": 0x5678}, cpu.popHl)
}

func TestJpPoXx(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010011}, cpu.jpPoXx)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010111}, cpu.jpPoXx)
}

func TestEx_Sp_Hl(t *testing.T) {
	resetAll()
	cpu.HL = 0x7012
	cpu.SP = 0x8856
	dmaX.SetMemoryBulk(0x8856, []uint8{0x11, 0x22})

	checkCpu(t, 19, map[string]uint16{"PC": 1, "HL": 0x2211, "SP": 0x8856}, cpu.ex_Sp_Hl)

	gotL, gotH := dmaX.GetMemory(0x8856), dmaX.GetMemory(0x8857)
	wantL, wantH := uint8(0x12), uint8(0x70)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

}

func TestCallPoXx(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010011}, cpu.callPoXx)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010111}, cpu.callPoXx)

	gotL, gotH = dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH = uint8(0x00), uint8(0x00)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestPushHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1234
	cpu.SP = 0x0000
	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0xfffe, "HL": 0x1234}, cpu.pushHl)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestAndX(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	dmaX.SetMemoryByte(0x0001, 0xa9)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x00, "Flags": 0b01010100}, cpu.andX)

	resetAll()
	cpu.setAcc(0xdf)
	dmaX.SetMemoryByte(0x0001, 0xb7)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x97, "Flags": 0b10010000}, cpu.andX)
}

func TestRetPe(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 11, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.retPe)

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 5, map[string]uint16{"PC": 0x1235, "SP": 0xfffc, "Flags": 0b11010011}, cpu.retPe)
}

func TestJp_Hl_(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.HL = 0x1234
	dmaX.SetMemoryBulk(0x1234, []uint8{0x78, 0x56})

	checkCpu(t, 4, map[string]uint16{"PC": 0x5678, "HL": 0x1234}, cpu.jp_Hl_)
}

func TestJpPeXx(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010111}, cpu.jpPeXx)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010011}, cpu.jpPeXx)
}

func TestExDeHl(t *testing.T) {
	resetAll()
	cpu.DE = 0x2822
	cpu.HL = 0x499a

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x499a, "HL": 0x2822}, cpu.exDeHl)
}

func TestCallPeXx(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.callPeXx)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010011}, cpu.callPeXx)

	gotL, gotH = dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH = uint8(0x00), uint8(0x00)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestXorX(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	dmaX.SetMemoryByte(0x0001, 0x56)

	checkCpu(t, 7, map[string]uint16{"PC": 2, "A": 0x00, "Flags": 0b01000100}, cpu.xorX)

	resetAll()
	cpu.setAcc(0x20)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x0001, 0xb7)

	checkCpu(t, 7, map[string]uint16{"PC": 2, "A": 0x97, "Flags": 0b10000000}, cpu.xorX)
}

func TestRetP(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 11, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b01010111}, cpu.retP)

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 5, map[string]uint16{"PC": 0x1235, "SP": 0xfffc, "Flags": 0b11010111}, cpu.retP)
}

func TestPopAf(t *testing.T) {
	resetAll()
	cpu.SP = 0xfffe
	cpu.AF = 0x1200
	dmaX.SetMemoryBulk(0xfffe, []uint8{0xd7, 0xab})

	checkCpu(t, 10, map[string]uint16{"PC": 1, "SP": 0x0000, "A": 0xab, "Flags": 0xd7}, cpu.popAf)
}

func TestJpPXx(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b01010111}, cpu.jpPXx)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010111}, cpu.jpPXx)
}

func TestDi(t *testing.T) {
	resetAll()
	cpu.enableInterrupts()

	checkCpu(t, 4, map[string]uint16{"PC": 1}, cpu.di)

	gotIFF1, gotIFF2 := cpu.checkInterrupts()
	wantIFF1, wantIFF2 := false, false

	if gotIFF1 != wantIFF1 || gotIFF2 != wantIFF2 {
		t.Errorf("got IFF1=%t, IFF2=%t, want IFF1=%t, IFF2=%t", gotIFF1, gotIFF2, wantIFF1, wantIFF2)
	}
}

func TestCallPXx(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b01010111}, cpu.callPXx)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010111}, cpu.callPXx)

	gotL, gotH = dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH = uint8(0x00), uint8(0x00)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestPushAf(t *testing.T) {
	resetAll()
	cpu.AF = 0x12d7
	cpu.SP = 0x0000
	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0xfffe, "A": 0x12, "Flags": 0xd7}, cpu.pushAf)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0xd7), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestOrX(t *testing.T) {
	resetAll()
	cpu.setAcc(0x00)
	dmaX.SetMemoryByte(0x0001, 0x00)

	checkCpu(t, 7, map[string]uint16{"PC": 2, "A": 0x00, "Flags": 0b01000100}, cpu.orX)

	resetAll()
	cpu.setAcc(0x84)
	dmaX.SetMemoryByte(0x0001, 0x13)

	checkCpu(t, 7, map[string]uint16{"PC": 2, "A": 0x97, "Flags": 0b10000000}, cpu.orX)
}

func TestRetM(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 11, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.retM)

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})

	checkCpu(t, 5, map[string]uint16{"PC": 0x1235, "SP": 0xfffc, "Flags": 0b01010111}, cpu.retM)
}

func TestLdSpHl(t *testing.T) {
	resetAll()
	cpu.SP = 0xfffc
	cpu.HL = 0x442e

	checkCpu(t, 6, map[string]uint16{"PC": 1, "SP": 0x442e, "HL": 0x442e}, cpu.ldSpHl)
}

func TestJpMXx(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010111}, cpu.jpZXx)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b01010111}, cpu.jpMXx)
}

func TestEi(t *testing.T) {
	resetAll()
	cpu.disableInterrupts()

	checkCpu(t, 4, map[string]uint16{"PC": 1}, cpu.ei)

	gotIFF1, gotIFF2 := cpu.checkInterrupts()
	wantIFF1, wantIFF2 := true, true

	if gotIFF1 != wantIFF1 || gotIFF2 != wantIFF2 {
		t.Errorf("got IFF1=%t, IFF2=%t, want IFF1=%t, IFF2=%t", gotIFF1, gotIFF2, wantIFF1, wantIFF2)
	}
}

func TestCallMXx(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.callMXx)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b01010111}, cpu.callMXx)

	gotL, gotH = dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH = uint8(0x00), uint8(0x00)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestInR_C_(t *testing.T) {
	expectedRegisterMap := map[byte]string{
		'B': "BC", 'C': "BC", 'D': "DE", 'E': "DE", 'H': "HL", 'L': "HL", 'A': "A",
	}
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A', ' '} {
		expectedValueMap := map[byte]uint16{
			'B': 0x8b34, 'C': 0x008b, 'D': 0x8b00, 'E': 0x008b, 'H': 0x8b00, 'L': 0x008b,
		}

		resetAll()
		cpu.setAcc(0x00)
		cpu.BC = 0x0034
		cpu.DE = 0x0000
		cpu.HL = 0x0000
		cpu.setFlags(0b00000001)
		cpu.setPort(0x34, 0x8b)

		switch register {
		case ' ':
			checkCpu(t, 12, map[string]uint16{"PC": 2, "BC": 0x0034, "Flags": 0b10000101}, cpu.inR_C_(register))
		case 'A':
			checkCpu(t, 12, map[string]uint16{"PC": 2, "BC": 0x0034, "A": 0x8b, "Flags": 0b10000101}, cpu.inR_C_(register))
		case 'B', 'C':
			checkCpu(t, 12, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b10000101}, cpu.inR_C_(register))
		default:
			checkCpu(t, 12, map[string]uint16{"PC": 2, "BC": 0x0034, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b10000101}, cpu.inR_C_(register))
		}

		expectedValueMap = map[byte]uint16{
			'B': 0x0034, 'C': 0xff00, 'D': 0x00ff, 'E': 0xff00, 'H': 0x00ff, 'L': 0xff00,
		}

		resetAll()
		cpu.setAcc(0x00)
		cpu.BC = 0xff34
		cpu.DE = 0xffff
		cpu.HL = 0xffff
		cpu.setFlags(0b00000000)
		cpu.setPort(0x34, 0x00)

		switch register {
		case ' ':
			checkCpu(t, 12, map[string]uint16{"PC": 2, "BC": 0xff34, "Flags": 0b01000100}, cpu.inR_C_(register))
		case 'A':
			checkCpu(t, 12, map[string]uint16{"PC": 2, "BC": 0xff34, "A": 0x00, "Flags": 0b01000100}, cpu.inR_C_(register))
		case 'B', 'C':
			checkCpu(t, 12, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b01000100}, cpu.inR_C_(register))
		default:
			checkCpu(t, 12, map[string]uint16{"PC": 2, "BC": 0xff34, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b01000100}, cpu.inR_C_(register))
		}
	}
}

func TestOut_C_R(t *testing.T) {
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A', ' '} {
		var want uint8

		resetAll()
		cpu.setAcc(0x8b)
		cpu.BC = 0x8b34
		cpu.DE = 0x8b8b
		cpu.HL = 0x8b8b

		checkCpu(t, 12, map[string]uint16{"PC": 2, "A": 0x8b, "BC": 0x8b34, "DE": 0x8b8b, "HL": 0x8b8b}, cpu.out_C_R(register))

		got := cpu.getPort(0x34)

		switch register {
		case ' ':
			want = 0
		case 'C':
			want = 0x34
		default:
			want = 0x8b
		}

		if got != want {
			t.Errorf("%c got %02x, want %02x", register, got, want)
		}
	}
}

func TestLd_Xx_Rr(t *testing.T) {
	for _, registerPair := range [4]string{"BC", "DE", "HL", "SP"} {
		resetAll()
		cpu.BC = 0x4644
		cpu.DE = 0x4644
		cpu.HL = 0x4644
		cpu.SP = 0x4644
		dmaX.SetMemoryBulk(0x0002, []uint8{0x20, 0x10})

		checkCpu(t, 20, map[string]uint16{"PC": 4, "BC": 0x4644, "DE": 0x4644, "HL": 0x4644}, cpu.ld_Xx_Rr(registerPair))

		gotL, gotH := dmaX.GetMemory(0x1020), dmaX.GetMemory(0x1021)
		wantL, wantH := uint8(0x44), uint8(0x46)

		if gotL != wantL || gotH != wantH {
			t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
		}
	}
}

func TestNeg(t *testing.T) {
	// A, ~A, C, N, PV, H, N, Z, S
	var negTruthTable [6][8]uint8 = [6][8]uint8{
		[8]uint8{0, 255, 1, 1, 0, 1, 0, 1},
		[8]uint8{1, 254, 1, 1, 0, 1, 0, 1},
		[8]uint8{127, 128, 1, 1, 0, 1, 0, 1},
		[8]uint8{128, 127, 1, 1, 0, 1, 0, 0},
		[8]uint8{129, 126, 1, 1, 0, 1, 0, 0},
		[8]uint8{255, 0, 1, 1, 0, 1, 1, 0},
	}

	for _, row := range negTruthTable {
		resetAll()
		cpu.setAcc(row[0])
		expectedFlags := 128*row[7] + 64*row[6] + 16*row[5] + 4*row[4] + 2*row[3] + row[2]

		checkCpu(t, 8, map[string]uint16{"PC": 2, "A": uint16(row[1]), "Flags": uint16(expectedFlags)}, cpu.neg)
	}
}

func TestRetn(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})
	cpu.States.IFF1 = true
	cpu.States.IFF2 = false

	checkCpu(t, 14, map[string]uint16{"PC": 0x5678, "SP": 0xfffe}, cpu.retn)

	gotIFF1, gotIFF2 := cpu.checkInterrupts()
	wantIFF1, wantIFF2 := false, false

	if gotIFF1 != wantIFF1 || gotIFF2 != wantIFF2 {
		t.Errorf("got IFF1=%t, IFF2=%t, want IFF1=%t, IFF2=%t", gotIFF1, gotIFF2, wantIFF1, wantIFF2)
	}
}

func TestReti(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0xfffc
	dmaX.SetMemoryBulk(0xfffc, []uint8{0x78, 0x56})
	cpu.States.IFF1 = true
	cpu.States.IFF2 = false

	checkCpu(t, 14, map[string]uint16{"PC": 0x5678, "SP": 0xfffe}, cpu.reti)

	gotIFF1, gotIFF2 := cpu.checkInterrupts()
	wantIFF1, wantIFF2 := false, false

	if gotIFF1 != wantIFF1 || gotIFF2 != wantIFF2 {
		t.Errorf("got IFF1=%t, IFF2=%t, want IFF1=%t, IFF2=%t", gotIFF1, gotIFF2, wantIFF1, wantIFF2)
	}
}

func TestIm(t *testing.T) {
	for im := 0; im <= 2; im++ {
		resetAll()
		cpu.States.IM = uint8(im + 1)

		checkCpu(t, 8, map[string]uint16{"PC": 2}, cpu.im(uint8(im)))

		got := cpu.States.IM
		want := uint8(im)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestLdIA(t *testing.T) {
	resetAll()

	cpu.setAcc(0x45)
	checkCpu(t, 9, map[string]uint16{"PC": 2, "A": 0x45, "I": 0x45}, cpu.ldIA)
}

func TestLdAI(t *testing.T) {
	resetAll()

	cpu.I = 0x99
	cpu.States.IFF2 = false
	cpu.setFlags(0b01010110)

	checkCpu(t, 9, map[string]uint16{"PC": 2, "A": 0x99, "I": 0x99, "Flags": 0b10000000}, cpu.ldAI)

	resetAll()

	cpu.setAcc(0x32)
	cpu.I = 0x00
	cpu.States.IFF2 = true
	cpu.setFlags(0b10010011)

	checkCpu(t, 9, map[string]uint16{"PC": 2, "A": 0x00, "I": 0x00, "Flags": 0b01000101}, cpu.ldAI)
}

func TestLdRr_Xx_(t *testing.T) {
	for _, registerPair := range [4]string{"BC", "DE", "HL", "SP"} {
		resetAll()
		cpu.BC = 0x0123
		cpu.DE = 0x4567
		cpu.HL = 0x89ab
		cpu.SP = 0xcdef
		dmaX.SetMemoryBulk(0x0002, []uint8{0x20, 0x10})
		dmaX.SetMemoryBulk(0x1020, []uint8{0x85, 0x24})

		checkCpu(t, 20, map[string]uint16{"PC": 4, registerPair: 0x2485}, cpu.ldRr_Xx_(registerPair))
	}
}

func TestLdRA(t *testing.T) {
	resetAll()

	cpu.setAcc(0x45)
	checkCpu(t, 9, map[string]uint16{"PC": 2, "A": 0x45, "R": 0x45}, cpu.ldRA)
}

func TestLdAR(t *testing.T) {
	resetAll()

	cpu.R = 0x99
	cpu.States.IFF2 = false
	cpu.setFlags(0b01010110)

	checkCpu(t, 9, map[string]uint16{"PC": 2, "A": 0x99, "R": 0x99, "Flags": 0b10000000}, cpu.ldAR)

	resetAll()

	cpu.setAcc(0x32)
	cpu.R = 0x00
	cpu.States.IFF2 = true
	cpu.setFlags(0b10010011)

	checkCpu(t, 9, map[string]uint16{"PC": 2, "A": 0x00, "R": 0x00, "Flags": 0b01000101}, cpu.ldAR)
}

func TestRrd(t *testing.T) {
	resetAll()

	cpu.setAcc(0x84)
	cpu.HL = 0x5000
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryByte(0x5000, 0x20)

	checkCpu(t, 18, map[string]uint16{"PC": 2, "A": 0x80, "HL": 0x5000, "Flags": 0b10000001}, cpu.rrd)

	got := dmaX.GetMemory(0x5000)
	want := uint8(0x42)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}

}

func TestRld(t *testing.T) {
	resetAll()

	cpu.setAcc(0x7a)
	cpu.HL = 0x5000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x5000, 0x31)

	checkCpu(t, 18, map[string]uint16{"PC": 2, "A": 0x73, "HL": 0x5000, "Flags": 0b00000001}, cpu.rld)

	got := dmaX.GetMemory(0x5000)
	want := uint8(0x1a)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestLdR_IXd_(t *testing.T) {
	resetAll()
	cpu.BC = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "BC": 0xab34, "IX": 0x5678}, cpu.ldR_Ss_('B', "IX"))

	resetAll()
	cpu.BC = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "BC": 0x12ab, "IX": 0x5678}, cpu.ldR_Ss_('C', "IX"))

	resetAll()
	cpu.DE = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "DE": 0xab34, "IX": 0x5678}, cpu.ldR_Ss_('D', "IX"))

	resetAll()
	cpu.DE = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "DE": 0x12ab, "IX": 0x5678}, cpu.ldR_Ss_('E', "IX"))

	resetAll()
	cpu.HL = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "HL": 0xab34, "IX": 0x5678}, cpu.ldR_Ss_('H', "IX"))

	resetAll()
	cpu.HL = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "HL": 0x12ab, "IX": 0x5678}, cpu.ldR_Ss_('L', "IX"))

	resetAll()
	cpu.setAcc(0x12)
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0xab, "IX": 0x5678}, cpu.ldR_Ss_('A', "IX"))
}

func TestLdR_IYd_(t *testing.T) {
	resetAll()
	cpu.BC = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "BC": 0xab34, "IY": 0x5678}, cpu.ldR_Ss_('B', "IY"))

	resetAll()
	cpu.BC = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "BC": 0x12ab, "IY": 0x5678}, cpu.ldR_Ss_('C', "IY"))

	resetAll()
	cpu.DE = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "DE": 0xab34, "IY": 0x5678}, cpu.ldR_Ss_('D', "IY"))

	resetAll()
	cpu.DE = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "DE": 0x12ab, "IY": 0x5678}, cpu.ldR_Ss_('E', "IY"))

	resetAll()
	cpu.HL = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "HL": 0xab34, "IY": 0x5678}, cpu.ldR_Ss_('H', "IY"))

	resetAll()
	cpu.HL = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "HL": 0x12ab, "IY": 0x5678}, cpu.ldR_Ss_('L', "IY"))

	resetAll()
	cpu.setAcc(0x12)
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5691, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0xab, "IY": 0x5678}, cpu.ldR_Ss_('A', "IY"))
}
