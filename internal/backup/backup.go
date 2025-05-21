package backup

import (
	"fmt"
	"snapkeep/internal/config"
	"snapkeep/pkg/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run(cfg *config.ApiConfig) error {
	dumpDir := "tmp"
	zipFileName := dumpDir + "/backup_" + fmt.Sprint(time.Now().UnixMilli()) + ".zip"

	db, err := openBackupDB()
	if err != nil {
		logger.Error("Failed to open backup database:", err)
		return err
	}

	if err := DumpDatabaseTablesToJson(db); err != nil {
		logger.Error("Failed to dump database tables to JSON:", err)
		return err
	}

	logger.Debug("All tables exported successfully.")

	err = ZipDirectory(dumpDir, zipFileName)
	if err != nil {
		logger.Error("Failed to create zip file:", zipFileName, "Error:", err)
		return err
	}

	logger.Debug("Created zip file:", zipFileName)

	return nil
}

func openBackupDB() (*gorm.DB, error) {
	envVariables, err := config.ParseEnv()
	if err != nil {
		logger.Error("Failed to parse environment variables:", err)
		return nil, err
	}

	return gorm.Open(postgres.Open(envVariables.BackupDbConnectionString), &gorm.Config{})
}
