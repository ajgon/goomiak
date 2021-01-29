package video

import (
	"z80/dma"
)

const fullWidth uint = 256 + 48 + 48
const fullHeight uint = 256 + 48 + 48
const screenWidth uint = 256
const screenHeight uint = 192

const borderLeftWidth uint = 48
const borderRightWidth uint = 48
const borderTopHeight uint = 48
const borderBottomHeight uint = 56

const borderRightPosition uint = 48 + 256
const borderBottomPosition uint = 48 + 192

const pixelsAddress uint16 = 0x4000
const attributesAddress uint16 = 0x5800
const pixelsByteCount uint = fullWidth * fullHeight

type PixelRenderer struct {
	dma           *dma.DMA
	pixels        [pixelsByteCount]byte
	memoryHandler dma.MemoryHandler

	screenPixelAddresses [screenHeight][screenWidth]uint16
}

func (pr *PixelRenderer) prepareScreenPixelAddresses() {
	var x, y uint
	var highAddr, lowAddr uint16

	for y = 0; y < screenHeight; y++ {
		for x = 0; x < screenWidth/8; x++ {
			highAddr = uint16(((y & 0b00000111) | 0b01000000) | ((y >> 3) & 0b00011000))
			lowAddr = uint16(((y << 2) & 0b11100000) | (x & 0b00011111))

			pr.screenPixelAddresses[y][x] = (highAddr << 8) | lowAddr
		}
	}
}

func (pr *PixelRenderer) PaintPixel(x, y uint) {
	var pixelPosition uint

	// @todo check port for border color, also paint only if color changes
	if x < borderLeftWidth || x >= borderRightPosition || y < borderTopHeight || y >= borderBottomPosition {
		pixelPosition = 4 * (fullWidth*y + x)

		pr.pixels[pixelPosition] = 207
		pr.pixels[pixelPosition+1] = 207
		pr.pixels[pixelPosition+2] = 207
		pr.pixels[pixelPosition+3] = 207

		return
	}

	var color uint8

	xIndex := (x - 48) / 8
	yIndex := (y - 48)

	pixelAddress := uint16(pr.screenPixelAddresses[y][xIndex])
	colorAddress := attributesAddress + uint16(yIndex*4+xIndex)

	if !pr.memoryHandler.CheckAddressDirtiness(pixelAddress) && !pr.memoryHandler.CheckAddressDirtiness(colorAddress) {
		return
	}

	value, _ := pr.dma.GetMemoryByte(pixelAddress)
	colorValue, _ := pr.dma.GetMemoryByte(colorAddress)

	if value&(128>>(x%8)) == (128 >> (x % 8)) {
		color = uint8(((colorValue >> 3) & 0b00001000) | (colorValue & 0b00000111)) // ink
	} else {
		color = uint8((colorValue >> 3) & 0b00001111) // paper
	}

	pixelPosition = 4 * (fullWidth*y + x)

	brightness := (color&0b00001000)*6 + 207

	if color&0b00000010 == 0b00000010 {
		pr.pixels[pixelPosition] = brightness
	} else {
		pr.pixels[pixelPosition] = 0
	}
	if color&0b00000100 == 0b00000100 {
		pr.pixels[pixelPosition+1] = brightness
	} else {
		pr.pixels[pixelPosition+1] = 0
	}
	if color&0b00000001 == 0b00000001 {
		pr.pixels[pixelPosition+2] = brightness
	} else {
		pr.pixels[pixelPosition+2] = 0
	}
	pr.pixels[pixelPosition+3] = 255
}

func (pr *PixelRenderer) Pixels() []byte {
	return pr.pixels[:]
}

func NewPixelRenderer(dma *dma.DMA) *PixelRenderer {
	pixelRenderer := new(PixelRenderer)
	pixelRenderer.dma = dma
	pixelRenderer.memoryHandler = dma.GetHandler("video")

	pixelRenderer.prepareScreenPixelAddresses()

	return pixelRenderer
}
