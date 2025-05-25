package config

import (
	"snapkeep/internal/db"
	"snapkeep/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	envVars, err := ParseEnv()
	if err != nil {
		logger.Fatal("Coultn't parse env vars, returning nil, err:", err)
		return nil, err
	}

	database, err := gorm.Open(postgres.Open(envVars.DbConnectionString), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database, err:", err)
		return nil, err
	}

	database.AutoMigrate(&db.Backup{})

	return database, nil
}
