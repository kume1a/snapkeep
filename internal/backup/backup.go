package backup

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"snapkeep/internal/config"
	"snapkeep/internal/db"
	"snapkeep/internal/logger"
	"snapkeep/internal/shared"
	"time"

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

	formattedNow := time.Now().Format("02_01_2006_15:04")
	zipedBackupDatabaseDestination := directoryPath + "/database_" + formattedNow + ".zip"
	zippedBackupFolderDestination := directoryPath + "/folder_" + formattedNow + ".zip"

	database, err := gorm.Open(postgres.Open(backupDBConnectionString), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to open backup database:", err)
		return err
	}

	if _, err := DumpDatabaseTablesToJson(database, directoryPath); err != nil {
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

	latestBackup, err := db.GetLatestActiveBackup(cfg.DB)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to get latest active backup:", err)
		return err
	}

	if err == nil && latestBackup != nil {
		if err := DeleteS3Folder(
			ctx, cfg.S3Client,
			envVariables.AWSS3BackupBucketName,
			latestBackup.BackupName,
		); err != nil {
			logger.Error("Failed to delete S3 folder:", err)
			return err
		}

		if err := db.MarkBackupAsDeleted(cfg.DB, latestBackup.ID); err != nil {
			logger.Error("Failed to mark backup as deleted in DB:", err)
			return err
		}
	}

	zippedBackupDbFileSize, err := shared.GetFileSize(zippedBackupDatabasePath)
	if err != nil {
		logger.Error("Failed to get zipped backup database file size:", err)
		return err
	}

	zippedBackupFolderSize, err := shared.GetFileSize(zippedBackupFolderPath)
	if err != nil {
		logger.Error("Failed to get zipped backup folder file size:", err)
		return err
	}

	s3Size, err := GetS3FolderSize(ctx, cfg.S3Client, envVariables.AWSS3BackupBucketName, "")
	if err != nil {
		logger.Error("Failed to get S3 folder size:", err)
		return err
	}

	totalSizeAfterUpload := s3Size.InBytes +
		zippedBackupDbFileSize.InBytes +
		zippedBackupFolderSize.InBytes

	if totalSizeAfterUpload > envVariables.AWSS3LimitBytes {
		logger.Error("S3 usage with new backup would exceed limit, aborting backup upload.")
		return fmt.Errorf("S3 usage with new backup would exceed limit, aborting backup upload")
	}

	zippedBackupDatabaseFile, err := os.Open(zippedBackupDatabasePath)
	if err != nil {
		logger.Error("Failed to open zipped backup database file for upload:", err)
		return err
	}

	uploadedBackupDatabaseZipURL, err := UploadFileToS3(UploadFileToS3Input{
		Context:     ctx,
		S3Client:    cfg.S3Client,
		Bucket:      envVariables.AWSS3BackupBucketName,
		Prefix:      backupName,
		Key:         filepath.Base(zippedBackupDatabasePath),
		Body:        zippedBackupDatabaseFile,
		ContentType: "application/zip",
	})
	zippedBackupDatabaseFile.Close()
	if err != nil {
		logger.Error("Failed to upload zipped backup database to S3:", err)
		return err
	}

	zippedBackupFolderFile, err := os.Open(zippedBackupFolderPath)
	if err != nil {
		logger.Error("Failed to open zipped backup folder file for upload:", err)
		return err
	}

	uploadedBackupFolderZipURL, err := UploadFileToS3(UploadFileToS3Input{
		Context:     ctx,
		S3Client:    cfg.S3Client,
		Bucket:      envVariables.AWSS3BackupBucketName,
		Prefix:      backupName,
		Key:         filepath.Base(zippedBackupFolderPath),
		Body:        zippedBackupFolderFile,
		ContentType: "application/zip",
	})
	zippedBackupFolderFile.Close()
	if err != nil {
		logger.Error("Failed to upload zipped backup folder to S3:", err)
		return err
	}

	logger.Debug("Uploaded zipped backup database URL:", uploadedBackupDatabaseZipURL)
	logger.Debug("Uploaded zipped backup folder URL:", uploadedBackupFolderZipURL)

	backupEntity := &db.Backup{
		BackupName:            backupName,
		BackupDBSizeBytes:     uint64(zippedBackupDbFileSize.InBytes),
		BackupDBUrl:           uploadedBackupDatabaseZipURL,
		BackupFolderSizeBytes: uint64(zippedBackupFolderSize.InBytes),
		BackupFolderUrl:       uploadedBackupFolderZipURL,
		Status:                db.BackupStatusActive,
	}
	if err := cfg.DB.Create(backupEntity).Error; err != nil {
		logger.Error("Failed to save backup entity to DB:", err)
	}

	logger.Info("Backup completed successfully.")

	if err := os.RemoveAll(directoryPath); err != nil {
		logger.Error("Failed to remove temporary directory:", directoryPath, "Error:", err)
		return err
	}

	return nil
}
