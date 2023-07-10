package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/disintegration/imaging"
	"image"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

func GenerateThumbnail(outName, bgImg, locationDimensions string) {
	// Coordinate to super-impose on. e.g. 200x500
	// locationX, locationY := parseCoordinates(locationDimensions, "x")

	// src := openImage(bgImg)

	// Resize the watermark to fit these dimensions, preserving aspect ratio.
	markFit := resizeImageKeepingAspectRatio(bgImg, locationDimensions)

	// Place the watermark over the background in the location
	// dst := imaging.Paste(src, markFit, image.Pt(locationX, locationY))

	// err := imaging.Save(dst, outName)
	err := imaging.Save(markFit, outName)

	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

func resizeImageKeepingAspectRatio(image, dimensions string) image.Image {

	width, height := parseCoordinates(dimensions, "x")
	src := openImage(image)
	return imaging.Fit(src, width, height, imaging.Lanczos)
}

func parseCoordinates(input, delimiter string) (int, int) {

	arr := strings.Split(input, delimiter)

	// convert a string to an int
	x, err := strconv.Atoi(arr[0])

	if err != nil {
		log.Fatalf("failed to parse x coordinate: %v", err)
	}

	y, err := strconv.Atoi(arr[1])

	if err != nil {
		log.Fatalf("failed to parse y coordinate: %v", err)
	}

	return x, y
}

func openImage(name string) image.Image {
	src, err := imaging.Open(name)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	return src
}
