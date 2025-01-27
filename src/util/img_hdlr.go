package util

import (
	"github.com/disintegration/imaging"
	"image"
	"log"
	"strconv"
	"strings"
)

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

/**
 * reference : https://medium.datadriveninvestor.com/build-image-watermark-app-for-branding-or-security-purpose-in-go-80f7ee15003b
 * https://stackoverflow.com/questions/16100023/manipulating-watermark-images-with-go
 */
/*
func GenerateThumbnailWithWatermark(srcImgPath, fileHash string) string {

	// imgUrl := UploadCloudinary(srcImgPath, fileHash)

		resizedSrcImg := resizeImageKeepingAspectRatio(srcImgPath, "512x512")
		markImage := readImage(gloval_consts.WATERMARK_THUMB_PATH) // step 2 ==> read mark image

		// step 3 ==> calculate position in center
		baseBound := resizedSrcImg.Bounds()
		markBound := markImage.Bounds()
		offset := image.Pt(
			(baseBound.Size().X/2)-(markBound.Size().X/2),
			(baseBound.Size().Y/2)-(markBound.Size().Y/2))

		// step 4 ==> put watermark with 50% opacity
		outputImage := image.NewRGBA(baseBound)
		draw.Draw(outputImage, outputImage.Bounds(), resizedSrcImg, image.ZP, draw.Src)
		draw.DrawMask(outputImage, markImage.Bounds().Add(offset), markImage, image.ZP,
			image.NewUniform(color.Alpha{128}), image.ZP, draw.Over)

		return encodeWebp(outputImage, fileHash)
}
*/

/*
func encodeWebp(m image.Image, fileHash string) string {

	var buf bytes.Buffer

	// Encode lossless webp
	if err := webp.Encode(&buf, m, &webp.Options{Lossless: true}); err != nil {
		log.Println(err)
	}
	var thmImgFileName = fmt.Sprintf(gloval_consts.DST_IMG_FILENAME, fileHash)
	if err := ioutil.WriteFile(thmImgFileName, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}
	// return ObjectPrivateUpload(thmImgFileName)
	return ""
}
*/

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
