package server

import (
	"fmt"
	"os"

	"github.com/conformitea-co/conformitea/internal/logger"
	"github.com/conformitea-co/conformitea/internal/server/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewRouter() *gin.Engine {
	router := gin.New()

	// Use custom logging middleware instead of gin.Logger()
	router.Use(middlewares.LoggingMiddleware(), gin.Recovery())

	middlewares.RegisterMiddlewares(router)
	RegisterRoutes(router)

	return router
}

func Start() error {
	// Initialize logger
	if err := logger.Initialize(viper.GetViper()); err != nil {
		// Fallback to fmt if logger fails to initialize
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		return err
	}

	// Log server startup
	logger.Info("Starting ConformiTea server",
		logger.Fields(map[string]interface{}{
			"version":     "1.0.0",
			"environment": viper.GetString("environment"),
		})...,
	)

	// Ensure logs are flushed on shutdown
	defer func() {
		if err := logger.Sync(); err != nil {
			// Ignore sync errors on stdout/stderr
			if err.Error() != "sync /dev/stdout: inappropriate ioctl for device" &&
				err.Error() != "sync /dev/stderr: inappropriate ioctl for device" {
				fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
			}
		}
	}()

	router := NewRouter()

	port := viper.GetString("server.port")
	if port == "" {
		port = ":8080"
	}

	logger.Info("Server listening", logger.Fields(map[string]interface{}{
		"port": port,
	})...)

	return router.Run(port)
}
