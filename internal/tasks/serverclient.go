package tasks

import (
	"log"

	"snapkeep/internal/config"
	"snapkeep/pkg/logger"
	"time"

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
	mux.HandleFunc(TypeBackupData, CreateBackupDataTaskHandler(cfg))

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	return srv, nil
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
