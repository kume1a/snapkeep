package backup

import (
	"fmt"
	"os"
	"snapkeep/internal/config"
	"snapkeep/internal/uploader"
	"time"
)

func Run(cfg config.Config) error {
	timestamp := time.Now().Format("20060102_150405")
	os.MkdirAll("tmp", 0755)
	dumpPath := fmt.Sprintf("tmp/db_%s.sql", timestamp)
	zipPath := fmt.Sprintf("tmp/backup_%s.zip", timestamp)

	if err := DumpDatabase(cfg.DBURL, dumpPath); err != nil {
		return err
	}
	if err := ZipFolder(cfg.PublicDir, zipPath); err != nil {
		return err
	}
	if err := ZipFolder(cfg.BackendDir, zipPath); err != nil {
		return err
	}
	return uploader.Upload(cfg.UploadTarget, zipPath)
}
