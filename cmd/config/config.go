package config

import (
	infrastructure "conformitea/infrastructure/config"
	server "conformitea/server/config"
)

type Config struct {
	GeneralConfig    server.GeneralConfig    `mapstructure:"general"`
	HTTPServerConfig server.HTTPServerConfig `mapstructure:"server"`
	RedisConfig      server.RedisConfig      `mapstructure:"redis"`

	LoggerConfig   infrastructure.LoggerConfig   `mapstructure:"logger"`
	DatabaseConfig infrastructure.DatabaseConfig `mapstructure:"database"`
	HydraConfig    infrastructure.HydraConfig    `mapstructure:"hydra"`
	OAuthConfig    infrastructure.OAuthConfig    `mapstructure:"oauth"`
}
