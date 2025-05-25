package webserver

import (
	"snapkeep/pkg/logger"

	"github.com/gin-gonic/gin"
)

func ConfigureWebServer() error {
	r := gin.Default()

	// healthcheck
	r.GET("/", HealthcheckHandler)

	// auth
	r.POST("/auth/signIn", LoginHandler)

	if err := r.Run(); err != nil {
		logger.Fatal("Failed to start HTTP server: ", err)
		return err
	}

	return nil
}
