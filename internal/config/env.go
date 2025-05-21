package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	// godotenv.Load(".env." + env + ".local")
	// if "test" != env {
	//   godotenv.Load(".env.local")
	// }

	if env == "development" {
		envPath := ".env." + env

		log.Println("Loading env file: " + envPath)

		godotenv.Load(envPath)
	}
}

type EnvVariables struct {
	IsDevelopment      bool
	IsProduction       bool
	Port               string
	DbConnectionString string
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

	return &EnvVariables{
		IsDevelopment:      environment == "development",
		IsProduction:       environment == "production",
		Port:               port,
		DbConnectionString: dbConnectionString,
	}, nil
}

// func getEnvInt(key string) (int64, error) {
// 	val, err := getEnv(key)
// 	if err != nil {
// 		return 0, err
// 	}

// 	valInt, err := strconv.ParseInt(val, 10, 64)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return valInt, nil
// }

func getEnv(key string) (string, error) {
	envVar := os.Getenv(key)
	if envVar == "" {
		msg := fmt.Sprintf("%v is not found in the env", key)

		log.Fatal(msg)
		return "", errors.New(msg)
	}

	return envVar, nil
}
