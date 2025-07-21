package config

import (
	"errors"
)

type Config struct {
	General    GeneralConfig    `mapstructure:"general"`
	HTTPServer HTTPServerConfig `mapstructure:"server"`
	Redis      RedisConfig      `mapstructure:"redis"`
}

func (c *Config) Validate() error {
	var errs []error

	if err := c.General.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.HTTPServer.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.Redis.Validate(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
