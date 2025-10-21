package tasks

import (
	"log"
	"time"

	"snapkeep/internal/config"
	"snapkeep/internal/logger"

	"github.com/hibiken/asynq"
)

func InitializeTaskClient() (*asynq.Client, error) {
	redisClientOpt, err := getRedisClientOpt()
	if err != nil {
		logger.Fatal("Failed to get Redis client options: ", err)
		return nil, err
	}

	client := asynq.NewClient(redisClientOpt)

	return client, nil
}

func InitializeTaskServer(cfg *config.ResourceConfig) (*asynq.Server, error) {
	redisClientOpt, err := getRedisClientOpt()
	if err != nil {
		logger.Fatal("Failed to get Redis client options: ", err)
		return nil, err
	}

	srv := asynq.NewServer(
		redisClientOpt,
		asynq.Config{
			Concurrency: 1,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			TaskCheckInterval: time.Second * 3,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeBackupAllApps, CreateBackupAllAppsTaskHandler(cfg))

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	return srv, nil
}

func InitializeTaskScheduler(cfg *config.ResourceConfig) (*asynq.Scheduler, error) {
	redisClientOpt, err := getRedisClientOpt()
	if err != nil {
		logger.Fatal("Failed to get Redis client options: ", err)
		return nil, err
	}

	location, err := time.LoadLocation("Asia/Tbilisi")
	if err != nil {
		logger.Fatal("Failed to load Asia/Tbilisi location: ", err)
		return nil, err
	}

	scheduler := asynq.NewScheduler(redisClientOpt, &asynq.SchedulerOpts{
		Location: location,
	})

	go func() {
		if err := scheduler.Run(); err != nil {
			log.Fatalf("could not run scheduler: %v", err)
		}
	}()

	return scheduler, nil
}

func ScheduleBackgroundTasks(
	client *asynq.Client,
	scheduler *asynq.Scheduler,
) error {
	envVars, err := config.ParseEnv()
	if err != nil {
		logger.Error("Failed to parse environment variables: ", err)
		return err
	}

	task, err := NewBackupAllAppsTask()
	if err != nil {
		logger.Fatal("Failed to create backup all apps task: ", err)
		return err
	}

	if envVars.RunBackupOnStart {
		log.Printf("Running backup on start for %d applications", len(envVars.Applications))
		_, err := client.Enqueue(task)
		if err != nil {
			logger.Error("Failed to enqueue backup task on start: ", err)
		} else {
			log.Printf("Successfully enqueued backup task on start")
		}
	}

	entryID, err := scheduler.Register("0 5 */3 * *", task)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Registered sequential backup task for %d applications with entry ID: %s", len(envVars.Applications), entryID)

	for i, app := range envVars.Applications {
		log.Printf("  App %d: %s (Database: %s, Folder: %s)",
			i+1,
			app.Name,
			"configured",
			func() string {
				if app.PublicFolderPath != "" {
					return "configured"
				}
				return "none"
			}())
	}

	return nil
}

func getRedisClientOpt() (asynq.RedisClientOpt, error) {
	envVars, err := config.ParseEnv()
	if err != nil {
		logger.Error("Failed to parse environment variables: ", err)
		return asynq.RedisClientOpt{}, err
	}

	return asynq.RedisClientOpt{
		Addr:     envVars.RedisAddress,
		Password: envVars.RedisPassword,
		DB:       0,
	}, nil
}
