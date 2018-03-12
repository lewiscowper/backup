package checksum

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"
)

func computeSha512(filePath string) (string, error) {
	var result string
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	h := sha512.New()
	if _, err := io.Copy(h, file); err != nil {
		return result, err
	}

	checksum := h.Sum(nil)

	result = hex.EncodeToString(checksum)

	return result, nil
}

// CreateFile creates the checksum file in the location specified
func CreateFile(filePath string, filename string) error {
	checksumFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer checksumFile.Close()

	hash, err := computeSha512(filePath)
	if err != nil {
		return err
	}

	checksumFile.WriteString(hash)

	return nil
}
