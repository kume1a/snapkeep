package backup

import (
	"context"
	"io"
	"snapkeep/internal/shared"
	"snapkeep/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
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
	fullKey := ""
	if input.Prefix != "" {
		fullKey = input.Prefix + "/" + input.Key
	} else {
		fullKey = input.Key
	}

	s3Input := &s3.PutObjectInput{
		Bucket:      &input.Bucket,
		Key:         &fullKey,
		Body:        input.Body,
		ContentType: &input.ContentType,
		ACL:         s3types.ObjectCannedACLPublicRead,
	}

	_, err := input.S3Client.PutObject(input.Context, s3Input)
	if err != nil {
		logger.Error("Failed to upload file to S3: ", err)
		return "", err
	}

	publicURL := "https://" + input.Bucket + ".s3.amazonaws.com/" + fullKey
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

func DeleteS3File(
	ctx context.Context,
	client *s3.Client, bucket,
	prefix, key string,
) error {
	fullKey := key
	if prefix != "" {
		fullKey = prefix + "/" + key
	}

	_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &fullKey,
	})
	return err
}

func DeleteS3Folder(ctx context.Context, client *s3.Client, bucket, prefix string) error {
	input := &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	}

	paginator := s3.NewListObjectsV2Paginator(client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		var objectsToDelete []s3types.ObjectIdentifier
		for _, obj := range page.Contents {
			objectsToDelete = append(objectsToDelete, s3types.ObjectIdentifier{Key: obj.Key})
		}

		if len(objectsToDelete) > 0 {
			quiet := true
			_, err := client.DeleteObjects(
				ctx,
				&s3.DeleteObjectsInput{
					Bucket: &bucket,
					Delete: &s3types.Delete{
						Objects: objectsToDelete,
						Quiet:   &quiet,
					},
				},
			)
			if err != nil {
				continue
			}
		}
	}

	return nil
}
