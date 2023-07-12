package utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
	"os"
	"sync"
)

const endPoint = "https://kr.object.ncloudstorage.com"
const regionName = "kr-standard"
const accessKey = "LR9hW8gPyStOvacROGl3"
const secretKey = "JZE1FXZObcu2tazxgEmQzx043XK8ZmSfv8JR1vlg"
const bucketName = "alpha.gen8.id"

type CustomReader struct {
	fp      *os.File
	size    int64
	read    int64
	signMap map[int64]struct{}
	mux     sync.Mutex
}

func (r *CustomReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

func (r *CustomReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.fp.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	r.mux.Lock()
	// Ignore the first signature call
	if _, ok := r.signMap[off]; ok {
		// Got the length have read( or means has uploaded), and you can construct your message
		r.read += int64(n)
		fmt.Printf("\rtotal read:%d    progress:%d%%", r.read, int(float32(r.read*100)/float32(r.size)))
	} else {
		r.signMap[off] = struct{}{}
	}
	r.mux.Unlock()
	return n, err
}

func (r *CustomReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}

// AWS_REGION=<region> go run -tags example putObjWithProcess.go <bucket> <key for object> <local file name>
// const endPoint = "https://kr.object.ncloudstorage.com"
// const regionName = "kr-standard"
// const accessKey = "LR9hW8gPyStOvacROGl3"
// const secretKey = "JZE1FXZObcu2tazxgEmQzx043XK8ZmSfv8JR1vlg"
func objectUpload(filename string) string {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           endPoint,
			SigningRegion: region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey, secretKey, "")))
	if err != nil {
		// handle error
	}

	client := s3.NewFromConfig(cfg)

	uploadFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file %v, %v", filename, err)
	}

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = true
	})

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		ACL:         types.ObjectCannedACLPublicRead, //  aws.String("public-read"),
		Key:         aws.String("pblc/" + filename),
		Body:        uploadFile,
		ContentType: aws.String("image/webp"),
	})
	if err != nil {
		log.Fatalf("failed to put file %v, %v", filename, err)
		return ""
	}

	fmt.Println()
	log.Println(result.Location)
	return result.Location

}
