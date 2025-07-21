package config

import (
	"errors"
)

type DatabaseConfig struct {
	URL                string `mapstructure:"url"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
}

func (d DatabaseConfig) Validate() error {
	var errs []error

	if d.URL == "" {
		errs = append(errs, errors.New("database URL is required"))
	}

	if d.MaxOpenConnections <= 0 {
		errs = append(errs, errors.New("max_open_connections must be >= 1"))
	}

	if d.MaxIdleConnections < 0 {
		errs = append(errs, errors.New("max_idle_connections must be non-negative"))
	}

	if d.MaxIdleConnections > d.MaxOpenConnections && d.MaxOpenConnections > 0 {
		errs = append(errs, errors.New("max_idle_connections cannot exceed max_open_connections"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
