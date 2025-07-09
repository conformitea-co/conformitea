package server

import (
	"fmt"
	"os"
	"strings"

	cftConfig "conformitea/server/config"
	"conformitea/server/internal/gateway/hydra"
	"conformitea/server/internal/gateway/microsoft"
	"conformitea/server/internal/middlewares"
	"conformitea/server/internal/routes"
	"conformitea/server/types"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server struct {
	config          cftConfig.Config
	hydraClient     *hydra.HydraClient
	logger          *zap.Logger
	microsoftClient *microsoft.OAuthClient
	redisStore      sessions.Store
	router          *gin.Engine
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

	var err error
	if err = s.router.Run(port); err == nil {
		s.logger.Info(fmt.Sprintf("Server listening on port %s", port))
	}

	return err
}

func (s *server) GetConfig() cftConfig.Config {
	return s.config
}

func (s *server) GetLogger() *zap.Logger {
	return s.logger
}

func (s *server) GetRouter() *gin.Engine {
	return s.router
}

func (s *server) GetSessionStore() sessions.Store {
	return s.redisStore
}

var cftServer *server

func Initialize(
	config cftConfig.Config,
	logger *zap.Logger,
	redisStore sessions.Store,
	hydraClient *hydra.HydraClient,
	microsoftClient *microsoft.OAuthClient,
) (types.PublicServer, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	logger.Info("starting ConformiTea server",
		zap.String("version", "0.1.0"),
	)

	cftServer = &server{
		config:          config,
		hydraClient:     hydraClient,
		logger:          logger,
		microsoftClient: microsoftClient,
		redisStore:      redisStore,
		router:          gin.New(),
	}

	middlewares.RegisterMiddlewares(cftServer)
	routes.RegisterRoutes(cftServer)

	return cftServer, nil
}
