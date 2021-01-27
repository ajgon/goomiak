package cpu

import (
	"testing"
)

var subTruthTable [36][9]uint8 = [36][9]uint8{
	// a, r, a-r, C, N, PV, H, Z, S
	[9]uint8{0, 0, 0, 0, 1, 0, 0, 1, 0},
	[9]uint8{0, 1, 255, 1, 1, 0, 1, 0, 1},
	[9]uint8{0, 127, 129, 1, 1, 0, 1, 0, 1},
	[9]uint8{0, 128, 128, 1, 1, 1, 0, 0, 1},
	[9]uint8{0, 129, 127, 1, 1, 0, 1, 0, 0},
	[9]uint8{0, 255, 1, 1, 1, 0, 1, 0, 0},
	[9]uint8{1, 0, 1, 0, 1, 0, 0, 0, 0},
	[9]uint8{1, 1, 0, 0, 1, 0, 0, 1, 0},
	[9]uint8{1, 127, 130, 1, 1, 0, 1, 0, 1},
	[9]uint8{1, 128, 129, 1, 1, 1, 0, 0, 1},
	[9]uint8{1, 129, 128, 1, 1, 1, 0, 0, 1},
	[9]uint8{1, 255, 2, 1, 1, 0, 1, 0, 0},
	[9]uint8{127, 0, 127, 0, 1, 0, 0, 0, 0},
	[9]uint8{127, 1, 126, 0, 1, 0, 0, 0, 0},
	[9]uint8{127, 127, 0, 0, 1, 0, 0, 1, 0},
	[9]uint8{127, 128, 255, 1, 1, 1, 0, 0, 1},
	[9]uint8{127, 129, 254, 1, 1, 1, 0, 0, 1},
	[9]uint8{127, 255, 128, 1, 1, 1, 0, 0, 1},
	[9]uint8{128, 0, 128, 0, 1, 0, 0, 0, 1},
	[9]uint8{128, 1, 127, 0, 1, 1, 1, 0, 0},
	[9]uint8{128, 127, 1, 0, 1, 1, 1, 0, 0},
	[9]uint8{128, 128, 0, 0, 1, 0, 0, 1, 0},
	[9]uint8{128, 129, 255, 1, 1, 0, 1, 0, 1},
	[9]uint8{128, 255, 129, 1, 1, 0, 1, 0, 1},
	[9]uint8{129, 0, 129, 0, 1, 0, 0, 0, 1},
	[9]uint8{129, 1, 128, 0, 1, 0, 0, 0, 1},
	[9]uint8{129, 127, 2, 0, 1, 1, 1, 0, 0},
	[9]uint8{129, 128, 1, 0, 1, 0, 0, 0, 0},
	[9]uint8{129, 129, 0, 0, 1, 0, 0, 1, 0},
	[9]uint8{129, 255, 130, 1, 1, 0, 1, 0, 1},
	[9]uint8{255, 0, 255, 0, 1, 0, 0, 0, 1},
	[9]uint8{255, 1, 254, 0, 1, 0, 0, 0, 1},
	[9]uint8{255, 127, 128, 0, 1, 0, 0, 0, 1},
	[9]uint8{255, 128, 127, 0, 1, 0, 0, 0, 0},
	[9]uint8{255, 129, 126, 0, 1, 0, 0, 0, 0},
	[9]uint8{255, 255, 0, 0, 1, 0, 0, 1, 0},
}

func TestSubRegister(t *testing.T) {
	for _, row := range subTruthTable {
		for _, register := range [11]byte{'B', 'C', 'D', 'E', 'H', 'L', 'A', 'X', 'x', 'Y', 'y'} {
			adjustPC := uint16(0)

			if register == 'A' {
				if row[0] != row[1] {
					continue
				}
			}

			if register == 'X' || register == 'x' || register == 'Y' || register == 'y' {
				adjustPC = 1
			}

			cpu.PC = 0
			cpu.tstates = 4
			cpu.setAcc(row[0])
			cpu.BC = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.DE = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.HL = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.IX = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.IY = (uint16(row[1]) << 8) | uint16(row[1])
			cpu.subR(register)()

			if cpu.getAcc() != row[2] || cpu.getC() != (row[3] == 1) || cpu.getN() != (row[4] == 1) || cpu.getPV() != (row[5] == 1) || cpu.getH() != (row[6] == 1) || cpu.getZ() != (row[7] == 1) || cpu.getS() != (row[8] == 1) {
				t.Errorf(
					"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d - %d)",
					cpu.getAcc(), cpu.getC(), cpu.getN(), cpu.getPV(), cpu.getH(), cpu.getZ(), cpu.getS(),
					row[2], row[3] == 1, row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[0], row[1],
				)
			}

			if cpu.PC != 1+adjustPC || cpu.tstates != 4 {
				t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, cpu.tstates, 1+adjustPC, 4)
			}
		}
	}
}

func TestSub_Hl_(t *testing.T) {
	cpu.HL = 0x1234

	for _, row := range subTruthTable {
		cpu.PC = 0
		cpu.tstates = 4
		cpu.setAcc(row[0])
		dmaX.SetMemoryByte(cpu.HL, row[1])
		cpu.sub_Ss_("HL")()

		if cpu.getAcc() != row[2] || cpu.getC() != (row[3] == 1) || cpu.getN() != (row[4] == 1) || cpu.getPV() != (row[5] == 1) || cpu.getH() != (row[6] == 1) || cpu.getZ() != (row[7] == 1) || cpu.getS() != (row[8] == 1) {
			t.Errorf(
				"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d + %d)",
				cpu.getAcc(), cpu.getC(), cpu.getN(), cpu.getPV(), cpu.getH(), cpu.getZ(), cpu.getS(),
				row[2], row[3] == 1, row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[0], row[1],
			)
		}

		if cpu.PC != 1 || cpu.tstates != 7 {
			t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, cpu.tstates, 1, 7)
		}
	}
}

func TestSub_Ix_(t *testing.T) {
	cpu.IX = 0x121b

	for _, row := range subTruthTable {
		cpu.PC = 0
		cpu.tstates = 8
		cpu.setAcc(row[0])
		dmaX.SetMemoryByte(0x1234, row[1])
		dmaX.SetMemoryByte(0x0002, 0x19)
		cpu.sub_Ss_("IX")()

		if cpu.getAcc() != row[2] || cpu.getC() != (row[3] == 1) || cpu.getN() != (row[4] == 1) || cpu.getPV() != (row[5] == 1) || cpu.getH() != (row[6] == 1) || cpu.getZ() != (row[7] == 1) || cpu.getS() != (row[8] == 1) {
			t.Errorf(
				"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d + %d)",
				cpu.getAcc(), cpu.getC(), cpu.getN(), cpu.getPV(), cpu.getH(), cpu.getZ(), cpu.getS(),
				row[2], row[3] == 1, row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[0], row[1],
			)
		}

		if cpu.PC != 3 || cpu.tstates != 19 {
			t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, cpu.tstates, 3, 19)
		}
	}
}

func TestSub_Iy_(t *testing.T) {
	cpu.IY = 0x121b

	for _, row := range subTruthTable {
		cpu.PC = 0
		cpu.tstates = 8
		cpu.setAcc(row[0])
		dmaX.SetMemoryByte(0x1234, row[1])
		dmaX.SetMemoryByte(0x0002, 0x19)
		cpu.sub_Ss_("IY")()

		if cpu.getAcc() != row[2] || cpu.getC() != (row[3] == 1) || cpu.getN() != (row[4] == 1) || cpu.getPV() != (row[5] == 1) || cpu.getH() != (row[6] == 1) || cpu.getZ() != (row[7] == 1) || cpu.getS() != (row[8] == 1) {
			t.Errorf(
				"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d + %d)",
				cpu.getAcc(), cpu.getC(), cpu.getN(), cpu.getPV(), cpu.getH(), cpu.getZ(), cpu.getS(),
				row[2], row[3] == 1, row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[0], row[1],
			)
		}

		if cpu.PC != 3 || cpu.tstates != 19 {
			t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, cpu.tstates, 3, 19)
		}
	}
}

func TestSubX(t *testing.T) {
	for _, row := range subTruthTable {
		cpu.PC = 0
		cpu.tstates = 4
		cpu.setAcc(row[0])
		dmaX.SetMemoryByte(0x0001, row[1])
		cpu.subN()

		if cpu.getAcc() != row[2] || cpu.getC() != (row[3] == 1) || cpu.getN() != (row[4] == 1) || cpu.getPV() != (row[5] == 1) || cpu.getH() != (row[6] == 1) || cpu.getZ() != (row[7] == 1) || cpu.getS() != (row[8] == 1) {
			t.Errorf(
				"\ngot:  A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t\nwant: A=0x%02x, C=%t, N=%t, PV=%t, H=%t, Z=%t, S=%t for (%d + %d)",
				cpu.getAcc(), cpu.getC(), cpu.getN(), cpu.getPV(), cpu.getH(), cpu.getZ(), cpu.getS(),
				row[2], row[3] == 1, row[4] == 1, row[5] == 1, row[6] == 1, row[7] == 1, row[8] == 1, row[0], row[1],
			)
		}

		if cpu.PC != 2 || cpu.tstates != 7 {
			t.Errorf("got PC=%d, %d T-states, want PC=%d, %d T-states", cpu.PC, cpu.tstates, 2, 7)
		}
	}
}
