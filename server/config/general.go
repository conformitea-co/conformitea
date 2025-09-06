package config

import (
	"errors"
	"fmt"
	"os"
)

type GeneralConfig struct {
	Environment string
}

func (g *GeneralConfig) Validate() error {
	var errs []error

	// Get environment from ENV variable, default to development
	g.Environment = os.Getenv("ENVIRONMENT")
	if g.Environment == "" {
		g.Environment = "development"
	}

	// Validate environment value
	if g.Environment != "development" && g.Environment != "production" && g.Environment != "staging" {
		errs = append(errs, fmt.Errorf("ENVIRONMENT must be one of: development, staging, production"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
