package memory

import "testing"

func TestBankGetByte(t *testing.T) {
	bank := NewMemoryBank(false, false)
	bank.bytes[0x1234] = 42

	got := bank.GetByte(0x1234)
	want := uint8(42)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestBankSetByte(t *testing.T) {
	bank := NewMemoryBank(false, false)
	bank.SetByte(0x2345, 55)

	got := bank.bytes[0x2345]
	want := uint8(55)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	bank = NewMemoryBank(false, true)
	bank.bytes[0x2345] = 77
	bank.SetByte(0x2345, 66)

	got = bank.bytes[0x2345]
	want = uint8(77)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
