package util

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"strings"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func openImage(name string) image.Image {
	src, err := imaging.Open(name)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	return src
}

func readImage(fileName string) (image image.Image) {
	baseFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}
	defer func(baseFile *os.File) {
		err := baseFile.Close()
		if err != nil {

		}
	}(baseFile)

	fileNameList := strings.Split(strings.ToLower(fileName), ".")

	switch fileNameList[1] {
	case "jpg":
		image, err = jpeg.Decode(baseFile)
	case "png":
		image, err = png.Decode(baseFile)
	default:
		image, err = png.Decode(baseFile)
		err = fmt.Errorf("invalid file type : %s", fileNameList[1])
	}

	if err != nil {
		log.Fatalf("failed to decode for %s: %s", fileName, err)
	}

	return
}

func writeImage(image image.Image, dstFileName string) (err error) {

	fileOut, err := os.Create(dstFileName)
	if err != nil {
		log.Println("Create", err)
		return
	}
	defer func(fileOut *os.File) {
		err := fileOut.Close()
		if err != nil {
			log.Println("Close", err)
		}
	}(fileOut)

	err = jpeg.Encode(fileOut, image, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		log.Println("jpg encode", err)
	}

	return
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		log.Println(err)
		return nil
	}
	return buf.Bytes()
}
