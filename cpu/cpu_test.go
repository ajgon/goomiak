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
