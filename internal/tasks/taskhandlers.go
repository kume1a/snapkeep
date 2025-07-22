package tasks

import (
	"context"
	"fmt"
	"snapkeep/internal/backup"
	"snapkeep/internal/config"
	"snapkeep/internal/logger"

	"github.com/hibiken/asynq"
)

func CreateBackupAllAppsTaskHandler(cfg *config.ResourceConfig) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		envVars, err := config.ParseEnv()
		if err != nil {
			logger.Error("Failed to parse environment variables: ", err)
			return err
		}

		logger.Info("Starting sequential backup for all applications")

		for i, app := range envVars.Applications {
			logger.Info(fmt.Sprintf("Starting backup for app '%s' (%d/%d)", app.Name, i+1, len(envVars.Applications)))

			if err := backup.Run(
				ctx, cfg,
				app.Name,
				app.DbConnectionString,
				app.PublicFolderPath,
			); err != nil {
				logger.Error("Failed to run backup for app ", app.Name, ": ", err)
				return err
			}

			logger.Info(fmt.Sprintf("Completed backup for app '%s' (%d/%d)", app.Name, i+1, len(envVars.Applications)))
		}

		logger.Info("Completed sequential backup for all applications")
		return nil
	}
}
