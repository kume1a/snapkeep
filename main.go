package main

import (
	"context"
	"os"
	"os/signal"
	"snapkeep/internal/config"
	"snapkeep/internal/db"
	"snapkeep/internal/tasks"
	"snapkeep/internal/webserver"
	"snapkeep/pkg/logger"
	"syscall"
)

func main() {
	ctx := context.Background()

	config.LoadEnv()

	db, err := db.InitializeDB()
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

	taskScheduler, err := tasks.InitializeTaskScheduler(cfg.ResourceConfig)
	if err != nil {
		logger.Fatal("Failed to initialize task scheduler: ", err)
		return
	}

	if err := tasks.EnqueueAndScheduleTasks(taskClient, taskScheduler); err != nil {
		logger.Fatal("Failed to enqueue and schedule tasks: ", err)
		return
	}

	if err := webserver.ConfigureWebServer(); err != nil {
		logger.Fatal("Failed to configure web server: ", err)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down task server and scheduler...")
	taskServer.Shutdown()
	taskScheduler.Shutdown()
}
