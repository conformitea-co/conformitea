package types

import "conformitea/server/internal/config/sections"

type Config struct {
	General    sections.GeneralConfig    `mapstructure:"general"`
	HTTPServer sections.HTTPServerConfig `mapstructure:"server"`
	Redis      sections.RedisConfig      `mapstructure:"redis"`
}
