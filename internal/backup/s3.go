package backup

import (
	"context"
	"io"
	"snapkeep/internal/logger"
	"snapkeep/internal/shared"

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

type progressReader struct {
	reader     io.Reader
	total      int64
	read       int64
	logEvery   int64
	lastLogged int64
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.reader.Read(b)
	if n > 0 {
		p.read += int64(n)
		if p.read-p.lastLogged >= p.logEvery || err == io.EOF {
			logger.Info("S3 upload progress: ", p.read/(1024*1024), "/", p.total/(1024*1024), " MB")
			p.lastLogged = p.read
		}
	}
	return n, err
}

func UploadFileToS3(input UploadFileToS3Input) (string, error) {
	fullKey := ""
	if input.Prefix != "" {
		fullKey = input.Prefix + "/" + input.Key
	} else {
		fullKey = input.Key
	}

	var body io.Reader = input.Body
	var contentLength *int64 = nil
	if seeker, ok := input.Body.(io.Seeker); ok {
		if size, err := seeker.Seek(0, io.SeekEnd); err == nil {
			_, _ = seeker.Seek(0, io.SeekStart)
			body = &progressReader{
				reader:   input.Body,
				total:    size,
				logEvery: 100 * 1024 * 1024, // log every 100MB
			}
			contentLength = &size
		}
	}

	s3Input := &s3.PutObjectInput{
		Bucket:      &input.Bucket,
		Key:         &fullKey,
		Body:        body,
		ContentType: &input.ContentType,
		ACL:         s3types.ObjectCannedACLPublicRead,
	}
	if contentLength != nil {
		s3Input.ContentLength = contentLength
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
