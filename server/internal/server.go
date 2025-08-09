package internal

import (
	"fmt"
	"os"
	"strings"

	"conformitea/server/config"
	"conformitea/server/internal/handlers/auth"
	"conformitea/server/internal/handlers/users"
	"conformitea/server/internal/middlewares"
	"conformitea/server/internal/routes"
	"conformitea/server/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server struct {
	authHandlers *auth.AuthHandlers
	logger       *zap.Logger
	router       *gin.Engine
	config       config.Config
}

func (s *server) Start() error {
	// Ensure logs are flushed on shutdown
	defer func() {
		if err := s.logger.Sync(); err != nil {
			// Ignore sync errors on stdout/stderr
			if err.Error() != "sync /dev/stdout: inappropriate ioctl for device" &&
				err.Error() != "sync /dev/stderr: inappropriate ioctl for device" {
				fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
			}
		}
	}()

	port := s.config.HTTPServer.Port
	if port == "" {
		port = "8080"
	}

	// Add colon prefix if not present
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	if err := s.router.Run(port); err != nil {
		return err
	}

	return nil
}

func Initialize(c config.Config, l *zap.Logger, appAuth types.AppAuth) (types.Server, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid server configuration: %w", err)
	}

	router := gin.New()

	if err := middlewares.RegisterMiddlewares(router, l, c); err != nil {
		return nil, fmt.Errorf("failed to register middlewares: %w", err)
	}

	authHandlers := auth.Initialize(appAuth, c)
	usersHandlers := users.Initialize(c)
	routes.RegisterRoutes(router, authHandlers, usersHandlers)

	return &server{
		authHandlers: authHandlers,
		logger:       l,
		router:       router,
		config:       c,
	}, nil
}
