package config

import (
	"net/http"
	"snapkeep/pkg/logger"

	"github.com/gin-gonic/gin"
)

func ConfigureWebServer() error {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	if err := r.Run(); err != nil {
		logger.Fatal("Failed to start HTTP server: ", err)
		return err
	}

	return nil
}
