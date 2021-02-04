package loader

type TapFile struct {
	blocks       [][]uint8
	currentBlock int
}

func (tf *TapFile) Loaded() bool {
	return len(tf.blocks) > 0
}

func (tf *TapFile) NextBlock() []uint8 {
	if tf.currentBlock >= len(tf.blocks) {
		return []uint8{}
	}

	tf.currentBlock++

	return tf.blocks[tf.currentBlock-1]
}

func (tf *TapFile) Rewind() {
	tf.currentBlock = 0
}

func NewTapFile(filePath string) *TapFile {
	var blockSize uint16
	var i int

	data := loadFile(filePath)

	tapFile := &TapFile{blocks: make([][]uint8, 0)}

	for i < len(data) {
		blockSize = (uint16(data[i+1]) << 8) | uint16(data[i])
		i += 2
		tapFile.blocks = append(tapFile.blocks, data[i:(i+int(blockSize))])
		i += int(blockSize)
	}

	return tapFile
}
