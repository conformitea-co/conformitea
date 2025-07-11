package sections

import (
	"errors"
	"strings"
)

type GeneralConfig struct {
	FrontendURL string `mapstructure:"frontend_url"`
}

func (g *GeneralConfig) Validate() error {
	var errs []string
	if g.FrontendURL == "" {
		errs = append(errs, "general.frontend_url is required")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
