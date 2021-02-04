package loader

import (
	"bufio"
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
	Border uint8
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
