package dma

import (
	"z80/memory"
)

type MemoryHandler interface {
	CheckAddressDirtiness(uint16) bool
	IsMemoryDirty() bool
	MarkAsFresh()
	MarkAsDirty(address uint16)
	MarkRangeAsDirty(start, end uint16)
	Name() string
}

type DMA struct {
	memory   *memory.Memory
	handlers map[string]MemoryHandler
}

func (dma *DMA) GetHandler(name string) MemoryHandler {
	return dma.handlers[name]
}

func (dma *DMA) GetMemoryByte(address uint16) (uint8, bool) {
	return dma.memory.GetByte(address)
}

func (dma *DMA) SetMemoryByte(address uint16, value uint8) (contended bool) {
	contended = dma.memory.SetByte(address, value)

	for _, handler := range dma.handlers {
		handler.MarkAsDirty(address)
	}

	return
}

// this function is used for testing only, shouldn't be used in production code
func (dma *DMA) SetMemoryBulk(address uint16, bytes []uint8) {
	for i := uint32(address); i < uint32(address)+uint32(len(bytes)); i++ {
		dma.memory.SetByte(uint16(i), bytes[uint16(i-uint32(address))])
	}

	for _, handler := range dma.handlers {
		handler.MarkRangeAsDirty(address, address+uint16(len(bytes))-1)
	}
}

func (dma *DMA) LoadData(startAddress uint16, data []byte) {
	dma.memory.LoadData(uint32(startAddress), data)
}

func NewDMA(memory *memory.Memory, handlers ...MemoryHandler) *DMA {
	dma := new(DMA)
	dma.memory = memory
	dma.handlers = make(map[string]MemoryHandler)

	for _, handler := range handlers {
		dma.handlers[handler.Name()] = handler
	}

	return dma
}
