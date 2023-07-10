package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func ExtractFileHash(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}(f)
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
