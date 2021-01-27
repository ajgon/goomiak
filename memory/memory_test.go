package memory

import (
	"testing"
)

func TestMemoryBankOptions(t *testing.T) {
	memory := NewMemory()

	for i := 0; i < 4; i++ {
		gotContended, gotReadOnly := memory.roms[i].contended, memory.roms[i].readOnly
		wantContended, wantReadOnly := false, true

		if gotContended != wantContended || gotReadOnly != wantReadOnly {
			t.Errorf("for rom %d, got %t/%t, want %t/%t", i, gotContended, gotReadOnly, wantContended, wantReadOnly)
		}
	}

	for i := 0; i < 8; i++ {
		gotContended, gotReadOnly := memory.banks[i].contended, memory.banks[i].readOnly
		wantContended, wantReadOnly := i&1 == 1, false

		if gotContended != wantContended || gotReadOnly != wantReadOnly {
			t.Errorf("for bank %d, got %t/%t, want %t/%t", i, gotContended, gotReadOnly, wantContended, wantReadOnly)
		}
	}
}

func TestMemoryGetByte(t *testing.T) {
	memory := NewMemory()
	memory.activeBanks[0].bytes[0x1234] = 42

	got, _ := memory.GetByte(0x1234)
	want := uint8(42)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	memory.activeBanks[2].bytes[0x2345] = 77

	got, _ = memory.GetByte(0xa345)
	want = uint8(77)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestMemorySetByte(t *testing.T) {
	memory := NewMemory()
	memory.SetByte(0x1234, 55)

	// first bank is a ROM
	got := memory.activeBanks[0].bytes[0x1234]
	want := uint8(0)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	memory.SetByte(0x6345, 55)

	got = memory.activeBanks[1].bytes[0x2345]
	want = uint8(55)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
