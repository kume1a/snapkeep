package config

import (
	"gorm.io/gorm"
)

type Backup struct {
	gorm.Model
	SizeKB uint64
	URL    string
}
