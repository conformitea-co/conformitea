package internal

import (
	"fmt"
	"os"
	"strings"

	"conformitea/server/internal/config"
	"conformitea/server/internal/handlers/auth"
	"conformitea/server/internal/middlewares"
	"conformitea/server/internal/routes"
	"conformitea/server/types"
	public "conformitea/server/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server struct {
	authHandlers *auth.AuthHandlers
	logger       *zap.Logger
	router       *gin.Engine
}

func (s *server) GetRouter() *gin.Engine {
	return s.router
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

	port := config.GetConfig().HTTPServer.Port
	if port == "" {
		port = "8080"
	}

	// Add colon prefix if not present
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	if err := s.router.Run(port); err == nil {
		s.logger.Info(fmt.Sprintf("Server listening on port %s", port))
		return nil
	} else {
		return err
	}
}

func Initialize(c public.Config, l *zap.Logger, appAuth types.AppAuth) (public.Server, error) {
	if err := config.Initialize(c); err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	router := gin.New()

	if err := middlewares.RegisterMiddlewares(router, l); err != nil {
		return nil, fmt.Errorf("failed to register middlewares: %w", err)
	}

	authHandlers := auth.Initialize(appAuth)
	routes.RegisterRoutes(router, authHandlers)

	return &server{
		authHandlers: authHandlers,
		logger:       l,
		router:       router,
	}, nil
}
