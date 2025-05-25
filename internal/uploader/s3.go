package uploader

import (
	"context"
	"io"
	"snapkeep/internal/shared"
	"snapkeep/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type UploadFileToS3Input struct {
	Context     context.Context
	S3Client    *s3.Client
	Bucket      string
	Prefix      string
	Key         string
	Body        io.Reader
	ContentType string
}

func UploadFileToS3(input UploadFileToS3Input) (string, error) {
	key := ""
	if input.Prefix != "" {
		key = input.Prefix + "/" + input.Key
	} else {
		key = input.Key
	}

	s3Input := &s3.PutObjectInput{
		Bucket:      &input.Bucket,
		Key:         &key,
		Body:        input.Body,
		ContentType: &input.ContentType,
		ACL:         types.ObjectCannedACLPublicRead,
	}

	_, err := input.S3Client.PutObject(input.Context, s3Input)
	if err != nil {
		logger.Error("Failed to upload file to S3: ", err)
		return "", err
	}

	publicURL := "https://" + input.Bucket + ".s3.amazonaws.com/" + key
	return publicURL, nil
}

func GetS3FolderSize(
	ctx context.Context,
	client *s3.Client,
	bucket string,
	prefix string,
) (*shared.FileSizeInUnits, error) {
	var total int64
	params := &s3.ListObjectsV2Input{
		Bucket: &bucket,
	}
	if prefix != "" {
		params.Prefix = &prefix
	}

	p := s3.NewListObjectsV2Paginator(client, params)
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			logger.Error("Failed to get next s3 page: ", err)
			return nil, err
		}

		for _, obj := range page.Contents {
			total += *obj.Size
		}
	}

	size := shared.ConvertBytes(total)

	return &size, nil
}
