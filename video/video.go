package video

import (
	"z80/dma"
)

const fullWidth uint16 = 256 + 48 + 48
const fullHeight uint16 = 192 + 48 + 56
const screenWidth uint16 = 256
const screenHeight uint16 = 192
const pixelsAddress uint16 = 0x4000
const colorsAddress uint16 = 0x5800
const pixelsByteCount uint32 = uint32(fullWidth) * uint32(fullHeight) * 4

type Video struct {
	dma           *dma.DMA
	lineAddresses [screenHeight]uint16
	pixels        [pixelsByteCount]byte
	memoryHandler dma.MemoryHandler
}

func getPixelAddress(x, y uint8) uint16 {
	h := ((y & 0b00000111) | 0b01000000) | ((y >> 3) & 0b00011000)
	l := ((y << 2) & 0b11100000) | (x & 0b00011111)

	return (uint16(h) << 8) | uint16(l)
}

func (v *Video) prepareLineAddresses() {
	for y := uint16(0); y < screenHeight; y++ {
		v.lineAddresses[y] = getPixelAddress(uint8(0), uint8(y))
	}
}

func (v *Video) PaintPixel(x, y uint64) {
	if x < 48 || x >= 304 || y < 48 || y >= 240 {
		addrr := 4 * uint32(uint32(fullWidth)*uint32(y)+uint32(x))
		//fmt.Println(x, y, addrr)
		v.pixels[addrr] = 207
		v.pixels[addrr+1] = 207
		v.pixels[addrr+2] = 207
		v.pixels[addrr+3] = 207
		return
	}
	var color uint8

	y -= 48
	x -= 48

	pixelAddress := uint16(v.lineAddresses[y])
	colorAddress := colorsAddress + (uint16(y)/8)*32 + uint16(x/8)
	value, _ := v.dma.GetMemoryByte(pixelAddress + uint16(x/8))
	colorValue, _ := v.dma.GetMemoryByte(colorAddress)
	ink := uint8(((colorValue >> 3) & 0b00001000) | (colorValue & 0b00000111))
	paper := uint8((colorValue >> 3) & 0b00001111)

	if value&(128>>(x%8)) == (128 >> (x % 8)) {
		color = ink
	} else {
		color = paper
	}

	addr := 4 * uint32(uint32(fullWidth)*uint32(y+48)+uint32(x+48))
	brightness := (color&0b00001000)*6 + 207

	if color&0b00000010 == 0b00000010 {
		v.pixels[addr] = brightness
	} else {
		v.pixels[addr] = 0
	}
	if color&0b00000100 == 0b00000100 {
		v.pixels[addr+1] = brightness
	} else {
		v.pixels[addr+1] = 0
	}
	if color&0b00000001 == 0b00000001 {
		v.pixels[addr+2] = brightness
	} else {
		v.pixels[addr+2] = 0
	}
	v.pixels[addr+3] = 255
}

func (v *Video) buildPixel(x, y uint16, color uint8) {
	addr := 4 * uint32(screenWidth*y+x)
	brightness := (color&0b00001000)*6 + 207

	if color&0b00000010 == 0b00000010 {
		v.pixels[addr] = brightness
	}
	if color&0b00000100 == 0b00000100 {
		v.pixels[addr+1] = brightness
	}
	if color&0b00000001 == 0b00000001 {
		v.pixels[addr+2] = brightness
	}
	v.pixels[addr+3] = 255
}

func (v *Video) NeedsRefresh() bool {
	return v.memoryHandler.IsMemoryDirty()
}

func (v *Video) AllPixels() []byte {
	return v.pixels[:]
}

func (v *Video) Pixels() []byte {
	if !v.memoryHandler.IsMemoryDirty() {
		return v.pixels[:]
	}

	memScreenWidth := screenWidth / 8

	for y := uint16(0); y < screenHeight; y++ {
		address := uint16(v.lineAddresses[y])
		for x := uint16(0); x < memScreenWidth; x++ {
			colorAddress := colorsAddress + (y/8)*32 + x
			if !v.memoryHandler.CheckAddressDirtiness(address+x) && !v.memoryHandler.CheckAddressDirtiness(colorAddress) {
				continue
			}
			value, _ := v.dma.GetMemoryByte(address + x)
			colorValue, _ := v.dma.GetMemoryByte(colorAddress)
			ink := uint8(((colorValue >> 3) & 0b00001000) | (colorValue & 0b00000111))
			paper := uint8((colorValue >> 3) & 0b00001111)

			if value&128 == 128 {
				v.buildPixel(x*8, y, ink)
			} else {
				v.buildPixel(x*8, y, paper)
			}

			if value&64 == 64 {
				v.buildPixel(x*8+1, y, ink)
			} else {
				v.buildPixel(x*8+1, y, paper)
			}

			if value&32 == 32 {
				v.buildPixel(x*8+2, y, ink)
			} else {
				v.buildPixel(x*8+2, y, paper)
			}

			if value&16 == 16 {
				v.buildPixel(x*8+3, y, ink)
			} else {
				v.buildPixel(x*8+3, y, paper)
			}

			if value&8 == 8 {
				v.buildPixel(x*8+4, y, ink)
			} else {
				v.buildPixel(x*8+4, y, paper)
			}

			if value&4 == 4 {
				v.buildPixel(x*8+5, y, ink)
			} else {
				v.buildPixel(x*8+5, y, paper)
			}

			if value&2 == 2 {
				v.buildPixel(x*8+6, y, ink)
			} else {
				v.buildPixel(x*8+6, y, paper)
			}

			if value&1 == 1 {
				v.buildPixel(x*8+7, y, ink)
			} else {
				v.buildPixel(x*8+7, y, paper)
			}
		}
	}

	v.memoryHandler.MarkAsFresh()
	return v.pixels[:]
}

func VideoNew(dma *dma.DMA) *Video {
	video := new(Video)
	video.dma = dma
	video.prepareLineAddresses()
	video.memoryHandler = dma.GetHandler("video")
	return video
}
