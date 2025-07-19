package infrastructure

import (
	"fmt"

	"conformitea/infrastructure/config"
	"conformitea/infrastructure/database"
	"conformitea/infrastructure/gateway/hydra"
	"conformitea/infrastructure/gateway/microsoft"
	"conformitea/infrastructure/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	config          config.Config
	logger          *zap.Logger
	database        *gorm.DB
	hydraClient     *hydra.HydraClient
	microsoftClient *microsoft.OAuthClient
}

var container Container

func Initialize(lc config.LoggerConfig, dc config.DatabaseConfig, hc config.HydraConfig, oc config.OAuthConfig) (Container, error) {
	l, err := logger.Initialize(lc)
	if err != nil {
		return Container{}, fmt.Errorf("failed to initialize logger: %w", err)
	}

	db, err := database.Initialize(dc, l)
	if err != nil {
		return Container{}, fmt.Errorf("failed to initialize database: %w", err)
	}

	h, err := hydra.Initialize(hc)
	if err != nil {
		return Container{}, fmt.Errorf("failed to initialize Hydra client: %w", err)
	}

	ms, err := microsoft.Initialize(oc.Microsoft)
	if err != nil {
		return Container{}, fmt.Errorf("failed to initialize Microsoft OAuth client: %w", err)
	}

	container.config = config.Config{
		LoggerConfig:   lc,
		DatabaseConfig: dc,
		HydraConfig:    hc,
		OAuthConfig:    oc,
	}
	container.logger = l
	container.database = db
	container.hydraClient = h
	container.microsoftClient = ms

	return container, nil
}

func (c *Container) GetLogger() *zap.Logger {
	return c.logger
}

func (c *Container) GetDatabase() *gorm.DB {
	return c.database
}

func (c *Container) GetHydraClient() *hydra.HydraClient {
	return c.hydraClient
}

func (c *Container) GetMicrosoftClient() *microsoft.OAuthClient {
	return c.microsoftClient
}
