package cpu

import (
	"testing"
	"z80/dma"
	"z80/memory"
)

var sbcTruthTable [72][10]uint8 = [72][10]uint8{
	// a, r, C before, a-r, C, N, PV, H, Z, S
	[10]uint8{0, 0, 0, 0, 0, 1, 0, 0, 1, 0},
	[10]uint8{0, 1, 0, 255, 1, 1, 0, 1, 0, 1},
	[10]uint8{0, 127, 0, 129, 1, 1, 0, 1, 0, 1},
	[10]uint8{0, 128, 0, 128, 1, 1, 1, 0, 0, 1},
	[10]uint8{0, 129, 0, 127, 1, 1, 0, 1, 0, 0},
	[10]uint8{0, 255, 0, 1, 1, 1, 0, 1, 0, 0},
	[10]uint8{1, 0, 0, 1, 0, 1, 0, 0, 0, 0},
	[10]uint8{1, 1, 0, 0, 0, 1, 0, 0, 1, 0},
	[10]uint8{1, 127, 0, 130, 1, 1, 0, 1, 0, 1},
	[10]uint8{1, 128, 0, 129, 1, 1, 1, 0, 0, 1},
	[10]uint8{1, 129, 0, 128, 1, 1, 1, 0, 0, 1},
	[10]uint8{1, 255, 0, 2, 1, 1, 0, 1, 0, 0},
	[10]uint8{127, 0, 0, 127, 0, 1, 0, 0, 0, 0},
	[10]uint8{127, 1, 0, 126, 0, 1, 0, 0, 0, 0},
	[10]uint8{127, 127, 0, 0, 0, 1, 0, 0, 1, 0},
	[10]uint8{127, 128, 0, 255, 1, 1, 1, 0, 0, 1},
	[10]uint8{127, 129, 0, 254, 1, 1, 1, 0, 0, 1},
	[10]uint8{127, 255, 0, 128, 1, 1, 1, 0, 0, 1},
	[10]uint8{128, 0, 0, 128, 0, 1, 0, 0, 0, 1},
	[10]uint8{128, 1, 0, 127, 0, 1, 1, 1, 0, 0},
	[10]uint8{128, 127, 0, 1, 0, 1, 1, 1, 0, 0},
	[10]uint8{128, 128, 0, 0, 0, 1, 0, 0, 1, 0},
	[10]uint8{128, 129, 0, 255, 1, 1, 0, 1, 0, 1},
	[10]uint8{128, 255, 0, 129, 1, 1, 0, 1, 0, 1},
	[10]uint8{129, 0, 0, 129, 0, 1, 0, 0, 0, 1},
	[10]uint8{129, 1, 0, 128, 0, 1, 0, 0, 0, 1},
	[10]uint8{129, 127, 0, 2, 0, 1, 1, 1, 0, 0},
	[10]uint8{129, 128, 0, 1, 0, 1, 0, 0, 0, 0},
	[10]uint8{129, 129, 0, 0, 0, 1, 0, 0, 1, 0},
	[10]uint8{129, 255, 0, 130, 1, 1, 0, 1, 0, 1},
	[10]uint8{255, 0, 0, 255, 0, 1, 0, 0, 0, 1},
	[10]uint8{255, 1, 0, 254, 0, 1, 0, 0, 0, 1},
	[10]uint8{255, 127, 0, 128, 0, 1, 0, 0, 0, 1},
	[10]uint8{255, 128, 0, 127, 0, 1, 0, 0, 0, 0},
	[10]uint8{255, 129, 0, 126, 0, 1, 0, 0, 0, 0},
	[10]uint8{255, 255, 0, 0, 0, 1, 0, 0, 1, 0},

	[10]uint8{0, 0, 1, 255, 1, 1, 0, 1, 0, 1},
	[10]uint8{0, 1, 1, 254, 1, 1, 0, 1, 0, 1},
	[10]uint8{0, 127, 1, 128, 1, 1, 0, 1, 0, 1},
	[10]uint8{0, 128, 1, 127, 1, 1, 0, 1, 0, 0},
	[10]uint8{0, 129, 1, 126, 1, 1, 0, 1, 0, 0},
	[10]uint8{0, 255, 1, 0, 1, 1, 0, 1, 1, 0},
	[10]uint8{1, 0, 1, 0, 0, 1, 0, 0, 1, 0},
	[10]uint8{1, 1, 1, 255, 1, 1, 0, 1, 0, 1},
	[10]uint8{1, 127, 1, 129, 1, 1, 0, 1, 0, 1},
	[10]uint8{1, 128, 1, 128, 1, 1, 1, 0, 0, 1},
	[10]uint8{1, 129, 1, 127, 1, 1, 0, 1, 0, 0},
	[10]uint8{1, 255, 1, 1, 1, 1, 0, 1, 0, 0},
	[10]uint8{127, 0, 1, 126, 0, 1, 0, 0, 0, 0},
	[10]uint8{127, 1, 1, 125, 0, 1, 0, 0, 0, 0},
	[10]uint8{127, 127, 1, 255, 1, 1, 0, 1, 0, 1},
	[10]uint8{127, 128, 1, 254, 1, 1, 1, 0, 0, 1},
	[10]uint8{127, 129, 1, 253, 1, 1, 1, 0, 0, 1},
	[10]uint8{127, 255, 1, 127, 1, 1, 0, 1, 0, 0},
	[10]uint8{128, 0, 1, 127, 0, 1, 1, 1, 0, 0},
	[10]uint8{128, 1, 1, 126, 0, 1, 1, 1, 0, 0},
	[10]uint8{128, 127, 1, 0, 0, 1, 1, 1, 1, 0},
	[10]uint8{128, 128, 1, 255, 1, 1, 0, 1, 0, 1},
	[10]uint8{128, 129, 1, 254, 1, 1, 0, 1, 0, 1},
	[10]uint8{128, 255, 1, 128, 1, 1, 0, 1, 0, 1},
	[10]uint8{129, 0, 1, 128, 0, 1, 0, 0, 0, 1},
	[10]uint8{129, 1, 1, 127, 0, 1, 1, 1, 0, 0},
	[10]uint8{129, 127, 1, 1, 0, 1, 1, 1, 0, 0},
	[10]uint8{129, 128, 1, 0, 0, 1, 0, 0, 1, 0},
	[10]uint8{129, 129, 1, 255, 1, 1, 0, 1, 0, 1},
	[10]uint8{129, 255, 1, 129, 1, 1, 0, 1, 0, 1},
	[10]uint8{255, 0, 1, 254, 0, 1, 0, 0, 0, 1},
	[10]uint8{255, 1, 1, 253, 0, 1, 0, 0, 0, 1},
	[10]uint8{255, 127, 1, 127, 0, 1, 1, 1, 0, 0},
	[10]uint8{255, 128, 1, 126, 0, 1, 0, 0, 0, 0},
	[10]uint8{255, 129, 1, 125, 0, 1, 0, 0, 0, 0},
	[10]uint8{255, 255, 1, 255, 1, 1, 0, 1, 0, 1},
}

func TestSbcRegister(t *testing.T) {
	var mem = memory.MemoryNew()
	var dmaX = dma.DMANew(mem)
	var cpu = CPUNew(dmaX)

	for _, row := range sbcTruthTable {
		for _, register := range [7]byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
			if register == 'A' {
				if row[0] != row[1] {
					continue
				}
			}

			cpu.PC = 0
			cpu.setC(row[2] == 1)
			cpu.setAcc(row[0])
			cpu.BC = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.DE = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.HL = (uint16(row[1]) << 8) | uint16(row[1])
			tstates := cpu.sbcR(register)()

			if cpu.getAcc() != row[3] || cpu.getC() != (row[4] == 1) || cpu.getN() != (row[5] == 1) || cpu.getPV() != (row[6] == 1) || cpu.getH() != (row[7] == 1) || cpu.getZ() != (row[8] == 1) || cpu.getS() != (row[9] == 1) {
				t.Errorf(
					"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d - %d - %d)",
					cpu.getAcc(), cpu.getC(), cpu.getN(), cpu.getPV(), cpu.getH(), cpu.getZ(), cpu.getS(),
					row[3], row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[9] == 1, row[0], row[1], row[2],
				)
			}

			if cpu.PC != 1 || tstates != 4 {
				t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, tstates, 1, 4)
			}
		}
	}
}

func TestSbc_Hl_(t *testing.T) {
	var mem = memory.MemoryNew()
	var dmaX = dma.DMANew(mem)
	var cpu = CPUNew(dmaX)
	cpu.HL = 0x1234

	for _, row := range sbcTruthTable {
		cpu.PC = 0
		cpu.setC(row[2] == 1)
		cpu.setAcc(row[0])
		dmaX.SetMemoryByte(cpu.HL, row[1])
		tstates := cpu.sbc_Hl_()

		if cpu.getAcc() != row[3] || cpu.getC() != (row[4] == 1) || cpu.getN() != (row[5] == 1) || cpu.getPV() != (row[6] == 1) || cpu.getH() != (row[7] == 1) || cpu.getZ() != (row[8] == 1) || cpu.getS() != (row[9] == 1) {
			t.Errorf(
				"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d - %d - %d)",
				cpu.getAcc(), cpu.getC(), cpu.getN(), cpu.getPV(), cpu.getH(), cpu.getZ(), cpu.getS(),
				row[3], row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[9] == 1, row[0], row[1], row[2],
			)
		}

		if cpu.PC != 1 || tstates != 7 {
			t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, tstates, 1, 7)
		}
	}
}

func TestSbcX(t *testing.T) {
	var mem = memory.MemoryNew()
	var dmaX = dma.DMANew(mem)
	var cpu = CPUNew(dmaX)

	for _, row := range sbcTruthTable {
		cpu.PC = 0
		cpu.setC(row[2] == 1)
		cpu.setAcc(row[0])
		dmaX.SetMemoryByte(0x01, row[1])
		tstates := cpu.sbcAX()

		if cpu.getAcc() != row[3] || cpu.getC() != (row[4] == 1) || cpu.getN() != (row[5] == 1) || cpu.getPV() != (row[6] == 1) || cpu.getH() != (row[7] == 1) || cpu.getZ() != (row[8] == 1) || cpu.getS() != (row[9] == 1) {
			t.Errorf(
				"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d - %d - %d)",
				cpu.getAcc(), cpu.getC(), cpu.getN(), cpu.getPV(), cpu.getH(), cpu.getZ(), cpu.getS(),
				row[3], row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[9] == 1, row[0], row[1], row[2],
			)
		}

		if cpu.PC != 2 || tstates != 7 {
			t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, tstates, 2, 7)
		}
	}
}
