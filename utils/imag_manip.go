package utils

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

const watermarkImgPath = "resources/watermark-pattern-gen8id.png"

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

/**
 * reference : https://medium.datadriveninvestor.com/build-image-watermark-app-for-branding-or-security-purpose-in-go-80f7ee15003b
 * https://stackoverflow.com/questions/16100023/manipulating-watermark-images-with-go
 */
func GenerateThumbnailWithWatermark(srcImgPath, dstImgPath string) {
	
	// clodinaryCaptioning(srcImgPath)

	resizedSrcImg := resizeImageKeepingAspectRatio(srcImgPath, "512x512")
	markImage := readImage(watermarkImgPath) // step 2 ==> read mark image

	// step 3 ==> calculate position in center
	baseBound := resizedSrcImg.Bounds()
	markBound := markImage.Bounds()
	offset := image.Pt(
		(baseBound.Size().X/2)-(markBound.Size().X/2),
		(baseBound.Size().Y/2)-(markBound.Size().Y/2))
	//

	// step 4 ==> put watermark with 50% opacity
	outputImage := image.NewRGBA(baseBound)
	draw.Draw(outputImage, outputImage.Bounds(), resizedSrcImg, image.ZP, draw.Src)
	draw.DrawMask(outputImage, markImage.Bounds().Add(offset), markImage, image.ZP,
		image.NewUniform(color.Alpha{128}), image.ZP, draw.Over)

	err := writeImage(outputImage, dstImgPath) // step 5 ==> write output to file image
	if err != nil {
		log.Println(err)
	}
	// updateExifMeta(dstImgPath)
	return
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

	/*
		exifProps := map[string]interface{}{
			"Artist":    "https://gen8.id",
			"Copyright": "https://gen8.id",
		}
		err = update.UpdateExif(fileOut, os.Stdout, exifProps)
		if err != nil {
			log.Println("UpdateExif", err)
			return
		}
	*/

	err = jpeg.Encode(fileOut, image, &jpeg.Options{Quality: 50}) // jpeg.DefaultQuality
	if err != nil {
		log.Println("jpg encode", err)
	}

	/*
		fh, _ := os.Open(dstFileName)
		defer func(fh *os.File) {
			err := fh.Close()
			if err != nil {
				log.Println("close file", err)
			}
		}(fh)

		err = imaging.Save(image, dstFileName)
		if err != nil {
			log.Println("Save", err)
			return
		}
	*/
	return
}

/**
 * from: https://ahmadrosid.com/blog/golang-img-crop
 */
func cropImage(srcImg, watermarkImg image.Image) image.Image {
	bounds := srcImg.Bounds()
	width := bounds.Dx()
	// height := bounds.Dy() you can use this to work with the height of the images
	cropSize := image.Rect(0, 0, width/2+100, width/2+100)
	cropSize = cropSize.Add(image.Point{100, 80})
	return srcImg.(SubImager).SubImage(cropSize)
}

func GenerateThumbnail(outName, bgImg string) {

	// Resize the watermark to fit these dimensions, preserving aspect ratio.
	resized := resizeImageKeepingAspectRatio(bgImg, "512x512")

	// err := imaging.Save(dst, outName)
	err := imaging.Save(resized, outName)

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
