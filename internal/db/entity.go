package db

import (
	"gorm.io/gorm"
)

type BackupStatus string

const (
	BackupStatusActive  BackupStatus = "active"
	BackupStatusDeleted BackupStatus = "deleted"
)

type Backup struct {
	gorm.Model
	BackupName            string
	BackupDBSizeBytes     uint64
	BackupDBUrl           string
	BackupFolderSizeBytes uint64
	BackupFolderUrl       string
	Status                BackupStatus
}
