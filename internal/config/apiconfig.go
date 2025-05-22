package config

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
)

type ResourceConfig struct {
	DB       *gorm.DB
	S3Client *s3.Client
}

type ApiConfig struct {
	*ResourceConfig
}
