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

func TestLdBcNn(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x01, 0x64, 0x32})

	checkCpu(t, 10, map[string]uint16{"PC": 3, "BC": 0x3264}, cpu.ldBcNn)
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

func TestLdBN(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "BC": 0x6400}, cpu.ldBN)
}

func TestRlca(t *testing.T) {
	resetAll()
	cpu.setAcc(0x8c)
	cpu.setFlags(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x19, "Flags": 0b00000001}, cpu.rlcR(' '))

	resetAll()
	cpu.setAcc(0x4d)
	cpu.setFlags(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x9a, "Flags": 0b10000100}, cpu.rlcR(' '))
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

	checkCpu(t, 11, map[string]uint16{"PC": 1, "BC": 0xa76c, "HL": 0x009f, "Flags": 0b00010001}, cpu.addSsRr("HL", "BC"))

	resetAll()
	cpu.BC = 0x7fff
	cpu.HL = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "BC": 0x7fff, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("HL", "BC"))
}

func TestAddIxBc(t *testing.T) {
	resetAll()
	cpu.BC = 0xa76c //  1010 0111 0110 1100
	cpu.IX = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "BC": 0xa76c, "IX": 0x009f, "Flags": 0b00010001}, cpu.addSsRr("IX", "BC"))

	resetAll()
	cpu.BC = 0x7fff
	cpu.IX = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "BC": 0x7fff, "IX": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("IX", "BC"))
}

func TestAddIyBc(t *testing.T) {
	resetAll()
	cpu.BC = 0xa76c //  1010 0111 0110 1100
	cpu.IY = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "BC": 0xa76c, "IY": 0x009f, "Flags": 0b00010001}, cpu.addSsRr("IY", "BC"))

	resetAll()
	cpu.BC = 0x7fff
	cpu.IY = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "BC": 0x7fff, "IY": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("IY", "BC"))
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

func TestLdCN(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "BC": 0x0064}, cpu.ldCN)
}

func TestRrca(t *testing.T) {
	resetAll()
	cpu.setAcc(0x8d)
	cpu.setFlags(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0xc6, "Flags": 0b10000101}, cpu.rrcR(' '))

	resetAll()
	cpu.setAcc(0x4c)
	cpu.setFlags(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x26, "Flags": 0b00000000}, cpu.rrcR(' '))
}

func TestDjnzN(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.BC = 0x1234
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 13, map[string]uint16{"PC": 0x37, "BC": 0x1134}, cpu.djnzN)

	resetAll()
	cpu.PC = 3
	cpu.BC = 0x0134
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 8, map[string]uint16{"PC": 0x05, "BC": 0x0034}, cpu.djnzN)

	resetAll()
	cpu.PC = 3
	cpu.BC = 0x0034
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 13, map[string]uint16{"PC": 0x37, "BC": 0xff34}, cpu.djnzN)

	resetAll()
	cpu.PC = 3
	cpu.BC = 0x0534
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0xfb})

	checkCpu(t, 13, map[string]uint16{"PC": 0x00, "BC": 0x0434}, cpu.djnzN)
}

func TestLdDeNn(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x01, 0x64, 0x32})

	checkCpu(t, 10, map[string]uint16{"PC": 3, "DE": 0x3264}, cpu.ldDeNn)
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

func TestLdDN(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "DE": 0x6400}, cpu.ldDN)
}

func TestRla(t *testing.T) {
	resetAll()
	cpu.setAcc(0x8c)
	cpu.setFlags(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x18, "Flags": 0b00000101}, cpu.rlR(' '))

	resetAll()
	cpu.setAcc(0x4d)
	cpu.setFlags(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x9b, "Flags": 0b10000000}, cpu.rlR(' '))
}

func TestJrN(t *testing.T) {
	resetAll()
	cpu.PC = 3
	dmaX.SetMemoryBulk(0x0003, []uint8{0x18, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37}, cpu.jrN)

	resetAll()
	cpu.PC = 3
	dmaX.SetMemoryBulk(0x0003, []uint8{0x18, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37}, cpu.jrN)

	resetAll()
	cpu.PC = 3
	dmaX.SetMemoryBulk(0x0003, []uint8{0x18, 0xfb})

	checkCpu(t, 12, map[string]uint16{"PC": 0x00}, cpu.jrN)
}

func TestAddHlDe(t *testing.T) {
	resetAll()
	cpu.DE = 0xa76c //  1010 0111 0110 1100
	cpu.HL = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "DE": 0xa76c, "HL": 0x009f, "Flags": 0b00010001}, cpu.addSsRr("HL", "DE"))

	resetAll()
	cpu.DE = 0x7fff
	cpu.HL = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "DE": 0x7fff, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("HL", "DE"))
}

func TestAddIxDe(t *testing.T) {
	resetAll()
	cpu.DE = 0xa76c //  1010 0111 0110 1100
	cpu.IX = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "DE": 0xa76c, "IX": 0x009f, "Flags": 0b00010001}, cpu.addSsRr("IX", "DE"))

	resetAll()
	cpu.DE = 0x7fff
	cpu.IX = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "DE": 0x7fff, "IX": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("IX", "DE"))
}

func TestAddIyDe(t *testing.T) {
	resetAll()
	cpu.DE = 0xa76c //  1010 0111 0110 1100
	cpu.IY = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "DE": 0xa76c, "IY": 0x009f, "Flags": 0b00010001}, cpu.addSsRr("IY", "DE"))

	resetAll()
	cpu.DE = 0x7fff
	cpu.IY = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "DE": 0x7fff, "IY": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("IY", "DE"))
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

func TestLdEN(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "DE": 0x0064}, cpu.ldEN)
}

func TestRra(t *testing.T) {
	resetAll()
	cpu.setAcc(0x8d)
	cpu.setFlags(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x46, "Flags": 0b00000001}, cpu.rrR(' '))

	resetAll()
	cpu.setAcc(0x4c)
	cpu.setFlags(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0xa6, "Flags": 0b10000100}, cpu.rrR(' '))
}

func TestJrNzN(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b10010111}, cpu.jrNzN)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 7, map[string]uint16{"PC": 0x05, "Flags": 0b11010111}, cpu.jrNzN)
}

func TestLdHlNn(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x01, 0x64, 0x32})

	checkCpu(t, 10, map[string]uint16{"PC": 3, "HL": 0x3264}, cpu.ldSsNn("HL"))
}

func TestLd_Nn_Hl(t *testing.T) {
	resetAll()
	cpu.HL = 0x483a
	dmaX.SetMemoryBulk(0x0000, []uint8{0x22, 0x29, 0xb2})

	checkCpu(t, 5, map[string]uint16{"PC": 3, "HL": 0x483a}, cpu.ld_Nn_Ss("HL"))

	gotH, gotL := dmaX.GetMemory(0xb229), dmaX.GetMemory(0xb22a)
	wantH, wantL := uint8(0x3a), uint8(0x48)

	if gotH != wantH || gotL != wantL {
		t.Errorf("got 0x%x%x, want 0x%x%x", gotH, gotL, wantH, wantL)
	}
}

func TestIncHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1020

	checkCpu(t, 6, map[string]uint16{"PC": 1, "HL": 0x1021}, cpu.incSs("HL"))
}

func TestIncIx(t *testing.T) {
	resetAll()
	cpu.IX = 0x1020

	checkCpu(t, 10, map[string]uint16{"PC": 2, "IX": 0x1021}, cpu.incSs("IX"))
}

func TestIncIy(t *testing.T) {
	resetAll()
	cpu.IY = 0x1020

	checkCpu(t, 10, map[string]uint16{"PC": 2, "IY": 0x1021}, cpu.incSs("IY"))
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

func TestLdHN(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "HL": 0x6400}, cpu.ldHN)
}

func TestJrZN(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b11010111}, cpu.jrZN)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 7, map[string]uint16{"PC": 0x05, "Flags": 0b10010111}, cpu.jrZN)
}

func TestAddHlHl(t *testing.T) {
	resetAll()
	cpu.HL = 0xae6c
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x5cd8, "Flags": 0b00010001}, cpu.addSsRr("HL", "HL"))

	resetAll()
	cpu.HL = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("HL", "HL"))
}

func TestAddIxIx(t *testing.T) {
	resetAll()
	cpu.IX = 0xae6c
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "IX": 0x5cd8, "Flags": 0b00010001}, cpu.addSsRr("IX", "IX"))

	resetAll()
	cpu.IX = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "IX": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("IX", "IX"))
}

func TestAddIyIy(t *testing.T) {
	resetAll()
	cpu.IY = 0xae6c
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "IY": 0x5cd8, "Flags": 0b00010001}, cpu.addSsRr("IY", "IY"))

	resetAll()
	cpu.IY = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "IY": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("IY", "IY"))
}

func TestLdHl_Nn_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x2a, 0x29, 0xb2})
	dmaX.SetMemoryBulk(0xb229, []uint8{0x37, 0xa1})

	checkCpu(t, 16, map[string]uint16{"PC": 3, "HL": 0xa137}, cpu.ldSs_Nn_("HL"))
}

func TestDecHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1000

	checkCpu(t, 6, map[string]uint16{"PC": 1, "HL": 0x0fff}, cpu.decSs("HL"))
}

func TestDecIx(t *testing.T) {
	resetAll()
	cpu.IX = 0x1000

	checkCpu(t, 10, map[string]uint16{"PC": 2, "IX": 0x0fff}, cpu.decSs("IX"))
}

func TestDecIy(t *testing.T) {
	resetAll()
	cpu.IY = 0x1000

	checkCpu(t, 10, map[string]uint16{"PC": 2, "IY": 0x0fff}, cpu.decSs("IY"))
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

func TestLdLN(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "HL": 0x0064}, cpu.ldLN)
}

func TestCpl(t *testing.T) {
	resetAll()
	cpu.setFlags(0b00000000)
	cpu.setAcc(0xe7)

	checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x18, "Flags": 0b00010010}, cpu.cpl)
}

func TestJrNcN(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b11010110}, cpu.jrNcN)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 7, map[string]uint16{"PC": 0x05, "Flags": 0b11010111}, cpu.jrNcN)
}

func TestLdSpNn(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x01, 0x64, 0x32})

	checkCpu(t, 10, map[string]uint16{"PC": 3, "SP": 0x3264}, cpu.ldSpNn)
}

func TestLd_Nn_A(t *testing.T) {
	resetAll()
	cpu.setAcc(0xd7)
	dmaX.SetMemoryBulk(0x0000, []uint8{0x32, 0x41, 0x31})

	checkCpu(t, 13, map[string]uint16{"PC": 3, "A": 0xd7}, cpu.ld_Nn_A)

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

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b00000001}, cpu.inc_Ss_("HL"))

	got := dmaX.GetMemory(0x3572)
	want := uint8(0x26)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0xff)
	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b01010000}, cpu.inc_Ss_("HL"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x7f)
	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b10010100}, cpu.inc_Ss_("HL"))

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

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b01000011}, cpu.dec_Ss_("HL"))

	got := dmaX.GetMemory(0x3572)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x00)
	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b10010010}, cpu.dec_Ss_("HL"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0xff)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x80)
	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b00010110}, cpu.dec_Ss_("HL"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x7f)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestLd_Hl_N(t *testing.T) {
	resetAll()
	cpu.HL = 0x1015
	dmaX.SetMemoryBulk(0x0000, []uint8{0x36, 0x28})

	checkCpu(t, 10, map[string]uint16{"PC": 2, "HL": 0x1015}, cpu.ld_Ss_N("HL"))

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

func TestJrCN(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x38, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b11010111}, cpu.jrCN)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x38, 0x32})

	checkCpu(t, 7, map[string]uint16{"PC": 0x05, "Flags": 0b11010110}, cpu.jrCN)
}

func TestAddHlSp(t *testing.T) {
	resetAll()
	cpu.SP = 0xa76c //  1010 0111 0110 1100
	cpu.HL = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0xa76c, "HL": 0x009f, "Flags": 0b00010001}, cpu.addSsRr("HL", "SP"))

	resetAll()
	cpu.SP = 0x7fff
	cpu.HL = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0x7fff, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("HL", "SP"))
}

func TestAddIxSp(t *testing.T) {
	resetAll()
	cpu.SP = 0xa76c //  1010 0111 0110 1100
	cpu.IX = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "SP": 0xa76c, "IX": 0x009f, "Flags": 0b00010001}, cpu.addSsRr("IX", "SP"))

	resetAll()
	cpu.SP = 0x7fff
	cpu.IX = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "SP": 0x7fff, "IX": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("IX", "SP"))
}

func TestAddIySp(t *testing.T) {
	resetAll()
	cpu.SP = 0xa76c //  1010 0111 0110 1100
	cpu.IY = 0x5933 //  0101 1001 0011 0011
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "SP": 0xa76c, "IY": 0x009f, "Flags": 0b00010001}, cpu.addSsRr("IY", "SP"))

	resetAll()
	cpu.SP = 0x7fff
	cpu.IY = 0x7fff
	cpu.setFlags(0b00000010)

	checkCpu(t, 15, map[string]uint16{"PC": 2, "SP": 0x7fff, "IY": 0xfffe, "Flags": 0b00010000}, cpu.addSsRr("IY", "SP"))
}

func TestLdA_Nn_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x3a, 0x57, 0x12})
	dmaX.SetMemoryByte(0x1257, 0x64)
	cpu.setAcc(0xff)

	checkCpu(t, 13, map[string]uint16{"PC": 3, "A": 0x64}, cpu.ldA_Nn_)
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

func TestLdAN(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "A": 0x64}, cpu.ldAN)
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
	checkCpu(t, 7, map[string]uint16{"PC": 1, "BC": 0x1234, "HL": 0x5678}, cpu.ld_Ss_R("HL", 'B'))

	got := dmaX.GetMemory(0x5678)
	want := uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.BC = 0x1234
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "BC": 0x1234, "HL": 0x5678}, cpu.ld_Ss_R("HL", 'C'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x34)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.DE = 0x1234
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "DE": 0x1234, "HL": 0x5678}, cpu.ld_Ss_R("HL", 'D'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.DE = 0x1234
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "DE": 0x1234, "HL": 0x5678}, cpu.ld_Ss_R("HL", 'E'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x34)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "HL": 0x5678}, cpu.ld_Ss_R("HL", 'H'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x56)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "HL": 0x5678}, cpu.ld_Ss_R("HL", 'L'))

	got = dmaX.GetMemory(0x5678)
	want = uint8(0x78)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.setAcc(0x12)
	cpu.HL = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x12, "HL": 0x5678}, cpu.ld_Ss_R("HL", 'A'))

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

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x00, "BC": 0xa9a9, "DE": 0xa9a9, "HL": 0xa9a9, "Flags": 0b01010100}, cpu.andR(register))

		resetAll()
		if register == 'A' {
			cpu.setAcc(0x97)
		} else {
			cpu.setAcc(0xdf)
		}
		cpu.BC = 0xb7b7
		cpu.DE = 0xb7b7
		cpu.HL = 0xb7b7

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x97, "BC": 0xb7b7, "DE": 0xb7b7, "HL": 0xb7b7, "Flags": 0b10010000}, cpu.andR(register))
	}
}

func TestAnd_Hl_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0xa9)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x00, "HL": 0x1234, "Flags": 0b01010100}, cpu.and_Ss_("HL"))

	resetAll()
	cpu.setAcc(0xdf)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0xb7)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x97, "HL": 0x1234, "Flags": 0b10010000}, cpu.and_Ss_("HL"))
}

func TestAnd_Ix_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	cpu.IX = 0x121b
	dmaX.SetMemoryByte(0x1234, 0xa9)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x00, "IX": 0x121b, "Flags": 0b01010100}, cpu.and_Ss_("IX"))

	resetAll()
	cpu.setAcc(0xdf)
	cpu.IX = 0x121b
	dmaX.SetMemoryByte(0x1234, 0xb7)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x97, "IX": 0x121b, "Flags": 0b10010000}, cpu.and_Ss_("IX"))
}

func TestAnd_Iy_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	cpu.IY = 0x121b
	dmaX.SetMemoryByte(0x1234, 0xa9)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x00, "IY": 0x121b, "Flags": 0b01010100}, cpu.and_Ss_("IY"))

	resetAll()
	cpu.setAcc(0xdf)
	cpu.IY = 0x121b
	dmaX.SetMemoryByte(0x1234, 0xb7)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x97, "IY": 0x121b, "Flags": 0b10010000}, cpu.and_Ss_("IY"))
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

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x00, "BC": 0x5656, "DE": 0x5656, "HL": 0x5656, "Flags": 0b01000100}, cpu.xorR(register))

		if register == 'A' {
			continue
		}

		resetAll()
		cpu.setAcc(0x20)
		cpu.BC = 0xb7b7
		cpu.DE = 0xb7b7
		cpu.HL = 0xb7b7

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x97, "BC": 0xb7b7, "DE": 0xb7b7, "HL": 0xb7b7, "Flags": 0b10000000}, cpu.xorR(register))
	}
}

func TestXor_Hl_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0x56)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x00, "HL": 0x1234, "Flags": 0b01000100}, cpu.xor_Ss_("HL"))

	resetAll()
	cpu.setAcc(0x20)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0xb7)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x97, "HL": 0x1234, "Flags": 0b10000000}, cpu.xor_Ss_("HL"))
}

func TestXor_Ix_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	cpu.IX = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x56)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x00, "IX": 0x121b, "Flags": 0b01000100}, cpu.xor_Ss_("IX"))

	resetAll()
	cpu.setAcc(0x20)
	cpu.IX = 0x121b
	dmaX.SetMemoryByte(0x1234, 0xb7)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x97, "IX": 0x121b, "Flags": 0b10000000}, cpu.xor_Ss_("IX"))
}

func TestXor_Iy_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	cpu.IY = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x56)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x00, "IY": 0x121b, "Flags": 0b01000100}, cpu.xor_Ss_("IY"))

	resetAll()
	cpu.setAcc(0x20)
	cpu.IY = 0x121b
	dmaX.SetMemoryByte(0x1234, 0xb7)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x97, "IY": 0x121b, "Flags": 0b10000000}, cpu.xor_Ss_("IY"))
}

func TestOrR(t *testing.T) {
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		resetAll()
		cpu.setAcc(0x00)
		cpu.BC = 0x0000
		cpu.DE = 0x0000
		cpu.HL = 0x0000

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x00, "BC": 0x0000, "DE": 0x0000, "HL": 0x0000, "Flags": 0b01000100}, cpu.orR(register))

		resetAll()
		if register == 'A' {
			cpu.setAcc(0x97)
		} else {
			cpu.setAcc(0x84)
		}

		cpu.BC = 0x1313
		cpu.DE = 0x1313
		cpu.HL = 0x1313

		checkCpu(t, 4, map[string]uint16{"PC": 1, "A": 0x97, "BC": 0x1313, "DE": 0x1313, "HL": 0x1313, "Flags": 0b10000000}, cpu.orR(register))
	}
}

func TestOr_Hl_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x00)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0x00)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x00, "HL": 0x1234, "Flags": 0b01000100}, cpu.or_Ss_("HL"))

	resetAll()
	cpu.setAcc(0x84)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0x13)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x97, "HL": 0x1234, "Flags": 0b10000000}, cpu.or_Ss_("HL"))
}

func TestOr_Ix_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x00)
	cpu.IX = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x00)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x00, "IX": 0x121b, "Flags": 0b01000100}, cpu.or_Ss_("IX"))

	resetAll()
	cpu.setAcc(0x84)
	cpu.IX = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x13)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x97, "IX": 0x121b, "Flags": 0b10000000}, cpu.or_Ss_("IX"))
}

func TestOr_Iy_(t *testing.T) {
	resetAll()
	cpu.setAcc(0x00)
	cpu.IY = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x00)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x00, "IY": 0x121b, "Flags": 0b01000100}, cpu.or_Ss_("IY"))

	resetAll()
	cpu.setAcc(0x84)
	cpu.IY = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x13)
	dmaX.SetMemoryByte(0x0002, 0x19)

	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x97, "IY": 0x121b, "Flags": 0b10000000}, cpu.or_Ss_("IY"))
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

func TestJpNzNn(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b10010111}, cpu.jpNzNn)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010111}, cpu.jpNzNn)
}

func TestJpNn(t *testing.T) {
	resetAll()
	cpu.PC = 3
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678}, cpu.jpNn)
}

func TestCallNzNn(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b10010111}, cpu.callNzNn)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x37), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010111}, cpu.callNzNn)

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
		wantL, wantH := uint8(0x35), uint8(0x12)

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

func TestJpZNn(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010111}, cpu.jpZNn)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b10010111}, cpu.jpZNn)
}

func TestCallZNn(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.callZNn)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x37), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b10010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b10010111}, cpu.callZNn)

	gotL, gotH = dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH = uint8(0x00), uint8(0x00)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestCallNn(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe}, cpu.callNn)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x37), uint8(0x12)

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

func TestJpNcNn(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010110}, cpu.jpNcNn)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010111}, cpu.jpNcNn)
}

func TestOut_N_A(t *testing.T) {
	resetAll()
	cpu.setAcc(0xaf)
	dmaX.SetMemoryByte(0x0001, 0x45)

	checkCpu(t, 11, map[string]uint16{"PC": 2, "A": 0xaf}, cpu.out_N_A)

	got := cpu.getPort(0x45)
	want := uint8(0xaf)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestCallNcNn(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010110}, cpu.callNcNn)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x37), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010111}, cpu.callNcNn)

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

func TestJpCNn(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010111}, cpu.jpCNn)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010110}, cpu.jpCNn)
}

func TestInA_N_(t *testing.T) {
	resetAll()
	cpu.setPort(0x45, 0xaf)
	dmaX.SetMemoryByte(0x0001, 0x45)

	checkCpu(t, 11, map[string]uint16{"PC": 2, "A": 0xaf}, cpu.inA_N_)
}

func TestCallCNn(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.callCNn)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x37), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010110}, cpu.callCNn)

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

	checkCpu(t, 10, map[string]uint16{"PC": 1, "SP": 0x0000, "HL": 0x5678}, cpu.popSs("HL"))
}

func TestJpPoNn(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010011}, cpu.jpPoNn)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010111}, cpu.jpPoNn)
}

func TestEx_Sp_Hl(t *testing.T) {
	resetAll()
	cpu.HL = 0x7012
	cpu.SP = 0x8856
	dmaX.SetMemoryBulk(0x8856, []uint8{0x11, 0x22})

	checkCpu(t, 19, map[string]uint16{"PC": 1, "HL": 0x2211, "SP": 0x8856}, cpu.ex_Sp_Ss("HL"))

	gotL, gotH := dmaX.GetMemory(0x8856), dmaX.GetMemory(0x8857)
	wantL, wantH := uint8(0x12), uint8(0x70)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

}

func TestCallPoNn(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010011}, cpu.callPoNn)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x37), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010111}, cpu.callPoNn)

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
	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0xfffe, "HL": 0x1234}, cpu.pushSs("HL"))

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestAndN(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	dmaX.SetMemoryByte(0x0001, 0xa9)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x00, "Flags": 0b01010100}, cpu.andN)

	resetAll()
	cpu.setAcc(0xdf)
	dmaX.SetMemoryByte(0x0001, 0xb7)

	checkCpu(t, 7, map[string]uint16{"PC": 1, "A": 0x97, "Flags": 0b10010000}, cpu.andN)
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

	checkCpu(t, 4, map[string]uint16{"PC": 0x5678, "HL": 0x1234}, cpu.jp_Ss_("HL"))
}

func TestJp_Ix_(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.IX = 0x1234
	dmaX.SetMemoryBulk(0x1234, []uint8{0x78, 0x56})

	checkCpu(t, 8, map[string]uint16{"PC": 0x5678, "IX": 0x1234}, cpu.jp_Ss_("IX"))
}

func TestJp_Iy_(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.IY = 0x1234
	dmaX.SetMemoryBulk(0x1234, []uint8{0x78, 0x56})

	checkCpu(t, 8, map[string]uint16{"PC": 0x5678, "IY": 0x1234}, cpu.jp_Ss_("IY"))
}

func TestJpPeNn(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010111}, cpu.jpPeNn)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010011}, cpu.jpPeNn)
}

func TestExDeHl(t *testing.T) {
	resetAll()
	cpu.DE = 0x2822
	cpu.HL = 0x499a

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x499a, "HL": 0x2822}, cpu.exDeSs("HL"))
}

func TestExDeIx(t *testing.T) {
	resetAll()
	cpu.DE = 0x2822
	cpu.IX = 0x499a

	checkCpu(t, 8, map[string]uint16{"PC": 2, "DE": 0x499a, "IX": 0x2822}, cpu.exDeSs("IX"))
}

func TestExDeIy(t *testing.T) {
	resetAll()
	cpu.DE = 0x2822
	cpu.IY = 0x499a

	checkCpu(t, 8, map[string]uint16{"PC": 2, "DE": 0x499a, "IY": 0x2822}, cpu.exDeSs("IY"))
}

func TestCallPeNn(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.callPeNn)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x37), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010011}, cpu.callPeNn)

	gotL, gotH = dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH = uint8(0x00), uint8(0x00)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestXorN(t *testing.T) {
	resetAll()
	cpu.setAcc(0x56)
	dmaX.SetMemoryByte(0x0001, 0x56)

	checkCpu(t, 7, map[string]uint16{"PC": 2, "A": 0x00, "Flags": 0b01000100}, cpu.xorN)

	resetAll()
	cpu.setAcc(0x20)
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x0001, 0xb7)

	checkCpu(t, 7, map[string]uint16{"PC": 2, "A": 0x97, "Flags": 0b10000000}, cpu.xorN)
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

func TestJpPNn(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b01010111}, cpu.jpPNn)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b11010111}, cpu.jpPNn)
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

func TestCallPNn(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b01010111}, cpu.callPNn)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x37), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b11010111}, cpu.callPNn)

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

func TestOrN(t *testing.T) {
	resetAll()
	cpu.setAcc(0x00)
	dmaX.SetMemoryByte(0x0001, 0x00)

	checkCpu(t, 7, map[string]uint16{"PC": 2, "A": 0x00, "Flags": 0b01000100}, cpu.orN)

	resetAll()
	cpu.setAcc(0x84)
	dmaX.SetMemoryByte(0x0001, 0x13)

	checkCpu(t, 7, map[string]uint16{"PC": 2, "A": 0x97, "Flags": 0b10000000}, cpu.orN)
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

	checkCpu(t, 6, map[string]uint16{"PC": 1, "SP": 0x442e, "HL": 0x442e}, cpu.ldSpSs("HL"))
}

func TestJpMNn(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x5678, "Flags": 0b11010111}, cpu.jpZNn)

	resetAll()
	cpu.PC = 3
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryBulk(0x0004, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x06, "Flags": 0b01010111}, cpu.jpMNn)
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

func TestCallMNn(t *testing.T) {
	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 17, map[string]uint16{"PC": 0x5678, "SP": 0xfffe, "Flags": 0b11010111}, cpu.callMNn)

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x37), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}

	resetAll()
	cpu.PC = 0x1234
	cpu.SP = 0x0000
	cpu.setFlags(0b01010111)
	dmaX.SetMemoryBulk(0x1235, []uint8{0x78, 0x56})

	checkCpu(t, 10, map[string]uint16{"PC": 0x1237, "SP": 0x0000, "Flags": 0b01010111}, cpu.callMNn)

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

func TestLd_Nn_Rr(t *testing.T) {
	for _, registerPair := range [4]string{"BC", "DE", "HL", "SP"} {
		resetAll()
		cpu.BC = 0x4644
		cpu.DE = 0x4644
		cpu.HL = 0x4644
		cpu.SP = 0x4644
		dmaX.SetMemoryBulk(0x0002, []uint8{0x20, 0x10})

		checkCpu(t, 20, map[string]uint16{"PC": 4, "BC": 0x4644, "DE": 0x4644, "HL": 0x4644, "SP": 0x4644}, cpu.ld_Nn_Rr(registerPair))

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

func TestLdRr_Nn_(t *testing.T) {
	for _, registerPair := range [4]string{"BC", "DE", "HL", "SP"} {
		resetAll()
		cpu.BC = 0x0123
		cpu.DE = 0x4567
		cpu.HL = 0x89ab
		cpu.SP = 0xcdef
		dmaX.SetMemoryBulk(0x0002, []uint8{0x20, 0x10})
		dmaX.SetMemoryBulk(0x1020, []uint8{0x85, 0x24})

		checkCpu(t, 20, map[string]uint16{"PC": 4, registerPair: 0x2485}, cpu.ldRr_Nn_(registerPair))
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

func TestLd_Ix_R(t *testing.T) {
	resetAll()
	cpu.BC = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "BC": 0x1234, "IX": 0x5678}, cpu.ld_Ss_R("IX", 'B'))

	got := dmaX.GetMemory(0x5691)
	want := uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.BC = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "BC": 0x1234, "IX": 0x5678}, cpu.ld_Ss_R("IX", 'C'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x34)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.DE = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "DE": 0x1234, "IX": 0x5678}, cpu.ld_Ss_R("IX", 'D'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.DE = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "DE": 0x1234, "IX": 0x5678}, cpu.ld_Ss_R("IX", 'E'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x34)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "HL": 0x1234, "IX": 0x5678}, cpu.ld_Ss_R("IX", 'H'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "HL": 0x1234, "IX": 0x5678}, cpu.ld_Ss_R("IX", 'L'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x34)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.setAcc(0x12)
	cpu.IX = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x12, "IX": 0x5678}, cpu.ld_Ss_R("IX", 'A'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}
}

func TestLd_Iy_R(t *testing.T) {
	resetAll()
	cpu.BC = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "BC": 0x1234, "IY": 0x5678}, cpu.ld_Ss_R("IY", 'B'))

	got := dmaX.GetMemory(0x5691)
	want := uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.BC = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "BC": 0x1234, "IY": 0x5678}, cpu.ld_Ss_R("IY", 'C'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x34)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.DE = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "DE": 0x1234, "IY": 0x5678}, cpu.ld_Ss_R("IY", 'D'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.DE = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "DE": 0x1234, "IY": 0x5678}, cpu.ld_Ss_R("IY", 'E'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x34)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "HL": 0x1234, "IY": 0x5678}, cpu.ld_Ss_R("IY", 'H'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "HL": 0x1234, "IY": 0x5678}, cpu.ld_Ss_R("IY", 'L'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x34)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}

	resetAll()
	cpu.setAcc(0x12)
	cpu.IY = 0x5678
	dmaX.SetMemoryByte(0x5678, 0xab)
	dmaX.SetMemoryByte(0x0002, 0x19)
	checkCpu(t, 19, map[string]uint16{"PC": 3, "A": 0x12, "IY": 0x5678}, cpu.ld_Ss_R("IY", 'A'))

	got = dmaX.GetMemory(0x5691)
	want = uint8(0x12)

	if got != want {
		t.Errorf("got 0x%02x, want %02x", got, want)
	}
}

func TestLd_Ix_N(t *testing.T) {
	resetAll()
	cpu.IX = 0x1015
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0x36, 0x19, 0x28})

	checkCpu(t, 19, map[string]uint16{"PC": 4, "IX": 0x1015}, cpu.ld_Ss_N("IX"))

	got := dmaX.GetMemory(0x102e)
	want := uint8(0x28)
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestLd_Iy_N(t *testing.T) {
	resetAll()
	cpu.IY = 0x1015
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0x36, 0x19, 0x28})

	checkCpu(t, 19, map[string]uint16{"PC": 4, "IY": 0x1015}, cpu.ld_Ss_N("IY"))

	got := dmaX.GetMemory(0x102e)
	want := uint8(0x28)
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestLdIxNn(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0x21, 0x64, 0x32})

	checkCpu(t, 14, map[string]uint16{"PC": 4, "IX": 0x3264}, cpu.ldSsNn("IX"))
}

func TestLdIyNn(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0x21, 0x64, 0x32})

	checkCpu(t, 14, map[string]uint16{"PC": 4, "IY": 0x3264}, cpu.ldSsNn("IY"))
}

func TestLdIx_Nn_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0x2a, 0x29, 0xb2})
	dmaX.SetMemoryBulk(0xb229, []uint8{0x37, 0xa1})

	checkCpu(t, 20, map[string]uint16{"PC": 4, "IX": 0xa137}, cpu.ldSs_Nn_("IX"))
}

func TestLdIy_Nn_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0x2a, 0x29, 0xb2})
	dmaX.SetMemoryBulk(0xb229, []uint8{0x37, 0xa1})

	checkCpu(t, 20, map[string]uint16{"PC": 4, "IY": 0xa137}, cpu.ldSs_Nn_("IY"))
}

func TestLd_Nn_Ix(t *testing.T) {
	resetAll()
	cpu.IX = 0x483a
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0x22, 0x29, 0xb2})

	checkCpu(t, 20, map[string]uint16{"PC": 4, "IX": 0x483a}, cpu.ld_Nn_Ss("IX"))

	gotH, gotL := dmaX.GetMemory(0xb229), dmaX.GetMemory(0xb22a)
	wantH, wantL := uint8(0x3a), uint8(0x48)

	if gotH != wantH || gotL != wantL {
		t.Errorf("got 0x%x%x, want 0x%x%x", gotH, gotL, wantH, wantL)
	}
}

func TestLd_Nn_Iy(t *testing.T) {
	resetAll()
	cpu.IY = 0x483a
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0x22, 0x29, 0xb2})

	checkCpu(t, 20, map[string]uint16{"PC": 4, "IY": 0x483a}, cpu.ld_Nn_Ss("IY"))

	gotH, gotL := dmaX.GetMemory(0xb229), dmaX.GetMemory(0xb22a)
	wantH, wantL := uint8(0x3a), uint8(0x48)

	if gotH != wantH || gotL != wantL {
		t.Errorf("got 0x%x%x, want 0x%x%x", gotH, gotL, wantH, wantL)
	}
}

func TestLdSpIx(t *testing.T) {
	resetAll()
	cpu.SP = 0xfffc
	cpu.IX = 0x442e

	checkCpu(t, 10, map[string]uint16{"PC": 2, "SP": 0x442e, "IX": 0x442e}, cpu.ldSpSs("IX"))
}

func TestLdSpIy(t *testing.T) {
	resetAll()
	cpu.SP = 0xfffc
	cpu.IY = 0x442e

	checkCpu(t, 10, map[string]uint16{"PC": 2, "SP": 0x442e, "IY": 0x442e}, cpu.ldSpSs("IY"))
}

func TestPushIx(t *testing.T) {
	resetAll()
	cpu.IX = 0x1234
	cpu.SP = 0x0000
	checkCpu(t, 15, map[string]uint16{"PC": 2, "SP": 0xfffe, "IX": 0x1234}, cpu.pushSs("IX"))

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestPushIy(t *testing.T) {
	resetAll()
	cpu.IY = 0x1234
	cpu.SP = 0x0000
	checkCpu(t, 15, map[string]uint16{"PC": 2, "SP": 0xfffe, "IY": 0x1234}, cpu.pushSs("IY"))

	gotL, gotH := dmaX.GetMemory(0xfffe), dmaX.GetMemory(0xffff)
	wantL, wantH := uint8(0x34), uint8(0x12)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestPopIx(t *testing.T) {
	resetAll()
	cpu.SP = 0xfffe
	cpu.IX = 0x1234
	dmaX.SetMemoryBulk(0xfffe, []uint8{0x78, 0x56})

	checkCpu(t, 14, map[string]uint16{"PC": 2, "SP": 0x0000, "IX": 0x5678}, cpu.popSs("IX"))
}

func TestPopIy(t *testing.T) {
	resetAll()
	cpu.SP = 0xfffe
	cpu.IY = 0x1234
	dmaX.SetMemoryBulk(0xfffe, []uint8{0x78, 0x56})

	checkCpu(t, 14, map[string]uint16{"PC": 2, "SP": 0x0000, "IY": 0x5678}, cpu.popSs("IY"))
}

func TestEx_Sp_Ix(t *testing.T) {
	resetAll()
	cpu.IX = 0x7012
	cpu.SP = 0x8856
	dmaX.SetMemoryBulk(0x8856, []uint8{0x11, 0x22})

	checkCpu(t, 23, map[string]uint16{"PC": 2, "IX": 0x2211, "SP": 0x8856}, cpu.ex_Sp_Ss("IX"))

	gotL, gotH := dmaX.GetMemory(0x8856), dmaX.GetMemory(0x8857)
	wantL, wantH := uint8(0x12), uint8(0x70)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestEx_Sp_Iy(t *testing.T) {
	resetAll()
	cpu.IY = 0x7012
	cpu.SP = 0x8856
	dmaX.SetMemoryBulk(0x8856, []uint8{0x11, 0x22})

	checkCpu(t, 23, map[string]uint16{"PC": 2, "IY": 0x2211, "SP": 0x8856}, cpu.ex_Sp_Ss("IY"))

	gotL, gotH := dmaX.GetMemory(0x8856), dmaX.GetMemory(0x8857)
	wantL, wantH := uint8(0x12), uint8(0x70)

	if gotL != wantL || gotH != wantH {
		t.Errorf("got 0x%02x%02x, want 0x%02x%02x", gotH, gotL, wantH, wantL)
	}
}

func TestInc_Ix_(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010111)
	cpu.IX = 0x353f
	dmaX.SetMemoryByte(0x0002, 0x33)
	dmaX.SetMemoryByte(0x3572, 0x25)

	checkCpu(t, 23, map[string]uint16{"PC": 3, "IX": 0x353f, "Flags": 0b00000001}, cpu.inc_Ss_("IX"))

	got := dmaX.GetMemory(0x3572)
	want := uint8(0x26)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.IX = 0x353f
	dmaX.SetMemoryByte(0x0002, 0x33)
	dmaX.SetMemoryByte(0x3572, 0xff)
	checkCpu(t, 23, map[string]uint16{"PC": 3, "IX": 0x353f, "Flags": 0b01010000}, cpu.inc_Ss_("IX"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.IX = 0x353f
	dmaX.SetMemoryByte(0x0002, 0x33)
	dmaX.SetMemoryByte(0x3572, 0x7f)
	checkCpu(t, 23, map[string]uint16{"PC": 3, "IX": 0x353f, "Flags": 0b10010100}, cpu.inc_Ss_("IX"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x80)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestInc_Iy_(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010111)
	cpu.IY = 0x353f
	dmaX.SetMemoryByte(0x0002, 0x33)
	dmaX.SetMemoryByte(0x3572, 0x25)

	checkCpu(t, 23, map[string]uint16{"PC": 3, "IY": 0x353f, "Flags": 0b00000001}, cpu.inc_Ss_("IY"))

	got := dmaX.GetMemory(0x3572)
	want := uint8(0x26)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b10000110)
	cpu.IY = 0x353f
	dmaX.SetMemoryByte(0x0002, 0x33)
	dmaX.SetMemoryByte(0x3572, 0xff)
	checkCpu(t, 23, map[string]uint16{"PC": 3, "IY": 0x353f, "Flags": 0b01010000}, cpu.inc_Ss_("IY"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b01000010)
	cpu.IY = 0x353f
	dmaX.SetMemoryByte(0x0002, 0x33)
	dmaX.SetMemoryByte(0x3572, 0x7f)
	checkCpu(t, 23, map[string]uint16{"PC": 3, "IY": 0x353f, "Flags": 0b10010100}, cpu.inc_Ss_("IY"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x80)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestDec_Ix_(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010101)
	cpu.IX = 0x353f
	dmaX.SetMemoryByte(0x3572, 0x01)
	dmaX.SetMemoryByte(0x0002, 0x33)

	checkCpu(t, 23, map[string]uint16{"PC": 3, "IX": 0x353f, "Flags": 0b01000011}, cpu.dec_Ss_("IX"))

	got := dmaX.GetMemory(0x3572)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.IX = 0x353f
	dmaX.SetMemoryByte(0x3572, 0x00)
	dmaX.SetMemoryByte(0x0002, 0x33)

	checkCpu(t, 23, map[string]uint16{"PC": 3, "IX": 0x353f, "Flags": 0b10010010}, cpu.dec_Ss_("IX"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0xff)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.IX = 0x353f
	dmaX.SetMemoryByte(0x3572, 0x80)
	dmaX.SetMemoryByte(0x0002, 0x33)

	checkCpu(t, 23, map[string]uint16{"PC": 3, "IX": 0x353f, "Flags": 0b00010110}, cpu.dec_Ss_("IX"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x7f)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestDec_Iy_(t *testing.T) {
	resetAll()
	cpu.setFlags(0b11010101)
	cpu.IY = 0x353f
	dmaX.SetMemoryByte(0x3572, 0x01)
	dmaX.SetMemoryByte(0x0002, 0x33)

	checkCpu(t, 23, map[string]uint16{"PC": 3, "IY": 0x353f, "Flags": 0b01000011}, cpu.dec_Ss_("IY"))

	got := dmaX.GetMemory(0x3572)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b01000100)
	cpu.IY = 0x353f
	dmaX.SetMemoryByte(0x3572, 0x00)
	dmaX.SetMemoryByte(0x0002, 0x33)

	checkCpu(t, 23, map[string]uint16{"PC": 3, "IY": 0x353f, "Flags": 0b10010010}, cpu.dec_Ss_("IY"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0xff)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.setFlags(0b11000000)
	cpu.IY = 0x353f
	dmaX.SetMemoryByte(0x3572, 0x80)
	dmaX.SetMemoryByte(0x0002, 0x33)

	checkCpu(t, 23, map[string]uint16{"PC": 3, "IY": 0x353f, "Flags": 0b00010110}, cpu.dec_Ss_("IY"))

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x7f)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestLdi(t *testing.T) {
	resetAll()
	cpu.HL = 0x1111
	cpu.DE = 0x2222
	cpu.BC = 0x0007
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryByte(0x1111, 0x88)
	dmaX.SetMemoryByte(0x2222, 0x66)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x1112, "DE": 0x2223, "BC": 0x0006, "Flags": 0b11000101}, cpu.ldi)

	got := dmaX.GetMemory(0x1111)
	want := uint8(0x88)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}

	got = dmaX.GetMemory(0x2222)
	want = uint8(0x88)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestCpi(t *testing.T) {
	resetAll()
	cpu.setAcc(0x3b)
	cpu.HL = 0x1111
	cpu.BC = 0x0001
	cpu.setFlags(0b01010001)
	dmaX.SetMemoryByte(0x1111, 0x3b)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "A": 0x3b, "HL": 0x1112, "BC": 0x0000, "Flags": 0b01000011}, cpu.cpi)

	got := dmaX.GetMemory(0x1111)
	want := uint8(0x3b)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}

	resetAll()
	cpu.setAcc(0x00)
	cpu.HL = 0x1111
	cpu.BC = 0x8000
	cpu.setFlags(0b01000000)
	dmaX.SetMemoryByte(0x1111, 0x7f)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "A": 0x00, "HL": 0x1112, "BC": 0x7fff, "Flags": 0b10010110}, cpu.cpi)

	got = dmaX.GetMemory(0x1111)
	want = uint8(0x7f)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestIni(t *testing.T) {
	resetAll()
	cpu.BC = 0x1007
	cpu.HL = 0x1000
	cpu.States.Ports[0x07] = 0x7b
	cpu.setFlags(0b01000000)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x1001, "BC": 0x0f07, "Flags": 0b00000010}, cpu.ini)

	got := dmaX.GetMemory(0x1000)
	want := uint8(0x7b)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}

	resetAll()
	cpu.BC = 0x0107
	cpu.HL = 0x1000
	cpu.States.Ports[0x07] = 0x7b
	cpu.setFlags(0b10010101)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x1001, "BC": 0x0007, "Flags": 0b11010111}, cpu.ini)

	got = dmaX.GetMemory(0x1000)
	want = uint8(0x7b)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestOuti(t *testing.T) {
	resetAll()
	cpu.BC = 0x1007
	cpu.HL = 0x1000
	cpu.setFlags(0b01000000)
	dmaX.SetMemoryByte(0x1000, 0x59)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x1001, "BC": 0x0f07, "Flags": 0b00000010}, cpu.outi)

	got := cpu.getPort(0x07)
	want := uint8(0x59)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}

	resetAll()
	cpu.BC = 0x0107
	cpu.HL = 0x1000
	cpu.setFlags(0b10010101)
	dmaX.SetMemoryByte(0x1000, 0x59)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x1001, "BC": 0x0007, "Flags": 0b11010111}, cpu.outi)

	got = cpu.getPort(0x07)
	want = uint8(0x59)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestLdd(t *testing.T) {
	resetAll()
	cpu.HL = 0x1111
	cpu.DE = 0x2222
	cpu.BC = 0x0007
	cpu.setFlags(0b11010011)
	dmaX.SetMemoryByte(0x1111, 0x88)
	dmaX.SetMemoryByte(0x2222, 0x66)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x1110, "DE": 0x2221, "BC": 0x0006, "Flags": 0b11000101}, cpu.ldd)

	got := dmaX.GetMemory(0x1111)
	want := uint8(0x88)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}

	got = dmaX.GetMemory(0x2222)
	want = uint8(0x88)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestCpd(t *testing.T) {
	resetAll()
	cpu.setAcc(0x3b)
	cpu.HL = 0x1111
	cpu.BC = 0x0001
	cpu.setFlags(0b01010001)
	dmaX.SetMemoryByte(0x1111, 0x3b)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "A": 0x3b, "HL": 0x1110, "BC": 0x0000, "Flags": 0b01000011}, cpu.cpd)

	got := dmaX.GetMemory(0x1111)
	want := uint8(0x3b)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}

	resetAll()
	cpu.setAcc(0x00)
	cpu.HL = 0x1111
	cpu.BC = 0x8000
	cpu.setFlags(0b01000000)
	dmaX.SetMemoryByte(0x1111, 0x7f)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "A": 0x00, "HL": 0x1110, "BC": 0x7fff, "Flags": 0b10010110}, cpu.cpd)

	got = dmaX.GetMemory(0x1111)
	want = uint8(0x7f)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestInd(t *testing.T) {
	resetAll()
	cpu.BC = 0x1007
	cpu.HL = 0x1000
	cpu.setPort(0x07, 0x7b)
	cpu.setFlags(0b01000000)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x0fff, "BC": 0x0f07, "Flags": 0b00000010}, cpu.ind)

	got := dmaX.GetMemory(0x1000)
	want := uint8(0x7b)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}

	resetAll()
	cpu.BC = 0x0107
	cpu.HL = 0x1000
	cpu.States.Ports[0x07] = 0x7b
	cpu.setFlags(0b10010101)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x0fff, "BC": 0x0007, "Flags": 0b11010111}, cpu.ind)

	got = dmaX.GetMemory(0x1000)
	want = uint8(0x7b)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestOutd(t *testing.T) {
	resetAll()
	cpu.BC = 0x1007
	cpu.HL = 0x1000
	cpu.setFlags(0b01000000)
	dmaX.SetMemoryByte(0x1000, 0x59)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x0fff, "BC": 0x0f07, "Flags": 0b00000010}, cpu.outd)

	got := cpu.getPort(0x07)
	want := uint8(0x59)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}

	resetAll()
	cpu.BC = 0x0107
	cpu.HL = 0x1000
	cpu.setFlags(0b10010101)
	dmaX.SetMemoryByte(0x1000, 0x59)

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x0fff, "BC": 0x0007, "Flags": 0b11010111}, cpu.outd)

	got = cpu.getPort(0x07)
	want = uint8(0x59)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestLdir(t *testing.T) {
	resetAll()
	cpu.HL = 0x1111
	cpu.DE = 0x2222
	cpu.BC = 0x0003
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1111, []uint8{0x88, 0x36, 0xa5})
	dmaX.SetMemoryBulk(0x2222, []uint8{0x66, 0x59, 0xc5})

	for cpu.BC > 1 {
		got := cpu.ldir()
		want := uint8(21)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x1114, "DE": 0x2225, "BC": 0x0000, "Flags": 0b11000001}, cpu.ldir)

	gotHLA, gotHLB, gotHLC := dmaX.GetMemory(0x1111), dmaX.GetMemory(0x1112), dmaX.GetMemory(0x1113)
	wantHLA, wantHLB, wantHLC := uint8(0x88), uint8(0x36), uint8(0xa5)

	if gotHLA != wantHLA || gotHLB != wantHLB || gotHLC != wantHLC {
		t.Errorf("got %02x%02x%02x, want %02x%02x%02x", gotHLA, gotHLB, gotHLC, wantHLA, wantHLB, wantHLC)
	}

	gotDEA, gotDEB, gotDEC := dmaX.GetMemory(0x2222), dmaX.GetMemory(0x2223), dmaX.GetMemory(0x2224)
	wantDEA, wantDEB, wantDEC := uint8(0x88), uint8(0x36), uint8(0xa5)

	if gotDEA != wantDEA || gotDEB != wantDEB || gotDEC != wantDEC {
		t.Errorf("got %02x%02x%02x, want %02x%02x%02x", gotDEA, gotDEB, gotDEC, wantDEA, wantDEB, wantDEC)
	}
}

func TestCpir(t *testing.T) {
	resetAll()
	cpu.HL = 0x1111
	cpu.BC = 0x0007
	cpu.setAcc(0xf3)
	cpu.setFlags(0b11010001)
	dmaX.SetMemoryBulk(0x1111, []uint8{0x52, 0x00, 0xf3})

	for i := 0; i < 2; i++ {
		got := cpu.cpir()
		want := uint8(21)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}

	checkCpu(t, 16, map[string]uint16{"PC": 2, "A": 0xf3, "HL": 0x1114, "BC": 0x0004, "Flags": 0b01000111}, cpu.cpir)

	gotHLA, gotHLB, gotHLC := dmaX.GetMemory(0x1111), dmaX.GetMemory(0x1112), dmaX.GetMemory(0x1113)
	wantHLA, wantHLB, wantHLC := uint8(0x52), uint8(0x00), uint8(0xf3)

	if gotHLA != wantHLA || gotHLB != wantHLB || gotHLC != wantHLC {
		t.Errorf("got %02x%02x%02x, want %02x%02x%02x", gotHLA, gotHLB, gotHLC, wantHLA, wantHLB, wantHLC)
	}
}

func TestInir(t *testing.T) {
	resetAll()
	cpu.BC = 0x0307
	cpu.HL = 0x1000
	ports := []uint8{0x51, 0xa9, 0x03}
	cpu.setFlags(0b01000000)

	for i := 0; cpu.BC > 512; i++ {
		cpu.setPort(0x07, ports[i])
		got := cpu.inir()
		want := uint8(21)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}

	cpu.setPort(0x07, ports[2])
	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x1003, "BC": 0x0007, "Flags": 0b01000010}, cpu.inir)

	gotA, gotB, gotC := dmaX.GetMemory(0x1000), dmaX.GetMemory(0x1001), dmaX.GetMemory(0x1002)
	wantA, wantB, wantC := uint8(0x51), uint8(0xa9), uint8(0x03)

	if gotA != wantA || gotB != wantB || gotC != wantC {
		t.Errorf("got %02x/%02x/%02x, want %02x/%02x/%02x", gotA, gotB, gotC, wantA, wantB, wantC)
	}
}

func TestOtir(t *testing.T) {
	resetAll()
	cpu.BC = 0x0307
	cpu.HL = 0x1000
	dmaX.SetMemoryBulk(0x1000, []uint8{0x51, 0xa9, 0x03})
	cpu.setFlags(0b01000000)

	for i := 0x1000; cpu.BC > 512; i++ {
		gotT, gotPort := cpu.otir(), cpu.getPort(0x07)
		wantT, wantPort := uint8(21), uint8(dmaX.GetMemory(uint16(i)))

		if gotT != wantT || gotPort != wantPort {
			t.Errorf("got %02x (%d), want %02x (%d)", gotPort, gotT, wantPort, wantT)
		}
	}

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x1003, "BC": 0x0007, "Flags": 0b01000010}, cpu.otir)

	got := cpu.getPort(0x07)
	want := uint8(0x03)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestLddr(t *testing.T) {
	resetAll()
	cpu.HL = 0x1114
	cpu.DE = 0x2225
	cpu.BC = 0x0003
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x1112, []uint8{0x88, 0x36, 0xa5})
	dmaX.SetMemoryBulk(0x2223, []uint8{0x66, 0x59, 0xc5})

	for cpu.BC > 1 {
		got := cpu.lddr()
		want := uint8(21)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x1111, "DE": 0x2222, "BC": 0x0000, "Flags": 0b11000001}, cpu.lddr)

	gotHLA, gotHLB, gotHLC := dmaX.GetMemory(0x1112), dmaX.GetMemory(0x1113), dmaX.GetMemory(0x1114)
	wantHLA, wantHLB, wantHLC := uint8(0x88), uint8(0x36), uint8(0xa5)

	if gotHLA != wantHLA || gotHLB != wantHLB || gotHLC != wantHLC {
		t.Errorf("got %02x%02x%02x, want %02x%02x%02x", gotHLA, gotHLB, gotHLC, wantHLA, wantHLB, wantHLC)
	}

	gotDEA, gotDEB, gotDEC := dmaX.GetMemory(0x2223), dmaX.GetMemory(0x2224), dmaX.GetMemory(0x2225)
	wantDEA, wantDEB, wantDEC := uint8(0x88), uint8(0x36), uint8(0xa5)

	if gotDEA != wantDEA || gotDEB != wantDEB || gotDEC != wantDEC {
		t.Errorf("got %02x%02x%02x, want %02x%02x%02x", gotDEA, gotDEB, gotDEC, wantDEA, wantDEB, wantDEC)
	}
}

func TestCpdr(t *testing.T) {
	resetAll()
	cpu.HL = 0x1118
	cpu.BC = 0x0007
	cpu.setAcc(0xf3)
	cpu.setFlags(0b11010001)
	dmaX.SetMemoryBulk(0x1116, []uint8{0xf3, 0x00, 0x52})

	for i := 0; i < 2; i++ {
		got := cpu.cpdr()
		want := uint8(21)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}

	checkCpu(t, 16, map[string]uint16{"PC": 2, "A": 0xf3, "HL": 0x1115, "BC": 0x0004, "Flags": 0b01000111}, cpu.cpdr)

	gotHLA, gotHLB, gotHLC := dmaX.GetMemory(0x1116), dmaX.GetMemory(0x1117), dmaX.GetMemory(0x1118)
	wantHLA, wantHLB, wantHLC := uint8(0xf3), uint8(0x00), uint8(0x52)

	if gotHLA != wantHLA || gotHLB != wantHLB || gotHLC != wantHLC {
		t.Errorf("got %02x%02x%02x, want %02x%02x%02x", gotHLA, gotHLB, gotHLC, wantHLA, wantHLB, wantHLC)
	}
}

func TestIndr(t *testing.T) {
	resetAll()
	cpu.BC = 0x0307
	cpu.HL = 0x1000
	ports := []uint8{0x51, 0xa9, 0x03}
	cpu.setFlags(0b01000000)

	for i := 0; cpu.BC > 512; i++ {
		cpu.setPort(0x07, ports[i])
		got := cpu.indr()
		want := uint8(21)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}

	cpu.setPort(0x07, ports[2])
	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x0ffd, "BC": 0x0007, "Flags": 0b01000010}, cpu.indr)

	gotA, gotB, gotC := dmaX.GetMemory(0x0ffe), dmaX.GetMemory(0x0fff), dmaX.GetMemory(0x1000)
	wantA, wantB, wantC := uint8(0x03), uint8(0xa9), uint8(0x51)

	if gotA != wantA || gotB != wantB || gotC != wantC {
		t.Errorf("got %02x/%02x/%02x, want %02x/%02x/%02x", gotA, gotB, gotC, wantA, wantB, wantC)
	}
}

func TestOtdr(t *testing.T) {
	resetAll()
	cpu.BC = 0x0307
	cpu.HL = 0x1000
	dmaX.SetMemoryBulk(0x0ffe, []uint8{0x51, 0xa9, 0x03})
	cpu.setFlags(0b01000000)

	for i := 0x1000; cpu.BC > 512; i-- {
		gotT, gotPort := cpu.otdr(), cpu.getPort(0x07)
		wantT, wantPort := uint8(21), uint8(dmaX.GetMemory(uint16(i)))

		if gotT != wantT || gotPort != wantPort {
			t.Errorf("got %02x (%d), want %02x (%d)", gotPort, gotT, wantPort, wantT)
		}
	}

	checkCpu(t, 16, map[string]uint16{"PC": 2, "HL": 0x0ffd, "BC": 0x0007, "Flags": 0b01000010}, cpu.otdr)

	got := cpu.getPort(0x07)
	want := uint8(0x51)

	if got != want {
		t.Errorf("got %02x, want %02x", got, want)
	}
}

func TestRlcR(t *testing.T) {
	expectedRegisterMap := map[byte]string{
		'B': "BC", 'C': "BC", 'D': "DE", 'E': "DE", 'H': "HL", 'L': "HL", 'A': "A",
	}
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		expectedValueMap := map[byte]uint16{
			'B': 0x198c, 'C': 0x8c19, 'D': 0x198c, 'E': 0x8c19, 'H': 0x198c, 'L': 0x8c19,
		}

		resetAll()
		cpu.setAcc(0x8c)
		cpu.BC = 0x8c8c
		cpu.DE = 0x8c8c
		cpu.HL = 0x8c8c
		cpu.setFlags(0b11010110)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x19, "Flags": 0b00000001}, cpu.rlcR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b00000001}, cpu.rlcR(register))
		}

		expectedValueMap = map[byte]uint16{
			'B': 0x9a4d, 'C': 0x4d9a, 'D': 0x9a4d, 'E': 0x4d9a, 'H': 0x9a4d, 'L': 0x4d9a,
		}

		resetAll()
		cpu.setAcc(0x4d)
		cpu.BC = 0x4d4d
		cpu.DE = 0x4d4d
		cpu.HL = 0x4d4d
		cpu.setFlags(0b11010111)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x9a, "Flags": 0b10000100}, cpu.rlcR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b10000100}, cpu.rlcR(register))
		}
	}
}

func TestRlcHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8c)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b00000001}, cpu.rlcSs("HL"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x19)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0x4d)
	cpu.setFlags(0b11010111)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b10000100}, cpu.rlcSs("HL"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x9a)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRlcIx(t *testing.T) {
	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b00000001}, cpu.rlcSs("IX"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x19)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IX = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x4d)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	cpu.setFlags(0b11010111)
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b10000100}, cpu.rlcSs("IX"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x9a)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRlcIy(t *testing.T) {
	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b00000001}, cpu.rlcSs("IY"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x19)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IY = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x4d)
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b10000100}, cpu.rlcSs("IY"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x9a)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRrcR(t *testing.T) {
	expectedRegisterMap := map[byte]string{
		'B': "BC", 'C': "BC", 'D': "DE", 'E': "DE", 'H': "HL", 'L': "HL", 'A': "A",
	}
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		expectedValueMap := map[byte]uint16{
			'B': 0xc68d, 'C': 0x8dc6, 'D': 0xc68d, 'E': 0x8dc6, 'H': 0xc68d, 'L': 0x8dc6,
		}

		resetAll()
		cpu.setAcc(0x8d)
		cpu.BC = 0x8d8d
		cpu.DE = 0x8d8d
		cpu.HL = 0x8d8d
		cpu.setFlags(0b11010110)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0xc6, "Flags": 0b10000101}, cpu.rrcR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b10000101}, cpu.rrcR(register))
		}

		expectedValueMap = map[byte]uint16{
			'B': 0x264c, 'C': 0x4c26, 'D': 0x264c, 'E': 0x4c26, 'H': 0x264c, 'L': 0x4c26,
		}

		resetAll()
		cpu.setAcc(0x4c)
		cpu.BC = 0x4c4c
		cpu.DE = 0x4c4c
		cpu.HL = 0x4c4c
		cpu.setFlags(0b11010111)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x26, "Flags": 0b00000000}, cpu.rrcR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b00000000}, cpu.rrcR(register))
		}
	}
}

func TestRrcHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8d)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b10000101}, cpu.rrcSs("HL"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0xc6)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0x4c)
	cpu.setFlags(0b11010111)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b00000000}, cpu.rrcSs("HL"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x26)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRrcIx(t *testing.T) {
	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8d)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b10000101}, cpu.rrcSs("IX"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0xc6)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IX = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x4c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	cpu.setFlags(0b11010111)
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b00000000}, cpu.rrcSs("IX"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x26)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRrcIy(t *testing.T) {
	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8d)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b10000101}, cpu.rrcSs("IY"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0xc6)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IY = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x4c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	cpu.setFlags(0b11010111)
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b00000000}, cpu.rrcSs("IY"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x26)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRlR(t *testing.T) {
	expectedRegisterMap := map[byte]string{
		'B': "BC", 'C': "BC", 'D': "DE", 'E': "DE", 'H': "HL", 'L': "HL", 'A': "A",
	}
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		expectedValueMap := map[byte]uint16{
			'B': 0x188c, 'C': 0x8c18, 'D': 0x188c, 'E': 0x8c18, 'H': 0x188c, 'L': 0x8c18,
		}

		resetAll()
		cpu.setAcc(0x8c)
		cpu.BC = 0x8c8c
		cpu.DE = 0x8c8c
		cpu.HL = 0x8c8c
		cpu.setFlags(0b11010110)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x18, "Flags": 0b00000101}, cpu.rlR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b00000101}, cpu.rlR(register))
		}

		expectedValueMap = map[byte]uint16{
			'B': 0x9b4d, 'C': 0x4d9b, 'D': 0x9b4d, 'E': 0x4d9b, 'H': 0x9b4d, 'L': 0x4d9b,
		}

		resetAll()
		cpu.setAcc(0x4d)
		cpu.BC = 0x4d4d
		cpu.DE = 0x4d4d
		cpu.HL = 0x4d4d
		cpu.setFlags(0b11010111)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x9b, "Flags": 0b10000000}, cpu.rlR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b10000000}, cpu.rlR(register))
		}
	}
}

func TestRlHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8c)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b00000101}, cpu.rlSs("HL"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x18)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	dmaX.SetMemoryByte(0x1234, 0x4d)
	cpu.setFlags(0b11010111)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b10000000}, cpu.rlSs("HL"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x9b)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRlIx(t *testing.T) {
	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b00000101}, cpu.rlSs("IX"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x18)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IX = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x4d)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	cpu.setFlags(0b11010111)
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b10000000}, cpu.rlSs("IX"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x9b)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRlIy(t *testing.T) {
	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b00000101}, cpu.rlSs("IY"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x18)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IY = 0x121b
	dmaX.SetMemoryByte(0x1234, 0x4d)
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b10000000}, cpu.rlSs("IY"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x9b)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRrR(t *testing.T) {
	expectedRegisterMap := map[byte]string{
		'B': "BC", 'C': "BC", 'D': "DE", 'E': "DE", 'H': "HL", 'L': "HL", 'A': "A",
	}
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		expectedValueMap := map[byte]uint16{
			'B': 0x468d, 'C': 0x8d46, 'D': 0x468d, 'E': 0x8d46, 'H': 0x468d, 'L': 0x8d46,
		}

		resetAll()
		cpu.setAcc(0x8d)
		cpu.BC = 0x8d8d
		cpu.DE = 0x8d8d
		cpu.HL = 0x8d8d
		cpu.setFlags(0b11010110)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x46, "Flags": 0b00000001}, cpu.rrR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b00000001}, cpu.rrR(register))
		}

		expectedValueMap = map[byte]uint16{
			'B': 0xa64c, 'C': 0x4ca6, 'D': 0xa64c, 'E': 0x4ca6, 'H': 0xa64c, 'L': 0x4ca6,
		}

		resetAll()
		cpu.setAcc(0x4c)
		cpu.BC = 0x4c4c
		cpu.DE = 0x4c4c
		cpu.HL = 0x4c4c
		cpu.setFlags(0b11010111)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0xa6, "Flags": 0b10000100}, cpu.rrR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b10000100}, cpu.rrR(register))
		}
	}
}

func TestRrHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8d)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b00000001}, cpu.rrSs("HL"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x46)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x4c)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b10000100}, cpu.rrSs("HL"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0xa6)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRrIx(t *testing.T) {
	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8d)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b00000001}, cpu.rrSs("IX"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x46)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x4c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b10000100}, cpu.rrSs("IX"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0xa6)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestRrIy(t *testing.T) {
	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x8d)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b00000001}, cpu.rrSs("IY"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x46)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x4c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b10000100}, cpu.rrSs("IY"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0xa6)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestSlaR(t *testing.T) {
	expectedRegisterMap := map[byte]string{
		'B': "BC", 'C': "BC", 'D': "DE", 'E': "DE", 'H': "HL", 'L': "HL", 'A': "A",
	}
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		expectedValueMap := map[byte]uint16{
			'B': 0x0080, 'C': 0x8000, 'D': 0x0080, 'E': 0x8000, 'H': 0x0080, 'L': 0x8000,
		}

		resetAll()
		cpu.setAcc(0x80)
		cpu.BC = 0x8080
		cpu.DE = 0x8080
		cpu.HL = 0x8080
		cpu.setFlags(0b10010110)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x00, "Flags": 0b01000101}, cpu.slaR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b01000101}, cpu.slaR(register))
		}

		expectedValueMap = map[byte]uint16{
			'B': 0x984c, 'C': 0x4c98, 'D': 0x984c, 'E': 0x4c98, 'H': 0x984c, 'L': 0x4c98,
		}

		resetAll()
		cpu.setAcc(0x4c)
		cpu.BC = 0x4c4c
		cpu.DE = 0x4c4c
		cpu.HL = 0x4c4c
		cpu.setFlags(0b01010111)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x98, "Flags": 0b10000000}, cpu.slaR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b10000000}, cpu.slaR(register))
		}
	}
}

func TestSlaHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x80)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b01000101}, cpu.slaSs("HL"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x4c)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b10000000}, cpu.slaSs("HL"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x98)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestSlaIx(t *testing.T) {
	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x80)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b01000101}, cpu.slaSs("IX"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x4c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b10000000}, cpu.slaSs("IX"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x98)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestSlaIy(t *testing.T) {
	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x80)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b01000101}, cpu.slaSs("IY"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x4c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b10000000}, cpu.slaSs("IY"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x98)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestSraR(t *testing.T) {
	expectedRegisterMap := map[byte]string{
		'B': "BC", 'C': "BC", 'D': "DE", 'E': "DE", 'H': "HL", 'L': "HL", 'A': "A",
	}
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		expectedValueMap := map[byte]uint16{
			'B': 0x0001, 'C': 0x0100, 'D': 0x0001, 'E': 0x0100, 'H': 0x0001, 'L': 0x0100,
		}

		resetAll()
		cpu.setAcc(0x01)
		cpu.BC = 0x0101
		cpu.DE = 0x0101
		cpu.HL = 0x0101
		cpu.setFlags(0b10010110)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x00, "Flags": 0b01000101}, cpu.sraR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b01000101}, cpu.sraR(register))
		}

		expectedValueMap = map[byte]uint16{
			'B': 0xc78e, 'C': 0x8ec7, 'D': 0xc78e, 'E': 0x8ec7, 'H': 0xc78e, 'L': 0x8ec7,
		}

		resetAll()
		cpu.setAcc(0x8e)
		cpu.BC = 0x8e8e
		cpu.DE = 0x8e8e
		cpu.HL = 0x8e8e
		cpu.setFlags(0b01010111)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0xc7, "Flags": 0b10000000}, cpu.sraR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b10000000}, cpu.sraR(register))
		}
	}
}

func TestSraHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x01)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b01000101}, cpu.sraSs("HL"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x8e)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b10000000}, cpu.sraSs("HL"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0xc7)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestSraIx(t *testing.T) {
	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x01)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b01000101}, cpu.sraSs("IX"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x8e)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b10000000}, cpu.sraSs("IX"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0xc7)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestSraIy(t *testing.T) {
	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x01)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b01000101}, cpu.sraSs("IY"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x8e)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b10000000}, cpu.sraSs("IY"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0xc7)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestSllR(t *testing.T) {
	expectedRegisterMap := map[byte]string{
		'B': "BC", 'C': "BC", 'D': "DE", 'E': "DE", 'H': "HL", 'L': "HL", 'A': "A",
	}
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		expectedValueMap := map[byte]uint16{
			'B': 0x0180, 'C': 0x8001, 'D': 0x0180, 'E': 0x8001, 'H': 0x0180, 'L': 0x8001,
		}

		resetAll()
		cpu.setAcc(0x80)
		cpu.BC = 0x8080
		cpu.DE = 0x8080
		cpu.HL = 0x8080
		cpu.setFlags(0b10010110)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x01, "Flags": 0b00000001}, cpu.sllR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b00000001}, cpu.sllR(register))
		}

		expectedValueMap = map[byte]uint16{
			'B': 0x994c, 'C': 0x4c99, 'D': 0x994c, 'E': 0x4c99, 'H': 0x994c, 'L': 0x4c99,
		}

		resetAll()
		cpu.setAcc(0x4c)
		cpu.BC = 0x4c4c
		cpu.DE = 0x4c4c
		cpu.HL = 0x4c4c
		cpu.setFlags(0b01010111)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x99, "Flags": 0b10000100}, cpu.sllR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b10000100}, cpu.sllR(register))
		}
	}
}

func TestSllHl(t *testing.T) {
	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x80)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b00000001}, cpu.sllSs("HL"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x01)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.HL = 0x1234
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x4c)
	checkCpu(t, 15, map[string]uint16{"PC": 2, "HL": 0x1234, "Flags": 0b10000100}, cpu.sllSs("HL"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x99)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestSllIx(t *testing.T) {
	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x80)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b00000001}, cpu.sllSs("IX"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x01)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IX = 0x121b
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x4c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IX": 0x121b, "Flags": 0b10000100}, cpu.sllSs("IX"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x99)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestSllIy(t *testing.T) {
	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010110)
	dmaX.SetMemoryByte(0x1234, 0x80)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b00000001}, cpu.sllSs("IY"))

	got := dmaX.GetMemory(0x1234)
	want := uint8(0x01)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.IY = 0x121b
	cpu.setFlags(0b11010111)
	dmaX.SetMemoryByte(0x1234, 0x4c)
	dmaX.SetMemoryBulk(0x0000, []uint8{0xdd, 0xcb, 0x06, 0x19})
	checkCpu(t, 23, map[string]uint16{"PC": 4, "IY": 0x121b, "Flags": 0b10000100}, cpu.sllSs("IY"))

	got = dmaX.GetMemory(0x1234)
	want = uint8(0x99)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestSrlR(t *testing.T) {
	expectedRegisterMap := map[byte]string{
		'B': "BC", 'C': "BC", 'D': "DE", 'E': "DE", 'H': "HL", 'L': "HL", 'A': "A",
	}
	for _, register := range []byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
		expectedValueMap := map[byte]uint16{
			'B': 0x0001, 'C': 0x0100, 'D': 0x0001, 'E': 0x0100, 'H': 0x0001, 'L': 0x0100,
		}

		resetAll()
		cpu.setAcc(0x01)
		cpu.BC = 0x0101
		cpu.DE = 0x0101
		cpu.HL = 0x0101
		cpu.setFlags(0b10010110)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x00, "Flags": 0b01000101}, cpu.srlR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b01000101}, cpu.srlR(register))
		}

		expectedValueMap = map[byte]uint16{
			'B': 0x67ce, 'C': 0xce67, 'D': 0x67ce, 'E': 0xce67, 'H': 0x67ce, 'L': 0xce67,
		}

		resetAll()
		cpu.setAcc(0xce)
		cpu.BC = 0xcece
		cpu.DE = 0xcece
		cpu.HL = 0xcece
		cpu.setFlags(0b01010111)

		switch register {
		case 'A':
			checkCpu(t, 8, map[string]uint16{"PC": 2, "A": 0x67, "Flags": 0b00000000}, cpu.srlR(register))
		default:
			checkCpu(t, 8, map[string]uint16{"PC": 2, expectedRegisterMap[register]: expectedValueMap[register], "Flags": 0b00000000}, cpu.srlR(register))
		}
	}
}
