package memory

import (
	"z80/loader"
)

type Memory struct {
	roms        [4]MemoryBank
	banks       [8]MemoryBank
	activeBanks [4]*MemoryBank
}

func (m *Memory) GetByte(address uint16) (uint8, bool) {
	bankNumber := uint8(address >> 14)
	return m.activeBanks[bankNumber].GetByte(address & 0x3fff), m.activeBanks[bankNumber].contended
}

func (m *Memory) SetByte(address uint16, value uint8) bool {
	bankNumber := uint8(address >> 14)
	m.activeBanks[bankNumber].SetByte(address&0x3fff, value)
	return m.activeBanks[bankNumber].contended
}

func (m *Memory) Clear() {
	for i := 0; i < 8; i++ {
		m.banks[i].Clear()
	}
}

func (m *Memory) LoadData(startAddress uint32, data []byte) {
	var endAddress uint32
	memoryLimit := 4 * uint32(memoryBankSize)

	if startAddress > memoryLimit {
		startAddress = memoryLimit
	}

	if startAddress+uint32(len(data)) > memoryLimit {
		endAddress = memoryLimit
	} else {
		endAddress = startAddress + uint32(len(data))
	}

	for address := startAddress; address < endAddress; address++ {
		m.activeBanks[address>>14].bytes[address&0x3fff] = data[address-startAddress]
	}
}

func NewMemory() *Memory {
	memory := new(Memory)
	for i := 0; i < 4; i++ {
		memory.roms[i] = *NewMemoryBank(false, true)
	}

	for i := 0; i < 8; i++ {
		// Odd memory banks are always contended
		memory.banks[i] = *NewMemoryBank(i%2 == 1, false)
	}

	memory.activeBanks = [4]*MemoryBank{
		&memory.roms[0], &memory.banks[5], &memory.banks[2], &memory.banks[0],
	}

	return memory
}

// only for testing, never used in real-case scenario
func NewWritableMemory() *Memory {
	memory := new(Memory)
	for i := 0; i < 4; i++ {
		memory.roms[i] = *NewMemoryBank(false, false)
	}

	for i := 0; i < 8; i++ {
		// Odd memory banks are always contended
		memory.banks[i] = *NewMemoryBank(i%2 == 1, false)
	}

	memory.activeBanks = [4]*MemoryBank{
		&memory.roms[0], &memory.banks[5], &memory.banks[2], &memory.banks[0],
	}

	return memory
}

func (m *Memory) LoadSnapshot(snapshot loader.Snapshot) {
	for ptr := 0; ptr < len(snapshot.Memory); ptr++ {
		m.SetByte(uint16(0x4000+ptr), snapshot.Memory[ptr])
	}
}
