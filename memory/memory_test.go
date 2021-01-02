package memory

import (
	"reflect"
	"testing"
)

var memory = MemoryNew()

func TestGet(t *testing.T) {
	memory.Clear()
	memory.bytes[0x1234] = 42

	got := memory.Get(0x1234)
	want := uint8(42)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestSetByte(t *testing.T) {
	memory.Clear()
	memory.SetByte(0x4321, 55)

	got := memory.bytes[0x4321]
	want := uint8(55)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestSetBulk(t *testing.T) {
	memory.Clear()
	memory.SetBulk(0x4321, []uint8{58, 68, 78})

	got1 := []uint8{memory.bytes[0x4321], memory.bytes[0x4322], memory.bytes[0x4323]}
	want1 := []uint8{58, 68, 78}

	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("got %v, want %v", got1, want1)
	}
}
