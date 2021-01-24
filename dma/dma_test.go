package dma

import (
	"reflect"
	"testing"
	"z80/memory"
)

var markRangeAsDirtyCall [2]uint16
var markAsDirtyCall uint16

type DummyHandler struct {
}

func (dh *DummyHandler) CheckAddressDirtiness(address uint16) bool {
	return true
}

func (dh *DummyHandler) IsMemoryDirty() bool {
	return true
}

func (dh *DummyHandler) MarkAsFresh() {
}

func (dh *DummyHandler) MarkAsDirty(address uint16) {
	markAsDirtyCall = address
}

func (dh *DummyHandler) MarkRangeAsDirty(start, end uint16) {
	markRangeAsDirtyCall = [2]uint16{start, end}
}

func (dh *DummyHandler) Name() string {
	return "dummy"
}

var mem = memory.NewMemory()
var dummyHandler = new(DummyHandler)
var dma = DMANew(mem, dummyHandler)

func TestGetHandler(t *testing.T) {
	mem.Clear()
	got := dma.GetHandler("dummy")
	want := dummyHandler

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetMemory(t *testing.T) {
	mem.Clear()
	mem.SetByte(0x4444, 42)

	got := dma.GetMemoryByte(0x4444)
	want := uint8(42)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestSetMemoryByte(t *testing.T) {
	mem.Clear()
	dma.SetMemoryByte(0x4321, 25)

	got1 := mem.GetByte(0x4321)
	want1 := uint8(25)

	if got1 != want1 {
		t.Errorf("got %d, want %d", got1, want1)
	}

	got2 := markAsDirtyCall
	want2 := uint16(0x4321)

	if got2 != want2 {
		t.Errorf("got %x want %x", got2, want2)
	}
}

func TestSetMemoryBulk(t *testing.T) {
	mem.Clear()
	dma.SetMemoryBulk(0x4321, []uint8{15, 19, 23})

	got1 := []uint8{mem.GetByte(0x4321), mem.GetByte(0x4322), mem.GetByte(0x4323)}
	want1 := []uint8{15, 19, 23}

	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("got %v, want %v", got1, want1)
	}

	got2 := markRangeAsDirtyCall
	want2 := [2]uint16{0x4321, 0x4323}

	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("got %v, want %v", got2, want2)
	}
}
