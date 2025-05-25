package backup

import (
	"context"
	"fmt"
	"os"
	"snapkeep/internal/config"
	"snapkeep/internal/uploader"
	"snapkeep/pkg/logger"
	"time"

	"path/filepath"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run(
	ctx context.Context,
	cfg *config.ResourceConfig,
	backupDBConnectionString string,
	backupFolderPath string,
	backupName string,
) error {
	envVariables, err := config.ParseEnv()
	if err != nil {
		logger.Error("Failed to parse environment variables:", err)
		return err
	}

	directoryPath := "./tmp/" + backupName

	timestamp := fmt.Sprint(time.Now().UnixMilli())
	zipedBackupDatabaseDestination := directoryPath + "/backup_" + timestamp + ".zip"
	zippedBackupFolderDestination := directoryPath + "/" + filepath.Base(backupFolderPath) + "_" + timestamp + ".zip"

	db, err := gorm.Open(postgres.Open(backupDBConnectionString), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to open backup database:", err)
		return err
	}

	if _, err := DumpDatabaseTablesToJson(db, directoryPath); err != nil {
		logger.Error("Failed to dump database tables to JSON:", err)
		return err
	}

	logger.Debug("All tables exported successfully.")

	zippedBackupDatabasePath, err := ZipDirectory(directoryPath, zipedBackupDatabaseDestination)
	if err != nil {
		logger.Error("Failed to create zip file:", zipedBackupDatabaseDestination, "Error:", err)
		return err
	}

	zippedBackupFolderPath, err := ZipDirectory(backupFolderPath, zippedBackupFolderDestination)
	if err != nil {
		logger.Error("Failed to create zip file:", zipedBackupDatabaseDestination, "Error:", err)
		return err
	}

	logger.Debug("Zipped backup database path:", zippedBackupDatabasePath)
	logger.Debug("Zipped backup folder path:", zippedBackupFolderPath)

	zippedBackupDatabaseFile, err := os.Open(zippedBackupDatabasePath)
	if err != nil {
		logger.Error("Failed to open zipped backup database file for upload:", err)
		return err
	}
	defer zippedBackupDatabaseFile.Close()

	uploadedBackupDatabaseZipURL, err := uploader.UploadFileToS3(uploader.UploadFileToS3Input{
		Context:     ctx,
		S3Client:    cfg.S3Client,
		Bucket:      envVariables.AWSS3BackupBucketName,
		Prefix:      backupName,
		Key:         filepath.Base(zippedBackupDatabasePath),
		Body:        zippedBackupDatabaseFile,
		ContentType: "application/zip",
	})
	if err != nil {
		logger.Error("Failed to upload zipped backup database to S3:", err)
		return err
	}

	zippedBackupFolderFile, err := os.Open(zippedBackupFolderPath)
	if err != nil {
		logger.Error("Failed to open zipped backup folder file for upload:", err)
		return err
	}
	defer zippedBackupFolderFile.Close()

	uploadedBackupFolderZipURL, err := uploader.UploadFileToS3(uploader.UploadFileToS3Input{
		Context:     ctx,
		S3Client:    cfg.S3Client,
		Bucket:      envVariables.AWSS3BackupBucketName,
		Prefix:      backupName,
		Key:         filepath.Base(zippedBackupFolderPath),
		Body:        zippedBackupFolderFile,
		ContentType: "application/zip",
	})
	if err != nil {
		logger.Error("Failed to upload zipped backup folder to S3:", err)
		return err
	}

	logger.Debug("Uploaded zipped backup database URL:", uploadedBackupDatabaseZipURL)
	logger.Debug("Uploaded zipped backup folder URL:", uploadedBackupFolderZipURL)

	logger.Info("Backup completed successfully.")

	if err := os.RemoveAll(directoryPath); err != nil {
		logger.Error("Failed to remove temporary directory:", directoryPath, "Error:", err)
		return err
	}

	return nil
}
