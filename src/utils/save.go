package utils

import (
	"bytes"
	"io"
	"os"
)

func SaveToFile(r io.ReadCloser, filePath string) {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(r)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	_, _ = file.Write(buf.Bytes())
}
