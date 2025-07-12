package types

import "conformitea/server/internal/config/sections"

type Config struct {
	General    sections.GeneralConfig    `mapstructure:"general"`
	HTTPServer sections.HTTPServerConfig `mapstructure:"server"`
	Hydra      sections.HydraConfig      `mapstructure:"hydra"`
	Logger     sections.LoggerConfig     `mapstructure:"logger"`
	OAuth      sections.OAuthConfig      `mapstructure:"oauth"`
	Redis      sections.RedisConfig      `mapstructure:"redis"`
}
