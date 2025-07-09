//go:build wireinject

package server

import (
	"conformitea/server/config"
	"conformitea/server/internal/gateway/gin_session"
	"conformitea/server/internal/gateway/hydra"
	"conformitea/server/internal/gateway/microsoft"
	"conformitea/server/internal/logger"
	"conformitea/server/types"

	"github.com/google/wire"
)

// ProvideLoggerConfig extracts logger configuration from the main config.
func ProvideLoggerConfig(config config.Config) config.LoggerConfig {
	return config.Logger
}

// ProvideRedisConfig extracts Redis configuration from the main config.
func ProvideRedisConfig(config config.Config) config.RedisConfig {
	return config.Redis
}

// ProvideHTTPServerConfig extracts HTTP server configuration from the main config.
func ProvideHTTPServerConfig(config config.Config) config.HTTPServerConfig {
	return config.HTTPServer
}

// ProvideHydraConfig extracts Hydra configuration from the main config.
func ProvideHydraConfig(config config.Config) config.HydraConfig {
	return config.Hydra
}

// ProvideMicrosoftOAuthConfig extracts Microsoft OAuth configuration from the main config.
func ProvideMicrosoftOAuthConfig(config config.Config) config.MicrosoftOAuthConfig {
	return config.OAuth.Microsoft
}

// ServerSet contains all the providers needed to build a server
var ServerSet = wire.NewSet(
	ProvideLoggerConfig,
	ProvideRedisConfig,
	ProvideHTTPServerConfig,
	ProvideHydraConfig,
	ProvideMicrosoftOAuthConfig,
	logger.ProvideLogger,
	gin_session.ProvideRedisStore,
	hydra.ProvideHydraClient,
	microsoft.ProvideMicrosoftClient,
	Initialize,
)

// Creates a server instance with all dependencies injected
func InitializeServer(cfg config.Config) (types.PublicServer, error) {
	wire.Build(ServerSet)
	return nil, nil
}
