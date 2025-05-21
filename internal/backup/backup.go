package backup

import (
	"fmt"
	"snapkeep/internal/config"
	"snapkeep/pkg/logger"
	"time"

	"path/filepath"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run(cfg *config.ApiConfig) error {
	envVariables, err := config.ParseEnv()
	if err != nil {
		logger.Error("Failed to parse environment variables:", err)
		return err
	}

	tmpDir := "tmp"
	timestamp := fmt.Sprint(time.Now().UnixMilli())
	zipFileName := tmpDir + "/backup_" + timestamp + ".zip"
	backupFolderPath := envVariables.BackupFolderPath
	backupFolderZipName := tmpDir + "/" + filepath.Base(backupFolderPath) + "_" + timestamp + ".zip"

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

	err = ZipDirectory(tmpDir, zipFileName)
	if err != nil {
		logger.Error("Failed to create zip file:", zipFileName, "Error:", err)
		return err
	}

	if err := ZipDirectory(backupFolderPath, backupFolderZipName); err != nil {
		logger.Error("Failed to create zip file:", zipFileName, "Error:", err)
		return err
	}

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
