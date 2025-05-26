package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"snapkeep/internal/backup"
	"snapkeep/internal/config"
	"snapkeep/internal/logger"

	"github.com/hibiken/asynq"
)

func CreateBackupDataTaskHandler(cfg *config.ResourceConfig) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p BackupDataPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}

		if err := backup.Run(
			ctx, cfg,
			p.BackupDBConnectionString,
			p.BackupFolderPath,
			p.BackupName,
		); err != nil {
			logger.Error("Failed to run backup: ", err)
			return err
		}

		return nil
	}
}
