package cpu

import (
	"testing"
	"z80/dma"
	"z80/memory"
)

var adcTruthTable [72][10]uint8 = [72][10]uint8{
	// a, r, C before, a+r, C after, N, PV, H, Z, S
	[10]uint8{0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
	[10]uint8{0, 1, 0, 1, 0, 0, 0, 0, 0, 0},
	[10]uint8{0, 127, 0, 127, 0, 0, 0, 0, 0, 0},
	[10]uint8{0, 128, 0, 128, 0, 0, 0, 0, 0, 1},
	[10]uint8{0, 129, 0, 129, 0, 0, 0, 0, 0, 1},
	[10]uint8{0, 255, 0, 255, 0, 0, 0, 0, 0, 1},
	[10]uint8{1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
	[10]uint8{1, 1, 0, 2, 0, 0, 0, 0, 0, 0},
	[10]uint8{1, 127, 0, 128, 0, 0, 1, 1, 0, 1},
	[10]uint8{1, 128, 0, 129, 0, 0, 0, 0, 0, 1},
	[10]uint8{1, 129, 0, 130, 0, 0, 0, 0, 0, 1},
	[10]uint8{1, 255, 0, 0, 1, 0, 0, 1, 1, 0},
	[10]uint8{127, 0, 0, 127, 0, 0, 0, 0, 0, 0},
	[10]uint8{127, 1, 0, 128, 0, 0, 1, 1, 0, 1},
	[10]uint8{127, 127, 0, 254, 0, 0, 1, 1, 0, 1},
	[10]uint8{127, 128, 0, 255, 0, 0, 0, 0, 0, 1},
	[10]uint8{127, 129, 0, 0, 1, 0, 0, 1, 1, 0},
	[10]uint8{127, 255, 0, 126, 1, 0, 0, 1, 0, 0},
	[10]uint8{128, 0, 0, 128, 0, 0, 0, 0, 0, 1},
	[10]uint8{128, 1, 0, 129, 0, 0, 0, 0, 0, 1},
	[10]uint8{128, 127, 0, 255, 0, 0, 0, 0, 0, 1},
	[10]uint8{128, 128, 0, 0, 1, 0, 1, 0, 1, 0},
	[10]uint8{128, 129, 0, 1, 1, 0, 1, 0, 0, 0},
	[10]uint8{128, 255, 0, 127, 1, 0, 1, 0, 0, 0},
	[10]uint8{129, 0, 0, 129, 0, 0, 0, 0, 0, 1},
	[10]uint8{129, 1, 0, 130, 0, 0, 0, 0, 0, 1},
	[10]uint8{129, 127, 0, 0, 1, 0, 0, 1, 1, 0},
	[10]uint8{129, 128, 0, 1, 1, 0, 1, 0, 0, 0},
	[10]uint8{129, 129, 0, 2, 1, 0, 1, 0, 0, 0},
	[10]uint8{129, 255, 0, 128, 1, 0, 0, 1, 0, 1},
	[10]uint8{255, 0, 0, 255, 0, 0, 0, 0, 0, 1},
	[10]uint8{255, 1, 0, 0, 1, 0, 0, 1, 1, 0},
	[10]uint8{255, 127, 0, 126, 1, 0, 0, 1, 0, 0},
	[10]uint8{255, 128, 0, 127, 1, 0, 1, 0, 0, 0},
	[10]uint8{255, 129, 0, 128, 1, 0, 0, 1, 0, 1},
	[10]uint8{255, 255, 0, 254, 1, 0, 0, 1, 0, 1},

	[10]uint8{0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
	[10]uint8{0, 1, 1, 2, 0, 0, 0, 0, 0, 0},
	[10]uint8{0, 127, 1, 128, 0, 0, 1, 1, 0, 1},
	[10]uint8{0, 128, 1, 129, 0, 0, 0, 0, 0, 1},
	[10]uint8{0, 129, 1, 130, 0, 0, 0, 0, 0, 1},
	[10]uint8{0, 255, 1, 0, 1, 0, 0, 1, 1, 0},
	[10]uint8{1, 0, 1, 2, 0, 0, 0, 0, 0, 0},
	[10]uint8{1, 1, 1, 3, 0, 0, 0, 0, 0, 0},
	[10]uint8{1, 127, 1, 129, 0, 0, 1, 1, 0, 1},
	[10]uint8{1, 128, 1, 130, 0, 0, 0, 0, 0, 1},
	[10]uint8{1, 129, 1, 131, 0, 0, 0, 0, 0, 1},
	[10]uint8{1, 255, 1, 1, 1, 0, 0, 1, 0, 0},
	[10]uint8{127, 0, 1, 128, 0, 0, 1, 1, 0, 1},
	[10]uint8{127, 1, 1, 129, 0, 0, 1, 1, 0, 1},
	[10]uint8{127, 127, 1, 255, 0, 0, 1, 1, 0, 1},
	[10]uint8{127, 128, 1, 0, 1, 0, 0, 1, 1, 0},
	[10]uint8{127, 129, 1, 1, 1, 0, 0, 1, 0, 0},
	[10]uint8{127, 255, 1, 127, 1, 0, 0, 1, 0, 0},
	[10]uint8{128, 0, 1, 129, 0, 0, 0, 0, 0, 1},
	[10]uint8{128, 1, 1, 130, 0, 0, 0, 0, 0, 1},
	[10]uint8{128, 127, 1, 0, 1, 0, 0, 1, 1, 0},
	[10]uint8{128, 128, 1, 1, 1, 0, 1, 0, 0, 0},
	[10]uint8{128, 129, 1, 2, 1, 0, 1, 0, 0, 0},
	[10]uint8{128, 255, 1, 128, 1, 0, 0, 1, 0, 1},
	[10]uint8{129, 0, 1, 130, 0, 0, 0, 0, 0, 1},
	[10]uint8{129, 1, 1, 131, 0, 0, 0, 0, 0, 1},
	[10]uint8{129, 127, 1, 1, 1, 0, 0, 1, 0, 0},
	[10]uint8{129, 128, 1, 2, 1, 0, 1, 0, 0, 0},
	[10]uint8{129, 129, 1, 3, 1, 0, 1, 0, 0, 0},
	[10]uint8{129, 255, 1, 129, 1, 0, 0, 1, 0, 1},
	[10]uint8{255, 0, 1, 0, 1, 0, 0, 1, 1, 0},
	[10]uint8{255, 1, 1, 1, 1, 0, 0, 1, 0, 0},
	[10]uint8{255, 127, 1, 127, 1, 0, 0, 1, 0, 0},
	[10]uint8{255, 128, 1, 128, 1, 0, 0, 1, 0, 1},
	[10]uint8{255, 129, 1, 129, 1, 0, 0, 1, 0, 1},
	[10]uint8{255, 255, 1, 255, 1, 0, 0, 1, 0, 1},
}

func TestAdcRegister(t *testing.T) {
	var mem = memory.MemoryNew()
	var dmaX = dma.DMANew(mem)
	var cpu = CPUNew(dmaX)

	for _, row := range adcTruthTable {
		for _, register := range [7]byte{'B', 'C', 'D', 'E', 'H', 'L', 'A'} {
			if register == 'A' {
				if row[0] != row[1] {
					continue
				}
			}

			cpu.PC = 0
			cpu.AF = uint16(row[0]) << 8
			cpu.BC = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.DE = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.HL = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.Flags.C = row[2] == 1
			tstates := cpu.adcAR(register)()

			if uint8(cpu.AF>>8) != row[3] || cpu.Flags.C != (row[4] == 1) || cpu.Flags.N != (row[5] == 1) || cpu.Flags.PV != (row[6] == 1) || cpu.Flags.H != (row[7] == 1) || cpu.Flags.Z != (row[8] == 1) || cpu.Flags.S != (row[9] == 1) {
				t.Errorf(
					"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d + %d + %d)",
					uint8(cpu.AF>>8), cpu.Flags.C, cpu.Flags.N, cpu.Flags.PV, cpu.Flags.H, cpu.Flags.Z, cpu.Flags.S,
					row[3], row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[9] == 1, row[0], row[1], row[2],
				)
			}

			if cpu.PC != 1 || tstates != 4 {
				t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, tstates, 1, 4)
			}
		}
	}
}

func TestAdc_Hl_(t *testing.T) {
	var mem = memory.MemoryNew()
	var dmaX = dma.DMANew(mem)
	var cpu = CPUNew(dmaX)
	cpu.HL = 0x1234

	for _, row := range adcTruthTable {
		cpu.PC = 0
		cpu.AF = uint16(row[0]) << 8
		cpu.Flags.C = row[2] == 1
		dmaX.SetMemoryByte(cpu.HL, row[1])
		tstates := cpu.adcA_Hl_()

		if uint8(cpu.AF>>8) != row[3] || cpu.Flags.C != (row[4] == 1) || cpu.Flags.N != (row[5] == 1) || cpu.Flags.PV != (row[6] == 1) || cpu.Flags.H != (row[7] == 1) || cpu.Flags.Z != (row[8] == 1) || cpu.Flags.S != (row[9] == 1) {
			t.Errorf(
				"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d + %d + %d)",
				uint8(cpu.AF>>8), cpu.Flags.C, cpu.Flags.N, cpu.Flags.PV, cpu.Flags.H, cpu.Flags.Z, cpu.Flags.S,
				row[3], row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[9] == 1, row[0], row[1], row[2],
			)
		}

		if cpu.PC != 1 || tstates != 7 {
			t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, tstates, 1, 7)
		}
	}
}
