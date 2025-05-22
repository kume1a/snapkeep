package main

import (
	"context"
	"snapkeep/internal/backup"
	"snapkeep/internal/config"
	"snapkeep/pkg/logger"
)

// https://github.com/hibiken/asynq

func main() {
	ctx := context.Background()

	config.LoadEnv()

	db, err := config.InitializeDB()
	if err != nil {
		logger.Fatal("Failed to initialize database: ", err)
		return
	}

	s3Client, err := config.InitializeS3Client(ctx)
	if err != nil {
		logger.Fatal("Failed to initialize S3 client: ", err)
		return
	}

	cfg := &config.ApiConfig{
		ResourceConfig: &config.ResourceConfig{
			DB:       db,
			S3Client: s3Client,
		},
	}

	if err := backup.Run(ctx, cfg); err != nil {
		logger.Fatal("Failed to run backup: ", err)
		return
	}
}
