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
			t.Errorf("AF: got %d, want %d", cpu.AF, af)
		}
	}

	if bc, ok := expected["BC"]; ok {
		if cpu.BC != bc {
			t.Errorf("BC: got %d, want %d", cpu.BC, bc)
		}
	}

	if flags, ok := expected["Flags"]; ok {
		if cpu.Flags != uint8(flags) {
			t.Errorf("Flags: got %d, want %d", cpu.Flags, flags)
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
