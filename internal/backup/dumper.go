package backup

import (
	"encoding/json"
	"os"
	"snapkeep/pkg/logger"

	"gorm.io/gorm"
)

func DumpDatabaseTablesToJson(db *gorm.DB) ([]string, error) {
	err := os.MkdirAll("tmp", 0755)
	if err != nil {
		logger.Error("Failed to create tmp directory:", err)
		return nil, err
	}

	var tableNames []string
	err = db.
		Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").
		Scan(&tableNames).
		Error

	if err != nil {
		logger.Error("Failed to retrieve table names:", err)
		return nil, err
	}

	var filePaths []string
	for _, tableName := range tableNames {
		filePath := "tmp/" + tableName + ".json"
		err := dumpDatabaseTableToJSON(db, tableName, filePath)
		if err != nil {
			logger.Error("Failed to dump table:", tableName, "Error:", err)
			return nil, err
		}
		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil
}

// TODO indent json objects properly
func dumpDatabaseTableToJSON(db *gorm.DB, tableName, filePath string) error {
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

		result := db.Table(tableName).Limit(batchSize).Offset(offset).Find(&rows)

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
	return nil
}
