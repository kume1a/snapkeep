package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeBackupAllApps = "data:backup_all_apps"
)

type BackupAllAppsPayload struct {
}

func NewBackupAllAppsTask() (*asynq.Task, error) {
	payload := BackupAllAppsPayload{}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeBackupAllApps, bytes), nil
}
