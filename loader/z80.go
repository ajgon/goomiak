package loader

import (
	"bufio"
	"fmt"
	"os"
)

type Snapshot struct {
	AF     uint16
	AF_    uint16
	BC     uint16
	BC_    uint16
	DE     uint16
	DE_    uint16
	HL     uint16
	HL_    uint16
	IFF1   bool
	IFF2   bool
	IM     uint8
	I      uint8
	IX     uint16
	IY     uint16
	R      uint8
	PC     uint16
	SP     uint16
	Memory []byte
}

func loadFile(filePath string) []byte {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic("error loading file!")
	}
	stat, _ := os.Stat(filePath)

	bytes := make([]byte, stat.Size())
	buf := bufio.NewReader(file)
	buf.Read(bytes)

	return bytes
}

func extractMemoryBlock(data []byte, fileOffset uint, isCompressed bool, unpackedLength uint) []byte {
	if !isCompressed {
		return data[fileOffset:(fileOffset + unpackedLength)]
	}

	var filePtr, memoryPtr uint

	fileBytes := data[fileOffset:]
	memoryBytes := make([]byte, unpackedLength)

	for memoryPtr < unpackedLength {
		if unpackedLength-memoryPtr >= 2 && fileBytes[filePtr] == 0xed && fileBytes[filePtr+1] == 0xed {
			count := uint8(fileBytes[filePtr+2])
			value := uint8(fileBytes[filePtr+3])
			for i := uint8(0); i < count; i++ {
				memoryBytes[memoryPtr] = value
				memoryPtr++
			}
			filePtr += 4
		} else {
			memoryBytes[memoryPtr] = fileBytes[filePtr]
			memoryPtr++
			filePtr++
		}
	}

	return memoryBytes
}

func Z80(filePath string) Snapshot {
	var header [30]uint16

	data := loadFile(filePath)

	for i := 0; i < 30; i++ {
		header[i] = uint16(data[i])
	}

	snapshot := Snapshot{
		AF:   header[0]*256 + header[1],
		AF_:  header[21]*256 + header[22],
		BC:   header[3]*256 + header[2],
		BC_:  header[16]*256 + header[15],
		DE:   header[14]*256 + header[13],
		DE_:  header[18]*256 + header[17],
		HL:   header[5]*256 + header[4],
		HL_:  header[20]*256 + header[19],
		IFF1: header[27] != 0,
		IFF2: header[28] != 0,
		IM:   uint8(header[29] & 0x03),
		I:    uint8(header[10]),
		IX:   header[26]*256 + header[25],
		IY:   header[24]*256 + header[23],
		R:    uint8(header[11]&0x7f | ((header[12] & 0x01) << 7)),
		PC:   header[7]*256 + header[6],
		SP:   header[9]*256 + header[8],
	}

	fmt.Printf("%+v\n", snapshot)

	if snapshot.PC != 0 {
		// version 1
		snapshot.Memory = extractMemoryBlock(data, 30, data[12]&0x20 == 0x20, 0xc000)

	} else {
		panic("Not implemented")
	}

	return snapshot
}
