package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"snapkeep/internal/logger"
)

func ZipDirectory(dirPath, zipFileName string) (string, error) {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		logger.Error("Failed to create zip file:", zipFileName, "Error:", err)
		return "", err
	}

	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	defer zipWriter.Close()

	if err := filepath.Walk(
		dirPath,
		zipWalkFunc(dirPath, zipWriter),
	); err != nil {
		return "", err
	}

	logger.Info("Zipping complete")

	absZipFileName, err := filepath.Abs(zipFileName)
	if err != nil {
		return "", err
	}

	return absZipFileName, nil
}

func zipWalkFunc(dirPath string, zipWriter *zip.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
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
			logger.Debug("Adding directory to zip:", relPath)
			header := &zip.FileHeader{
				Name:   relPath + "/",
				Method: zip.Store,
			}
			if _, err := zipWriter.CreateHeader(header); err != nil {
				logger.Error("Failed to create zip entry for directory:", relPath, "Error:", err)
				return err
			}
			return nil
		}

		logger.Debug("Adding file to zip:", relPath)
		fileToZip, err := os.Open(path)
		if err != nil {
			logger.Error("Failed to open file for zipping:", path, "Error:", err)
			return err
		}
		defer fileToZip.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			logger.Error("Failed to get file info header for:", relPath, "Error:", err)
			return err
		}

		header.Name = relPath
		header.Method = zip.Deflate
		header.UncompressedSize64 = uint64(info.Size())

		zipEntryWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			logger.Error("Failed to create zip entry for:", relPath, "Error:", err)
			return err
		}

		if _, err := io.Copy(zipEntryWriter, fileToZip); err != nil {
			logger.Error("Failed to write file to zip:", relPath, "Error:", err)
			return err
		}

		return nil
	}
}
