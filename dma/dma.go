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

func (dma *DMA) GetMemory(address uint16) uint8 {
	return dma.memory.Get(address)
}

func (dma *DMA) SetMemoryByte(address uint16, value uint8) {
	dma.memory.SetByte(address, value)
	for _, handler := range dma.handlers {
		handler.MarkAsDirty(address)
	}
}

func (dma *DMA) SetMemoryBulk(address uint16, bytes []uint8) {
	dma.memory.SetBulk(address, bytes)

	for _, handler := range dma.handlers {
		handler.MarkRangeAsDirty(address, address+uint16(len(bytes))-1)
	}
}

func DMANew(memory *memory.Memory, handlers ...MemoryHandler) *DMA {
	dma := new(DMA)
	dma.memory = memory
	dma.handlers = make(map[string]MemoryHandler)

	for _, handler := range handlers {
		dma.handlers[handler.Name()] = handler
	}

	return dma
}
