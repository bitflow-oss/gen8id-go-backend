package utils

// Import the required packages for upload and admin.

// Import the required packages for upload and admin.

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"log"
)

func clodinaryCaptioning(srcImgPath string) {

	var ctx = context.Background()

	cld, _ := cloudinary.NewFromParams("dbqltwqac", "696441825692637", "YRXvf_DaiwJ7tGpLcAaXxyRgXxo")

	// Detection:   "openimages", "captioning"
	resp, err := cld.Upload.Upload(ctx, srcImgPath, uploader.UploadParams{
		Detection: "openimages", AutoTagging: 0.6})

	//resp, err := cld.Upload.Upload(ctx, srcImgPath, uploader.UploadParams{
	//	Detection: "openimages", AutoTagging: 0.6})

	if err != nil {
		log.Fatalf("Failed to upload file, %v\n", err)
	}

	fmt.Printf("resp %+v", resp)
	// log.Println("SecureURL", resp.SecureURL)

}
