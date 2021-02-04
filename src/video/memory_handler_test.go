package video

import (
	"reflect"
	"testing"
)

var vmh = VideoMemoryHandlerNew()

func TestMarkAsDirty(t *testing.T) {
	vmh.MarkAsFresh()
	vmh.MarkAsDirty(0x5900)

	got := vmh.dirtyAddresses[0x1900]
	want := true

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}

func TestMarkAsDirtyWithAddressesOutside(t *testing.T) {
	vmh.MarkAsFresh()
	vmh.MarkAsDirty(0x3000)

	for i := 0x0000; i < 0x1b00; i++ {
		if vmh.dirtyAddresses[i] != false {
			t.Errorf("got %t (at %x), want %t", vmh.dirtyAddresses[i], i, false)
		}
	}
}

func TestMarkRangeAsDirty(t *testing.T) {
	vmh.MarkAsFresh()
	vmh.MarkRangeAsDirty(0x4500, 0x4502)

	got := []bool{vmh.dirtyAddresses[0x04ff], vmh.dirtyAddresses[0x500], vmh.dirtyAddresses[0x501], vmh.dirtyAddresses[0x502], vmh.dirtyAddresses[0x503]}
	want := []bool{false, true, true, true, false}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMarkRangeAsDirtyWithAddressesOutside(t *testing.T) {
	vmh.MarkAsFresh()
	vmh.MarkRangeAsDirty(0x3000, 0x4000)

	got := vmh.dirtyAddresses[0x0000]
	want := true

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}

	for i := 0x0001; i < 0x1b00; i++ {
		if vmh.dirtyAddresses[i] != false {
			t.Errorf("got %t (at %x), want %t", vmh.dirtyAddresses[i], i, false)
		}
	}
}

func TestMarkAsFresh(t *testing.T) {
	vmh.dirtyAddresses[0x0100] = true
	vmh.MarkAsFresh()

	got := vmh.dirtyAddresses[0x0100]
	want := false

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}
