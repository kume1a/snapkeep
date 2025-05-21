package config

import (
	"gorm.io/gorm"
)

type ResourceConfig struct {
	DB *gorm.DB
}

type ApiConfig struct {
	*ResourceConfig
}
