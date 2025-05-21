package uploader

import "fmt"

func UploadGDrive(filePath string) error {
	fmt.Println("Uploading to Google Drive:", filePath)
	return nil
}
