package config

import (
	"context"
	"snapkeep/internal/logger"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitializeS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)

	if err != nil {
		logger.Fatal("Failed to load AWS config: ", err)
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}
