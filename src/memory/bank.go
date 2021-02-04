package memory

const memoryBankSize uint16 = 0x4000

type MemoryBank struct {
	bytes     [memoryBankSize]uint8
	contended bool
	readOnly  bool
}

func (mb *MemoryBank) loadData(data []byte, startAddress uint16) {
	var endAddress uint16

	if startAddress+uint16(len(data)) > memoryBankSize {
		endAddress = memoryBankSize
	} else {
		endAddress = startAddress + uint16(len(data))
	}

	for addr := startAddress; addr < endAddress; addr++ {
		mb.bytes[addr] = data[addr-startAddress]
	}
}

func (mb *MemoryBank) GetByte(address uint16) uint8 {
	return mb.bytes[address]
}

func (mb *MemoryBank) SetByte(address uint16, value uint8) {
	if mb.readOnly {
		return
	}

	mb.bytes[address] = value
}

func (mb *MemoryBank) Clear() {
	mb.bytes = [memoryBankSize]uint8{}
}

func NewMemoryBank(contented bool, readOnly bool) *MemoryBank {
	return &MemoryBank{contended: contented, readOnly: readOnly}
}
