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
	var expectedSP, expectedAF, expectedAF_, expectedBC, expectedDE, expectedHL uint16
	var expectedFlags uint8

	if sp, ok := expected["SP"]; ok {
		expectedSP = sp
	} else {
		cpu.SP = uint16(rand.Uint32())
		expectedSP = cpu.SP
	}

	if af, ok := expected["AF"]; ok {
		expectedAF = af
	} else {
		cpu.AF = uint16(rand.Uint32())
		expectedAF = cpu.AF
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

	if de, ok := expected["DE"]; ok {
		expectedDE = de
	} else {
		cpu.DE = uint16(rand.Uint32())
		expectedDE = cpu.DE
	}

	if hl, ok := expected["HL"]; ok {
		expectedHL = hl
	} else {
		cpu.HL = uint16(rand.Uint32())
		expectedHL = cpu.HL
	}

	if flags, ok := expected["Flags"]; ok {
		expectedFlags = uint8(flags)
	} else {
		expectedFlags = uint8(rand.Uint32() & 0b11010111)
		cpu.Flags.fromRegister(expectedFlags)
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

	if cpu.AF != expectedAF {
		t.Errorf("AF: got %x, want %x", cpu.AF, expectedAF)
	}

	if cpu.AF_ != expectedAF_ {
		t.Errorf("AF': got %x, want %x", cpu.AF_, expectedAF_)
	}

	if cpu.BC != expectedBC {
		t.Errorf("BC: got %x, want %x", cpu.BC, expectedBC)
	}

	if cpu.DE != expectedDE {
		t.Errorf("DE: got %x, want %x", cpu.DE, expectedDE)
	}

	if cpu.HL != expectedHL {
		t.Errorf("HL: got %x, want %x", cpu.HL, expectedHL)
	}

	if cpu.Flags.toRegister() != expectedFlags {
		t.Errorf("Flags: got %08b, want %08b", cpu.Flags.toRegister(), expectedFlags)
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
	cpu.AF = 0x7A05
	cpu.BC = 0x1015

	checkCpu(t, 7, map[string]uint16{"PC": 1, "AF": 0x7a05, "BC": 0x1015}, cpu.ld_Bc_A)

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
	cpu.Flags.fromRegister(0b11010111)
	cpu.BC = 0x1002

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1102, "Flags": 0b00000001}, cpu.incB)

	resetAll()
	cpu.Flags.fromRegister(0b10000110)
	cpu.BC = 0xff02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0002, "Flags": 0b01010000}, cpu.incB)

	resetAll()
	cpu.Flags.fromRegister(0b01000010)
	cpu.BC = 0x7f02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x8002, "Flags": 0b10010100}, cpu.incB)
}

func TestDecB(t *testing.T) {
	resetAll()
	cpu.Flags.fromRegister(0b11010101)
	cpu.BC = 0x0102

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0002, "Flags": 0b01000011}, cpu.decB)

	resetAll()
	cpu.Flags.fromRegister(0b01000100)
	cpu.BC = 0x0002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0xff02, "Flags": 0b10010010}, cpu.decB)

	resetAll()
	cpu.Flags.fromRegister(0b11000000)
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
	cpu.AF = 0x8c05
	cpu.Flags.fromRegister(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x1905, "Flags": 0b11000101}, cpu.rlca)

	resetAll()
	cpu.AF = 0x4d05
	cpu.Flags.fromRegister(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x9a05, "Flags": 0b11000100}, cpu.rlca)
}

func TestExAfAf_(t *testing.T) {
	resetAll()
	cpu.AF = 0x1234
	cpu.AF_ = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x5678, "AF_": 0x1234}, cpu.exAfAf_)
}

func TestAddHlBc(t *testing.T) {
	resetAll()
	cpu.BC = 0xa76c //  1010 0111 0110 1100
	cpu.HL = 0x5933 //  0101 1001 0011 0011
	cpu.Flags.fromRegister(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "BC": 0xa76c, "HL": 0x009f, "Flags": 0b00010001}, cpu.addHlBc)

	resetAll()
	cpu.BC = 0x7fff
	cpu.HL = 0x7fff
	cpu.Flags.fromRegister(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "BC": 0x7fff, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addHlBc)
}

func TestLdA_Bc_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryByte(0x1257, 0x64)
	cpu.AF = 0xffff
	cpu.BC = 0x1257

	checkCpu(t, 7, map[string]uint16{"PC": 1, "AF": 0x64ff, "BC": 0x1257}, cpu.ldA_Bc_)
}

func TestDecBc(t *testing.T) {
	resetAll()
	cpu.BC = 0x1000

	checkCpu(t, 6, map[string]uint16{"PC": 1, "BC": 0x0fff}, cpu.decBc)
}

func TestIncC(t *testing.T) {
	resetAll()
	cpu.Flags.fromRegister(0b11010111)
	cpu.BC = 0x0210

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0211, "Flags": 0b00000001}, cpu.incC)

	resetAll()
	cpu.Flags.fromRegister(0b10000110)
	cpu.BC = 0x02ff
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0200, "Flags": 0b01010000}, cpu.incC)

	resetAll()
	cpu.Flags.fromRegister(0b01000010)
	cpu.BC = 0x027f
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0280, "Flags": 0b10010100}, cpu.incC)
}

func TestDecC(t *testing.T) {
	resetAll()
	cpu.Flags.fromRegister(0b11010101)
	cpu.BC = 0x0201

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0200, "Flags": 0b01000011}, cpu.decC)

	resetAll()
	cpu.Flags.fromRegister(0b01000100)
	cpu.BC = 0x0200
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x02ff, "Flags": 0b10010010}, cpu.decC)

	resetAll()
	cpu.Flags.fromRegister(0b11000000)
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
	cpu.AF = 0x8d05
	cpu.Flags.fromRegister(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0xc605, "Flags": 0b11000101}, cpu.rrca)

	resetAll()
	cpu.AF = 0x4c05
	cpu.Flags.fromRegister(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x2605, "Flags": 0b11000100}, cpu.rrca)
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
	cpu.AF = 0x7A05
	cpu.DE = 0x1015

	checkCpu(t, 7, map[string]uint16{"PC": 1, "AF": 0x7a05, "DE": 0x1015}, cpu.ld_De_A)

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
	cpu.Flags.fromRegister(0b11010111)
	cpu.DE = 0x1002

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1102, "Flags": 0b00000001}, cpu.incD)

	resetAll()
	cpu.Flags.fromRegister(0b10000110)
	cpu.DE = 0xff02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0002, "Flags": 0b01010000}, cpu.incD)

	resetAll()
	cpu.Flags.fromRegister(0b01000010)
	cpu.DE = 0x7f02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x8002, "Flags": 0b10010100}, cpu.incD)
}

func TestDecD(t *testing.T) {
	resetAll()
	cpu.Flags.fromRegister(0b11010101)
	cpu.DE = 0x0102

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0002, "Flags": 0b01000011}, cpu.decD)

	resetAll()
	cpu.Flags.fromRegister(0b01000100)
	cpu.DE = 0x0002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0xff02, "Flags": 0b10010010}, cpu.decD)

	resetAll()
	cpu.Flags.fromRegister(0b11000000)
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
	cpu.AF = 0x8c05
	cpu.Flags.fromRegister(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x1805, "Flags": 0b11000101}, cpu.rla)

	resetAll()
	cpu.AF = 0x4d05
	cpu.Flags.fromRegister(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x9b05, "Flags": 0b11000100}, cpu.rla)
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
	cpu.Flags.fromRegister(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "DE": 0xa76c, "HL": 0x009f, "Flags": 0b00010001}, cpu.addHlDe)

	resetAll()
	cpu.DE = 0x7fff
	cpu.HL = 0x7fff
	cpu.Flags.fromRegister(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "DE": 0x7fff, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addHlDe)
}

func TestLdA_De_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryByte(0x1257, 0x64)
	cpu.AF = 0xffff
	cpu.DE = 0x1257

	checkCpu(t, 7, map[string]uint16{"PC": 1, "AF": 0x64ff, "DE": 0x1257}, cpu.ldA_De_)
}

func TestDecDe(t *testing.T) {
	resetAll()
	cpu.DE = 0x1000

	checkCpu(t, 6, map[string]uint16{"PC": 1, "DE": 0x0fff}, cpu.decDe)
}

func TestIncE(t *testing.T) {
	resetAll()
	cpu.Flags.fromRegister(0b11010111)
	cpu.DE = 0x0210

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0211, "Flags": 0b00000001}, cpu.incE)

	resetAll()
	cpu.Flags.fromRegister(0b10000110)
	cpu.DE = 0x02ff
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0200, "Flags": 0b01010000}, cpu.incE)

	resetAll()
	cpu.Flags.fromRegister(0b01000010)
	cpu.DE = 0x027f
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0280, "Flags": 0b10010100}, cpu.incE)
}

func TestDecE(t *testing.T) {
	resetAll()
	cpu.Flags.fromRegister(0b11010101)
	cpu.DE = 0x0201

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0200, "Flags": 0b01000011}, cpu.decE)

	resetAll()
	cpu.Flags.fromRegister(0b01000100)
	cpu.DE = 0x0200
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x02ff, "Flags": 0b10010010}, cpu.decE)

	resetAll()
	cpu.Flags.fromRegister(0b11000000)
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
	cpu.AF = 0x8d05
	cpu.Flags.fromRegister(0b11010110)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x4605, "Flags": 0b11000101}, cpu.rra)

	resetAll()
	cpu.AF = 0x4c05
	cpu.Flags.fromRegister(0b11010111)
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0xa605, "Flags": 0b11000100}, cpu.rra)
}

func TestJrNzX(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.Flags.fromRegister(0b10010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b10010111}, cpu.jrNzX)

	resetAll()
	cpu.PC = 3
	cpu.Flags.fromRegister(0b11010111)
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
	cpu.Flags.fromRegister(0b11010111)
	cpu.HL = 0x1002

	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x1102, "Flags": 0b00000001}, cpu.incH)

	resetAll()
	cpu.Flags.fromRegister(0b10000110)
	cpu.HL = 0xff02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0002, "Flags": 0b01010000}, cpu.incH)

	resetAll()
	cpu.Flags.fromRegister(0b01000010)
	cpu.HL = 0x7f02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x8002, "Flags": 0b10010100}, cpu.incH)
}

func TestDecH(t *testing.T) {
	resetAll()
	cpu.Flags.fromRegister(0b11010101)
	cpu.HL = 0x0102

	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0002, "Flags": 0b01000011}, cpu.decH)

	resetAll()
	cpu.Flags.fromRegister(0b01000100)
	cpu.HL = 0x0002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0xff02, "Flags": 0b10010010}, cpu.decH)

	resetAll()
	cpu.Flags.fromRegister(0b11000000)
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
	cpu.Flags.fromRegister(0b11010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b11010111}, cpu.jrZX)

	resetAll()
	cpu.PC = 3
	cpu.Flags.fromRegister(0b10010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 7, map[string]uint16{"PC": 0x05, "Flags": 0b10010111}, cpu.jrZX)
}

func TestAddHlHl(t *testing.T) {
	resetAll()
	cpu.HL = 0xae6c
	cpu.Flags.fromRegister(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x5cd8, "Flags": 0b00010001}, cpu.addHlHl)

	resetAll()
	cpu.HL = 0x7fff
	cpu.Flags.fromRegister(0b00000010)

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
	cpu.Flags.fromRegister(0b11010111)
	cpu.HL = 0x0210

	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0211, "Flags": 0b00000001}, cpu.incL)

	resetAll()
	cpu.Flags.fromRegister(0b10000110)
	cpu.HL = 0x02ff
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0200, "Flags": 0b01010000}, cpu.incL)

	resetAll()
	cpu.Flags.fromRegister(0b01000010)
	cpu.HL = 0x027f
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0280, "Flags": 0b10010100}, cpu.incL)
}

func TestDecL(t *testing.T) {
	resetAll()
	cpu.Flags.fromRegister(0b11010101)
	cpu.HL = 0x0201

	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x0200, "Flags": 0b01000011}, cpu.decL)

	resetAll()
	cpu.Flags.fromRegister(0b01000100)
	cpu.HL = 0x0200
	checkCpu(t, 4, map[string]uint16{"PC": 1, "HL": 0x02ff, "Flags": 0b10010010}, cpu.decL)

	resetAll()
	cpu.Flags.fromRegister(0b11000000)
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
	cpu.Flags.fromRegister(0b00000000)
	cpu.AF = 0xe725

	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x1825, "Flags": 0b00010010}, cpu.cpl)
}

func TestJrNcX(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.Flags.fromRegister(0b11010110)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b11010110}, cpu.jrNcX)

	resetAll()
	cpu.PC = 3
	cpu.Flags.fromRegister(0b11010111)
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
	cpu.AF = 0xd73a
	dmaX.SetMemoryBulk(0x0000, []uint8{0x32, 0x41, 0x31})

	checkCpu(t, 13, map[string]uint16{"PC": 3, "AF": 0xd73a}, cpu.ld_Xx_A)

	got := dmaX.GetMemory(0x3141)
	want := uint8(0xd7)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}
}

func TestIncSp(t *testing.T) {
	resetAll()
	cpu.SP = 0x1020

	checkCpu(t, 6, map[string]uint16{"PC": 1, "SP": 0x1021}, cpu.incSP)
}

func TestInc_Hl_(t *testing.T) {
	resetAll()
	cpu.Flags.fromRegister(0b11010111)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x25)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b00000001}, cpu.inc_Hl_)

	got := dmaX.GetMemory(0x3572)
	want := uint8(0x26)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.Flags.fromRegister(0b10000110)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0xff)
	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b01010000}, cpu.inc_Hl_)

	got = dmaX.GetMemory(0x3572)
	want = uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.Flags.fromRegister(0b01000010)
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
	cpu.Flags.fromRegister(0b11010101)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x01)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b01000011}, cpu.dec_Hl_)

	got := dmaX.GetMemory(0x3572)
	want := uint8(0x00)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.Flags.fromRegister(0b01000100)
	cpu.HL = 0x3572
	dmaX.SetMemoryByte(0x3572, 0x00)
	checkCpu(t, 11, map[string]uint16{"PC": 1, "HL": 0x3572, "Flags": 0b10010010}, cpu.dec_Hl_)

	got = dmaX.GetMemory(0x3572)
	want = uint8(0xff)

	if got != want {
		t.Errorf("got 0x%x, want 0x%x", got, want)
	}

	resetAll()
	cpu.Flags.fromRegister(0b11000000)
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
	cpu.Flags.fromRegister(0b11010110)

	checkCpu(t, 4, map[string]uint16{"PC": 1, "Flags": 0b11000101}, cpu.scf)

}

func TestJrCX(t *testing.T) {
	resetAll()
	cpu.PC = 3
	cpu.Flags.fromRegister(0b11010111)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x38, 0x32})

	checkCpu(t, 12, map[string]uint16{"PC": 0x37, "Flags": 0b11010111}, cpu.jrCX)

	resetAll()
	cpu.PC = 3
	cpu.Flags.fromRegister(0b11010110)
	dmaX.SetMemoryBulk(0x0003, []uint8{0x38, 0x32})

	checkCpu(t, 7, map[string]uint16{"PC": 0x05, "Flags": 0b11010110}, cpu.jrCX)
}

func TestAddHlSp(t *testing.T) {
	resetAll()
	cpu.SP = 0xa76c //  1010 0111 0110 1100
	cpu.HL = 0x5933 //  0101 1001 0011 0011
	cpu.Flags.fromRegister(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0xa76c, "HL": 0x009f, "Flags": 0b00010001}, cpu.addHlSp)

	resetAll()
	cpu.SP = 0x7fff
	cpu.HL = 0x7fff
	cpu.Flags.fromRegister(0b00000010)

	checkCpu(t, 11, map[string]uint16{"PC": 1, "SP": 0x7fff, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addHlSp)
}

func TestLdA_Xx_(t *testing.T) {
	resetAll()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x3a, 0x57, 0x12})
	dmaX.SetMemoryByte(0x1257, 0x64)
	cpu.AF = 0xffff

	checkCpu(t, 13, map[string]uint16{"PC": 3, "AF": 0x64ff}, cpu.ldA_Xx_)
}
