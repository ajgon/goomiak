package video

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"reflect"
	"testing"
	"z80/dma"
	"z80/memory"
)

var mem = memory.MemoryNew()
var videoMemoryHandler = VideoMemoryHandlerNew()
var dmaAddresser = dma.DMANew(mem, videoMemoryHandler)
var video = VideoNew(dmaAddresser)

func loadFileToMemory(dma *dma.DMA, address uint16, filePath string) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic("error loading file!")
	}

	bytes := make([]byte, 6912)
	buf := bufio.NewReader(file)
	buf.Read(bytes)

	dma.SetMemoryBulk(address, bytes)
}

func TestBuildPixel(t *testing.T) {
	pixelsMap := map[uint8][]byte{
		0:  []byte{0, 0, 0, 255},
		1:  []byte{0, 0, 207, 255},
		2:  []byte{207, 0, 0, 255},
		3:  []byte{207, 0, 207, 255},
		4:  []byte{0, 207, 0, 255},
		5:  []byte{0, 207, 207, 255},
		6:  []byte{207, 207, 0, 255},
		7:  []byte{207, 207, 207, 255},
		8:  []byte{0, 0, 0, 255},
		9:  []byte{0, 0, 255, 255},
		10: []byte{255, 0, 0, 255},
		11: []byte{255, 0, 255, 255},
		12: []byte{0, 255, 0, 255},
		13: []byte{0, 255, 255, 255},
		14: []byte{255, 255, 0, 255},
		15: []byte{255, 255, 255, 255},
	}

	for i := uint8(0); i < 16; i++ {
		video.pixels = [pixelsByteCount]byte{}
		video.buildPixel(4, 8, i)
		got := video.pixels[8208:8212]
		want := pixelsMap[i]

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v (%d), want %v", got, i, want)
		}
	}
}

func TestGetPixelAddress(t *testing.T) {
	got := getPixelAddress(14, 179)
	want := uint16(0x53ce)

	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestNeedsRefresh(t *testing.T) {
	mem.Clear()
	got := video.NeedsRefresh()
	want := false

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}

	dmaAddresser.SetMemoryBulk(0x4000, []uint8{10})

	got = video.NeedsRefresh()
	want = true

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}

func TestPixels(t *testing.T) {
	mem.Clear()
	loadFileToMemory(dmaAddresser, 0x4000, "./example.scr")

	shasum := sha256.Sum256(video.Pixels())
	got := hex.EncodeToString(shasum[:])
	want := "8c409441e6797652a172638610b987503d7849828413f46b3befb50dc25c39b5"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func BenchmarkPixels(b *testing.B) {
	mem.Clear()
	loadFileToMemory(dmaAddresser, 0x4000, "./example.scr")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		videoMemoryHandler.isDirty = true
		video.Pixels()
	}
}
