package db

import (
	"snapkeep/internal/config"
	"snapkeep/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	envVars, err := config.ParseEnv()
	if err != nil {
		logger.Fatal("Coultn't parse env vars, returning nil, err:", err)
		return nil, err
	}

	database, err := gorm.Open(postgres.Open(envVars.DbConnectionString), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database, err:", err)
		return nil, err
	}

	database.AutoMigrate(&Backup{})

	return database, nil
}
