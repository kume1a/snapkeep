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

	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		logger.Error("Failed to read directory:", dirPath, "Error:", err)
		return err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join(dirPath, entry.Name())
		fileToZip, err := os.Open(filePath)
		if err != nil {
			logger.Error("Failed to open file for zipping:", filePath, "Error:", err)
			return err
		}
		defer fileToZip.Close()

		zipEntryWriter, err := zipWriter.Create(entry.Name())
		if err != nil {
			logger.Error("Failed to create zip entry for:", entry.Name(), "Error:", err)
			return err
		}

		_, err = io.Copy(zipEntryWriter, fileToZip)
		if err != nil {
			logger.Error("Failed to write file to zip:", entry.Name(), "Error:", err)
			return err
		}
	}
	return nil
}
