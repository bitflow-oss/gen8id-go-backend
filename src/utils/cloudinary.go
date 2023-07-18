package utils

// Import the required packages for upload and admin.
import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"log"
)

func UploadCloudinary(srcImgPath, fileHash string) string {

	var ctx = context.Background()
	var imgPid = fileHash

	cld, _ := cloudinary.NewFromParams("dbqltwqac", "696441825692637", "YRXvf_DaiwJ7tGpLcAaXxyRgXxo")

	// Amazon Rekognition AI Moderation
	// Amazon Rekognition Celebrity Detection
	// Google Auto Tagging
	// Imagga Auto Tagging
	// Detection:   "openimages", "captioning", "aws_rek_tagging"
	resp1, err := cld.Upload.Upload(ctx, srcImgPath, uploader.UploadParams{
		PublicID: imgPid, Detection: "openimages", AutoTagging: 0.6})
	if err != nil {
		log.Fatalf("Failed to upload file, %v\n", err)
	}
	fmt.Printf("resp1 %+v", resp1)

	// Get details about the image with PublicID "my_image" and log the secure URL.
	resp2, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: imgPid})
	if err != nil {
		fmt.Println("error1", resp2)
	}

	// Instantiate an object for the asset with public ID "my_image"
	updImg, err := cld.Image(imgPid)
	if err != nil {
		fmt.Println("error2", err)
	}

	// reference: https://cloudinary.com/documentation/go_integration
	// case1) generative fill option - gen_fill, b_gen_fill, gen_replace, e_gen_replace, e_gen_remove
	// https://cloudinary.com/blog/generative-replace-object-replacement-with-ai
	// .imageTag("docs/bench-house.jpg");
	// updImg.Transformation = "ar_16:9,b_gen_fill:prompt_mug%20of%20coffee%20and%20cookies,c_pad,w_1920,h_1080"

	// Generate and print the delivery URL
	myURL, err := updImg.String()
	if err != nil {
		fmt.Println("error3", err)
	}
	return myURL

}
