package cpu

import (
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

	cycles := instructionCall()

	if pc, ok := expected["PC"]; ok {
		if cpu.PC != pc {
			t.Errorf("PC: got %d, want %d", cpu.PC, pc)
		}
	} else {
		panic("Every mnemonic test should validate PC!")
	}

	if af, ok := expected["AF"]; ok {
		if cpu.AF != af {
			t.Errorf("AF: got %x, want %x", cpu.AF, af)
		}
	}

	if af_, ok := expected["AF_"]; ok {
		if cpu.AF_ != af_ {
			t.Errorf("AF': got %x, want %x", cpu.AF_, af_)
		}
	}

	if bc, ok := expected["BC"]; ok {
		if cpu.BC != bc {
			t.Errorf("BC: got %x, want %x", cpu.BC, bc)
		}
	}

	if de, ok := expected["DE"]; ok {
		if cpu.DE != de {
			t.Errorf("DE: got %x, want %x", cpu.DE, de)
		}
	}

	if hl, ok := expected["HL"]; ok {
		if cpu.HL != hl {
			t.Errorf("HL: got %x, want %x", cpu.HL, hl)
		}
	}

	if flags, ok := expected["Flags"]; ok {
		if cpu.Flags != uint8(flags) {
			t.Errorf("Flags: got %b, want %b", cpu.Flags, flags)
		}
	}

	if cycles != expectedCycles {
		t.Errorf("cycles: got %d, want %d", cycles, expectedCycles)
	}
}

func TestNop(t *testing.T) {
	cpu.Reset()
	checkCpu(t, 4, map[string]uint16{"PC": 1}, cpu.nop)
}

func TestLdBcXx(t *testing.T) {
	cpu.Reset()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x01, 0x64, 0x32})

	checkCpu(t, 10, map[string]uint16{"PC": 3, "BC": 0x3264}, cpu.ldBcXx)
}

func TestLdBcA(t *testing.T) {
	cpu.Reset()
	cpu.AF = 0x7A05
	cpu.BC = 0x1015

	checkCpu(t, 7, map[string]uint16{"PC": 1}, cpu.ldBcA)

	got := dmaX.GetMemory(0x1015)
	want := uint8(0x7A)
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestIncBc(t *testing.T) {
	cpu.Reset()
	cpu.BC = 0x1020

	checkCpu(t, 6, map[string]uint16{"PC": 1, "BC": 0x1021}, cpu.incBc)
}

func TestIncB(t *testing.T) {
	cpu.Reset()
	cpu.Flags = 0b11010111
	cpu.BC = 0x1002

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x1102, "Flags": 0b00000001}, cpu.incB)

	cpu.Reset()
	cpu.Flags = 0b10000110
	cpu.BC = 0xff02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0002, "Flags": 0b01010000}, cpu.incB)

	cpu.Reset()
	cpu.Flags = 0b01000010
	cpu.BC = 0x7f02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x8002, "Flags": 0b10010100}, cpu.incB)
}

func TestDecB(t *testing.T) {
	cpu.Reset()
	cpu.Flags = 0b11010101
	cpu.BC = 0x0102

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0002, "Flags": 0b01000011}, cpu.decB)

	cpu.Reset()
	cpu.Flags = 0b01000100
	cpu.BC = 0x0002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0xff02, "Flags": 0b10010010}, cpu.decB)

	cpu.Reset()
	cpu.Flags = 0b11000000
	cpu.BC = 0x8002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x7f02, "Flags": 0b00010110}, cpu.decB)
}

func TestLdBX(t *testing.T) {
	cpu.Reset()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "BC": 0x6400}, cpu.ldBX)
}

func TestRlca(t *testing.T) {
	cpu.Reset()
	cpu.AF = 0x8c05
	cpu.Flags = 0b11010110
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x1905, "Flags": 0b11000101}, cpu.rlca)

	cpu.Reset()
	cpu.AF = 0x4d05
	cpu.Flags = 0b11010111
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x9a05, "Flags": 0b11000100}, cpu.rlca)
}

func TestExAfAf_(t *testing.T) {
	cpu.Reset()
	cpu.AF = 0x1234
	cpu.AF_ = 0x5678
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x5678, "AF_": 0x1234}, cpu.exAfAf_)
}

func TestAddHlBc(t *testing.T) {
	cpu.Reset()
	cpu.BC = 0xa76c //  1010 0111 0110 1100
	cpu.HL = 0x5933 //  0101 1001 0011 0011
	cpu.Flags = 0b00000010

	checkCpu(t, 11, map[string]uint16{"PC": 1, "BC": 0xa76c, "HL": 0x009f, "Flags": 0b00010001}, cpu.addHlBc)

	cpu.Reset()
	cpu.BC = 0x7fff
	cpu.HL = 0x7fff
	cpu.Flags = 0b00000010

	checkCpu(t, 11, map[string]uint16{"PC": 1, "BC": 0x7fff, "HL": 0xfffe, "Flags": 0b00010000}, cpu.addHlBc)
}

func TestLdABc(t *testing.T) {
	cpu.Reset()
	dmaX.SetMemoryByte(0x1257, 0x64)
	cpu.AF = 0xffff
	cpu.BC = 0x1257

	checkCpu(t, 7, map[string]uint16{"PC": 1, "AF": 0x64ff, "BC": 0x1257}, cpu.ldABc)
}

func TestDecBc(t *testing.T) {
	cpu.Reset()
	cpu.BC = 0x1000

	checkCpu(t, 6, map[string]uint16{"PC": 1, "BC": 0x0fff}, cpu.decBc)
}

func TestIncC(t *testing.T) {
	cpu.Reset()
	cpu.Flags = 0b11010111
	cpu.BC = 0x0210

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0211, "Flags": 0b00000001}, cpu.incC)

	cpu.Reset()
	cpu.Flags = 0b10000110
	cpu.BC = 0x02ff
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0200, "Flags": 0b01010000}, cpu.incC)

	cpu.Reset()
	cpu.Flags = 0b01000010
	cpu.BC = 0x027f
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0280, "Flags": 0b10010100}, cpu.incC)
}

func TestDecC(t *testing.T) {
	cpu.Reset()
	cpu.Flags = 0b11010101
	cpu.BC = 0x0201

	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x0200, "Flags": 0b01000011}, cpu.decC)

	cpu.Reset()
	cpu.Flags = 0b01000100
	cpu.BC = 0x0200
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x02ff, "Flags": 0b10010010}, cpu.decC)

	cpu.Reset()
	cpu.Flags = 0b11000000
	cpu.BC = 0x0280
	checkCpu(t, 4, map[string]uint16{"PC": 1, "BC": 0x027f, "Flags": 0b00010110}, cpu.decC)
}

func TestLdCX(t *testing.T) {
	cpu.Reset()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "BC": 0x0064}, cpu.ldCX)
}

func TestRrca(t *testing.T) {
	cpu.Reset()
	cpu.AF = 0x8d05
	cpu.Flags = 0b11010110
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0xc605, "Flags": 0b11000101}, cpu.rrca)

	cpu.Reset()
	cpu.AF = 0x4c05
	cpu.Flags = 0b11010111
	checkCpu(t, 4, map[string]uint16{"PC": 1, "AF": 0x2605, "Flags": 0b11000100}, cpu.rrca)
}

func TestDjnzX(t *testing.T) {
	cpu.Reset()
	cpu.PC = 3
	cpu.BC = 0x1234
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 13, map[string]uint16{"PC": 0x37, "BC": 0x1134}, cpu.djnzX)

	cpu.Reset()
	cpu.PC = 3
	cpu.BC = 0x0134
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 8, map[string]uint16{"PC": 0x05, "BC": 0x0034}, cpu.djnzX)

	cpu.Reset()
	cpu.PC = 3
	cpu.BC = 0x0034
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0x32})

	checkCpu(t, 13, map[string]uint16{"PC": 0x37, "BC": 0xff34}, cpu.djnzX)

	cpu.Reset()
	cpu.PC = 3
	cpu.BC = 0x0534
	dmaX.SetMemoryBulk(0x0003, []uint8{0x10, 0xfb})

	checkCpu(t, 13, map[string]uint16{"PC": 0x00, "BC": 0x0434}, cpu.djnzX)
}

func TestLdDeXx(t *testing.T) {
	cpu.Reset()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x01, 0x64, 0x32})

	checkCpu(t, 10, map[string]uint16{"PC": 3, "DE": 0x3264}, cpu.ldDeXx)
}

func TestLdDeA(t *testing.T) {
	cpu.Reset()
	cpu.AF = 0x7A05
	cpu.DE = 0x1015

	checkCpu(t, 7, map[string]uint16{"PC": 1}, cpu.ldDeA)

	got := dmaX.GetMemory(0x1015)
	want := uint8(0x7A)
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestIncDe(t *testing.T) {
	cpu.Reset()
	cpu.DE = 0x1020

	checkCpu(t, 6, map[string]uint16{"PC": 1, "DE": 0x1021}, cpu.incDe)
}

func TestIncD(t *testing.T) {
	cpu.Reset()
	cpu.Flags = 0b11010111
	cpu.DE = 0x1002

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x1102, "Flags": 0b00000001}, cpu.incD)

	cpu.Reset()
	cpu.Flags = 0b10000110
	cpu.DE = 0xff02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0002, "Flags": 0b01010000}, cpu.incD)

	cpu.Reset()
	cpu.Flags = 0b01000010
	cpu.DE = 0x7f02
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x8002, "Flags": 0b10010100}, cpu.incD)
}

func TestDecD(t *testing.T) {
	cpu.Reset()
	cpu.Flags = 0b11010101
	cpu.DE = 0x0102

	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x0002, "Flags": 0b01000011}, cpu.decD)

	cpu.Reset()
	cpu.Flags = 0b01000100
	cpu.DE = 0x0002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0xff02, "Flags": 0b10010010}, cpu.decD)

	cpu.Reset()
	cpu.Flags = 0b11000000
	cpu.DE = 0x8002
	checkCpu(t, 4, map[string]uint16{"PC": 1, "DE": 0x7f02, "Flags": 0b00010110}, cpu.decD)
}

func TestLdDX(t *testing.T) {
	cpu.Reset()
	dmaX.SetMemoryBulk(0x0000, []uint8{0x06, 0x64})

	checkCpu(t, 7, map[string]uint16{"PC": 2, "DE": 0x6400}, cpu.ldDX)
}
