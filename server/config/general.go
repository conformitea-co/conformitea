package config

import (
	"errors"
	"fmt"
)

type GeneralConfig struct {
	FrontendURL string `mapstructure:"frontend_url"`
}

func (g *GeneralConfig) Validate() error {
	var errs []error
	if g.FrontendURL == "" {
		errs = append(errs, fmt.Errorf("general.frontend_url is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
