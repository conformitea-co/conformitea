package config

import (
	"errors"
)

type HydraConfig struct {
	AdminURL    string               `mapstructure:"admin_url"`
	PublicURL   string               `mapstructure:"public_url"`
	RedirectURL string               `mapstructure:"redirect_url"`
	Microsoft   HydraMicrosoftClient `mapstructure:"microsoft"`
}

type HydraMicrosoftClient struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

func (h *HydraConfig) Validate() error {
	var errs []error
	if h.AdminURL == "" {
		errs = append(errs, errors.New("hydra.admin_url is required"))
	}

	if h.PublicURL == "" {
		errs = append(errs, errors.New("hydra.public_url is required"))
	}

	if h.RedirectURL == "" {
		errs = append(errs, errors.New("hydra.redirect_url is required"))
	}

	if h.Microsoft.ClientID == "" {
		errs = append(errs, errors.New("hydra.microsoft.client_id is required"))
	}

	if h.Microsoft.ClientSecret == "" {
		errs = append(errs, errors.New("hydra.microsoft.client_secret is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
