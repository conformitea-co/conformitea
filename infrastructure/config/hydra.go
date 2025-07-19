package config

import (
	"errors"
)

type HydraConfig struct {
	AdminURL  string `mapstructure:"admin_url"`
	PublicURL string `mapstructure:"public_url"`
}

func (h *HydraConfig) Validate() error {
	var errs []error
	if h.AdminURL == "" {
		errs = append(errs, errors.New("hydra.admin_url is required"))
	}

	if h.PublicURL == "" {
		errs = append(errs, errors.New("hydra.public_url is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
