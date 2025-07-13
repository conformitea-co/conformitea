package config

import (
	"errors"

	"conformitea/server/types"
)

var BUILD = "development"

var roConfig types.Config

func Initialize(c types.Config) error {
	var errs []error
	if err := c.General.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.HTTPServer.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.Database.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.Redis.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.Hydra.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.OAuth.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.Logger.Validate(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	roConfig = c

	return nil
}

func GetConfig() types.Config {
	return roConfig
}
