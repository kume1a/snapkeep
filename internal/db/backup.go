package db

import (
	"gorm.io/gorm"
)

func GetLatestActiveBackupByAppName(dbConn *gorm.DB, appName string) (*Backup, error) {
	var latest Backup

	err := dbConn.Model(&Backup{}).
		Where("status = ? AND backup_name = ?", BackupStatusActive, appName).
		Order("created_at desc").
		First(&latest).Error

	if err != nil {
		return nil, err
	}

	return &latest, nil
}

func MarkBackupAsDeleted(dbConn *gorm.DB, backupID uint) error {
	return dbConn.Model(&Backup{}).
		Where("id = ?", backupID).
		Update("status", BackupStatusDeleted).Error
}
