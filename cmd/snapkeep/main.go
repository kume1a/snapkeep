package main

import (
	"snapkeep/internal/backup"
	"snapkeep/internal/config"
	"snapkeep/pkg/logger"
)

func main() {
	config.LoadEnv()

	db, err := config.InitializeDB()
	if err != nil {
		logger.Fatal("Failed to initialize database: ", err)
		return
	}

	cfg := &config.ApiConfig{
		ResourceConfig: &config.ResourceConfig{
			DB: db,
		},
	}

	if err := backup.Run(cfg); err != nil {
		logger.Fatal("Failed to run backup: ", err)
		return
	}

	logger.Info("Backup completed successfully")
}
