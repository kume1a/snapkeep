package webserver

import (
	"snapkeep/internal/config"
	"snapkeep/internal/logger"

	"github.com/gin-gonic/gin"
)

func ConfigureWebServer() error {
	envVars, err := config.ParseEnv()
	if err != nil {
		logger.Fatal("Failed to parse environment variables: ", err)
		return err
	}

	r := gin.Default()

	// healthcheck
	r.GET("/", HealthcheckHandler)

	// auth
	r.POST("/auth/signIn", LoginHandler)

	if envVars.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := r.Run(); err != nil {
		logger.Fatal("Failed to start HTTP server: ", err)
		return err
	}

	return nil
}
