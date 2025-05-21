package uploader

import "fmt"

func UploadS3(filePath string) error {
	fmt.Println("Uploading to S3:", filePath)
	return nil
}
