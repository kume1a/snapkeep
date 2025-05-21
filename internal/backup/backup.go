package backup

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"snapkeep/internal/config"
	"snapkeep/pkg/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run(cfg *config.ApiConfig) error {
	timestamp := time.Now().UnixMilli()
	os.MkdirAll("tmp", 0755)

	logger.Debug("timestamp:", timestamp)

	tempDbConnectionURL := "..."
	backupDB, err := gorm.Open(postgres.Open(tempDbConnectionURL), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to backup database:", err)
		return err
	}

	var tableNames []string
	err = backupDB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableNames).Error
	if err != nil {
		logger.Error("Failed to retrieve table names:", err)
		return nil
	}

	logger.Debug("Table names:", tableNames)

	for _, tableName := range tableNames {
		filePath := "tmp/" + tableName + ".json"
		logger.Debug("Exporting table:", tableName, "to file:", filePath)

		file, err := os.Create(filePath)
		if err != nil {
			logger.Error("Failed to create file:", filePath, "Error:", err)
			return err
		}
		defer file.Close()

		file.WriteString("[")
		batchSize := 1000
		offset := 0
		rowCount := 0
		for {
			var rows []map[string]any
			result := backupDB.Table(tableName).Limit(batchSize).Offset(offset).Find(&rows)
			if result.Error != nil {
				logger.Error("Failed to query table:", tableName, "Error:", result.Error)
				return result.Error
			}

			if len(rows) == 0 {
				break
			}

			for _, row := range rows {
				if rowCount == 0 {
					file.WriteString("\n")
				}
				jsonBytes, err := json.MarshalIndent(row, "", "  ")
				if err != nil {
					logger.Error("Failed to marshal row to JSON in table:", tableName, "Error:", err)
					return err
				}

				if rowCount > 0 {
					file.WriteString(",\n")
				}

				file.Write(jsonBytes)
				rowCount++
			}
			offset += batchSize
		}
		if rowCount > 0 {
			file.WriteString("\n]")
		} else {
			file.WriteString("]")
		}
		logger.Debug("Exported table:", tableName, "to file:", filePath)
	}

	logger.Debug("All tables exported successfully.")

	zipFileName := "backup_" + fmt.Sprint(timestamp) + ".zip"
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		logger.Error("Failed to create zip file:", zipFileName, "Error:", err)
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	dirEntries, err := os.ReadDir("tmp")
	if err != nil {
		logger.Error("Failed to read tmp directory:", err)
		return err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		filePath := "tmp/" + entry.Name()
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

	logger.Debug("Created zip file:", zipFileName)

	return nil
}
