package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeBackupData = "data:backup"
)

type BackupDataPayload struct {
	BackupDBConnectionString string
	BackupFolderPath         string
}

func NewBackupDataTask(payload BackupDataPayload) (*asynq.Task, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeBackupData, bytes), nil
}
