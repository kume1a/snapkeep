package config

import (
	"errors"
	"fmt"
	"os"
	"snapkeep/internal/logger"
	"strconv"
	"strings"

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

type BackupAppConfig struct {
	Name               string
	DbConnectionString string
	PublicFolderPath   string // Optional - empty string means no folder backup
}

type EnvVariables struct {
	IsDevelopment         bool
	IsProduction          bool
	Port                  string
	DbConnectionString    string
	Applications          []BackupAppConfig
	AWSS3LimitBytes       int64
	AWSS3BackupBucketName string
	RedisAddress          string
	RedisPassword         string
	AdminPassword         string
	AccessTokenSecret     string
	AccessTokenExpSeconds int64
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

	applications, err := parseApplicationsConfig()
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

	adminPassword, err := getEnv("ADMIN_PASSWORD")
	if err != nil {
		return nil, err
	}

	accessTokenSecret, err := getEnv("ACCESS_TOKEN_SECRET")
	if err != nil {
		return nil, err
	}

	accessTokenExpSeconds, err := getEnvInt("ACCESS_TOKEN_EXP_SECONDS")
	if err != nil {
		return nil, err
	}

	return &EnvVariables{
		IsDevelopment:         environment == "development",
		IsProduction:          environment == "production",
		Port:                  port,
		DbConnectionString:    dbConnectionString,
		Applications:          applications,
		AWSS3LimitBytes:       awsS3LimitBytes,
		AWSS3BackupBucketName: awsS3BackupBucketName,
		RedisAddress:          redisAddress,
		RedisPassword:         redisPassword,
		AdminPassword:         adminPassword,
		AccessTokenSecret:     accessTokenSecret,
		AccessTokenExpSeconds: accessTokenExpSeconds,
	}, nil
}

func parseApplicationsConfig() ([]BackupAppConfig, error) {
	appNames, err := getEnv("APP_NAMES")
	if err != nil {
		return nil, err
	}

	names := strings.Split(appNames, ",")
	applications := make([]BackupAppConfig, 0, len(names))

	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		dbConnString, err := getEnv(fmt.Sprintf("APP_%s_DB_CONNECTION_STRING", strings.ToUpper(name)))
		if err != nil {
			return nil, err
		}

		// Public folder is optional
		publicFolderPath, err := getEnvOptional(fmt.Sprintf("APP_%s_PUBLIC_FOLDER_PATH", strings.ToUpper(name)))
		if err != nil {
			return nil, err
		}

		applications = append(applications, BackupAppConfig{
			Name:               name,
			DbConnectionString: dbConnString,
			PublicFolderPath:   publicFolderPath,
		})
	}

	return applications, nil
}

func getEnvOptional(key string) (string, error) {
	envVar := os.Getenv(key)
	return envVar, nil
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
