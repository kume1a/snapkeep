package config

import (
	"errors"
	"fmt"
	"os"
	"snapkeep/pkg/logger"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	if env == "development" {
		envPath := ".env." + env

		logger.Debug("Loading env file: " + envPath)

		godotenv.Load(envPath)
	}
}

type EnvVariables struct {
	IsDevelopment            bool
	IsProduction             bool
	Port                     string
	DbConnectionString       string
	BackupDbConnectionString string
	BackupFolderPath         string
	AWSS3LimitBytes          int64
	AWSS3BackupBucketName    string
	RedisAddress             string
	RedisPassword            string
}

func ParseEnv() (*EnvVariables, error) {
	environment, err := getEnv("ENVIRONMENT")
	if err != nil {
		return nil, err
	}

	port, err := getEnv("PORT")
	if err != nil {
		return nil, err
	}

	dbConnectionString, err := getEnv("DB_CONNECTION_STRING")
	if err != nil {
		return nil, err
	}

	backupDbConnectionString, err := getEnv("BACKUP_DB_CONNECTION_STRING")
	if err != nil {
		return nil, err
	}

	backupFolderPath, err := getEnv("BACKUP_FOLDER_PATH")
	if err != nil {
		return nil, err
	}

	awsS3LimitBytes, err := getEnvInt("AWS_S3_LIMIT_BYTES")
	if err != nil {
		return nil, err
	}

	awsS3BackupBucketName, err := getEnv("AWS_S3_BACKUP_BUCKET_NAME")
	if err != nil {
		return nil, err
	}

	redisAddress, err := getEnv("REDIS_ADDRESS")
	if err != nil {
		return nil, err
	}

	redisPassword, err := getEnv("REDIS_PASSWORD")
	if err != nil {
		return nil, err
	}

	return &EnvVariables{
		IsDevelopment:            environment == "development",
		IsProduction:             environment == "production",
		Port:                     port,
		DbConnectionString:       dbConnectionString,
		BackupDbConnectionString: backupDbConnectionString,
		BackupFolderPath:         backupFolderPath,
		AWSS3LimitBytes:          awsS3LimitBytes,
		AWSS3BackupBucketName:    awsS3BackupBucketName,
		RedisAddress:             redisAddress,
		RedisPassword:            redisPassword,
	}, nil
}

func getEnvInt(key string) (int64, error) {
	val, err := getEnv(key)
	if err != nil {
		return 0, err
	}

	valInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return valInt, nil
}

func getEnv(key string) (string, error) {
	envVar := os.Getenv(key)
	if envVar == "" {
		msg := fmt.Sprintf("%v is not found in the env", key)

		logger.Fatal(msg)
		return "", errors.New(msg)
	}

	return envVar, nil
}
