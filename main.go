package main

import (
	"context"
	"os"
	"os/signal"
	"snapkeep/internal/config"
	"snapkeep/internal/tasks"
	"snapkeep/pkg/logger"
	"syscall"
)

func main() {
	ctx := context.Background()

	config.LoadEnv()

	envVars, err := config.ParseEnv()
	if err != nil {
		logger.Fatal("Failed to parse environment variables: ", err)
		return
	}

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

	taskClient, err := tasks.InitializeTaskClient()
	if err != nil {
		logger.Fatal("Failed to initialize task client: ", err)
		return
	}
	defer taskClient.Close()

	cfg := &config.ApiConfig{
		ResourceConfig: &config.ResourceConfig{
			DB:         db,
			S3Client:   s3Client,
			TaskClient: taskClient,
		},
	}

	taskServer, err := tasks.InitializeTaskServer(cfg.ResourceConfig)
	if err != nil {
		logger.Fatal("Failed to initialize task server: ", err)
		return
	}

	task, err := tasks.NewBackupDataTask(
		tasks.BackupDataPayload{
			BackupName:               envVars.BackupName,
			BackupDBConnectionString: envVars.BackupDbConnectionString,
			BackupFolderPath:         envVars.BackupFolderPath,
		},
	)
	if err != nil {
		logger.Fatal("Failed to create backup data task: ", err)
		return
	}

	info, err := taskClient.Enqueue(task)
	if err != nil {
		logger.Fatal("Failed to enqueue backup data task: ", err)
		return
	}

	logger.Info("Enqueued task: id=" + info.ID + ", type=" + info.Type + ", queue=" + info.Queue)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down task server...")
	taskServer.Shutdown()
}
