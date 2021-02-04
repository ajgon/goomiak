package loader

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
		AF:     header[0]*256 + header[1],
		AF_:    header[21]*256 + header[22],
		BC:     header[3]*256 + header[2],
		BC_:    header[16]*256 + header[15],
		DE:     header[14]*256 + header[13],
		DE_:    header[18]*256 + header[17],
		HL:     header[5]*256 + header[4],
		HL_:    header[20]*256 + header[19],
		IFF1:   header[27] != 0,
		IFF2:   header[28] != 0,
		IM:     uint8(header[29] & 0x03),
		I:      uint8(header[10]),
		IX:     header[26]*256 + header[25],
		IY:     header[24]*256 + header[23],
		R:      uint8(header[11]&0x7f | ((header[12] & 0x01) << 7)),
		PC:     header[7]*256 + header[6],
		SP:     header[9]*256 + header[8],
		Border: uint8(header[12]&0x0e) >> 1,
	}

	if snapshot.PC != 0 {
		// version 1
		snapshot.Memory = extractMemoryBlock(data, 30, data[12]&0x20 == 0x20, 0xc000)
	} else {
		var offset uint
		snapshot.Memory = make([]byte, 0xc000)
		additionalHeaderLength := uint16(data[31])*256 + uint16(data[30])
		// @todo there is a lot to handle here, but only if spectrum model is different than 48k
		snapshot.PC = uint16(data[33])*256 + uint16(data[32])
		offset = 32 + uint(additionalHeaderLength)

		pageOffsetMap := map[uint8]int{
			4: 0x8000, 5: 0xc000, 8: 0x4000,
		}

		for offset < uint(len(data)) {
			compressedLength := uint16(data[offset+1])*256 + uint16(data[offset])
			isCompressed := true
			if compressedLength == 0xffff {
				compressedLength = 0x4000
				isCompressed = false
			}
			pageId := data[offset+2]
			if _, ok := pageOffsetMap[pageId]; ok {
				pageData := extractMemoryBlock(data, offset+3, isCompressed, 0x4000)
				for b := range pageData {
					snapshot.Memory[b+pageOffsetMap[pageId]-0x4000] = pageData[b]
				}
				offset += uint(compressedLength) + 3
			}
		}
	}

	return snapshot
}
