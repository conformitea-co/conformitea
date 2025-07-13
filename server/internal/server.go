package internal

import (
	"fmt"
	"os"
	"strings"

	"conformitea/server/internal/config"
	"conformitea/server/internal/database"
	"conformitea/server/internal/gateways"
	"conformitea/server/internal/logger"
	"conformitea/server/internal/middlewares"
	"conformitea/server/internal/routes"
	public "conformitea/server/types"

	"github.com/gin-gonic/gin"
)

type serverDependencies struct {
	router *gin.Engine
}

func (s *serverDependencies) GetRouter() *gin.Engine {
	return s.router
}

func (s *serverDependencies) Start() error {
	// Ensure logs are flushed on shutdown
	defer func() {
		if err := logger.GetLogger().Sync(); err != nil {
			// Ignore sync errors on stdout/stderr
			if err.Error() != "sync /dev/stdout: inappropriate ioctl for device" &&
				err.Error() != "sync /dev/stderr: inappropriate ioctl for device" {
				fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
			}
		}
	}()

	port := config.GetConfig().HTTPServer.Port
	if port == "" {
		port = "8080"
	}

	// Add colon prefix if not present
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	var err error
	if err = s.router.Run(port); err == nil {
		logger.GetLogger().Info(fmt.Sprintf("Server listening on port %s", port))
	}

	return err
}

var cftServer *serverDependencies

func Initialize(c public.Config) (public.Server, error) {
	if err := config.Initialize(c); err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	if err := logger.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := database.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	gateways.Initialize()

	cftServer = &serverDependencies{
		router: gin.New(),
	}

	if err := middlewares.RegisterMiddlewares(cftServer); err != nil {
		return nil, fmt.Errorf("failed to register middlewares: %w", err)
	}

	routes.RegisterRoutes(cftServer)

	return cftServer, nil
}
