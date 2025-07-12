package sections

import (
	"errors"
	"strings"
)

type HydraConfig struct {
	AdminURL  string `mapstructure:"admin_url"`
	PublicURL string `mapstructure:"public_url"`
}

func (h *HydraConfig) Validate() error {
	var errs []string
	if h.AdminURL == "" {
		errs = append(errs, "hydra.admin_url is required")
	}

	if h.PublicURL == "" {
		errs = append(errs, "hydra.public_url is required")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
