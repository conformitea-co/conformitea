package config

import (
	"errors"
)

var BUILD = "development"

type Config struct {
	LoggerConfig   LoggerConfig   `mapstructure:"logger"`
	DatabaseConfig DatabaseConfig `mapstructure:"database"`
	HydraConfig    HydraConfig    `mapstructure:"hydra"`
	OAuthConfig    OAuthConfig    `mapstructure:"oauth"`
}

func (c *Config) Validate() error {
	var errs []error

	if err := c.LoggerConfig.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.DatabaseConfig.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.HydraConfig.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.OAuthConfig.Validate(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
