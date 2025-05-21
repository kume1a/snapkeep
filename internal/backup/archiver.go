package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"snapkeep/pkg/logger"
)

func ZipDirectory(dirPath, zipFileName string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		logger.Error("Failed to create zip file:", zipFileName, "Error:", err)
		return err
	}

	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	defer zipWriter.Close()

	return filepath.Walk(
		dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logger.Error("Failed to access path:", path, "Error:", err)
				return err
			}

			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				logger.Error("Failed to get relative path for:", path, "Error:", err)
				return err
			}

			if relPath == "." {
				return nil // skip root
			}

			if info.IsDir() {
				// Add directory entry (with trailing slash)
				if _, err := zipWriter.Create(relPath + "/"); err != nil {
					logger.Error("Failed to create zip entry for directory:", relPath, "Error:", err)
					return err
				}

				return nil
			}

			fileToZip, err := os.Open(path)
			if err != nil {
				logger.Error("Failed to open file for zipping:", path, "Error:", err)
				return err
			}

			defer fileToZip.Close()

			zipEntryWriter, err := zipWriter.Create(relPath)
			if err != nil {
				logger.Error("Failed to create zip entry for:", relPath, "Error:", err)
				return err
			}

			if _, err := io.Copy(zipEntryWriter, fileToZip); err != nil {
				logger.Error("Failed to write file to zip:", relPath, "Error:", err)
				return err
			}

			return nil
		},
	)
}
