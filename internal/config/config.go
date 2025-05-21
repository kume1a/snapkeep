package config

import (
	"flag"
)

type Config struct {
	DBURL        string
	PublicDir    string
	BackendDir   string
	UploadTarget string
}

func Load() Config {
	var cfg Config
	flag.StringVar(&cfg.DBURL, "db", "postgres://user:pass@localhost/db", "Database URL")
	flag.StringVar(&cfg.PublicDir, "public", "./public", "Public directory to backup")
	flag.StringVar(&cfg.BackendDir, "backend", "../app-backend", "Backend directory to backup")
	flag.StringVar(&cfg.UploadTarget, "upload", "s3", "Upload target (s3/gdrive)")
	flag.Parse()
	return cfg
}
