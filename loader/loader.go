package loader

import (
	"bufio"
	"os"
)

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
