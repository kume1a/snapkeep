package uploader

func Upload(target, filePath string) error {
	switch target {
	case "s3":
		return UploadS3(filePath)
	case "gdrive":
		return UploadGDrive(filePath)
	default:
		return nil
	}
}
