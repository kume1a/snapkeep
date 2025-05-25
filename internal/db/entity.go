package db

import (
	"gorm.io/gorm"
)

type Backup struct {
	gorm.Model
	BackupDBSizeBytes     uint64
	BackupDBUrl           string
	BackupFolderSizeBytes uint64
	BackupFolderUrl       string
}
