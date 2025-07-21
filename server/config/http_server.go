package config

import (
	"errors"
	"fmt"
)

type HTTPServerConfig struct {
	Port    string        `mapstructure:"port"`
	Session SessionConfig `mapstructure:"session"`
}

func (h *HTTPServerConfig) Validate() error {
	var errs []error
	if h.Port == "" {
		errs = append(errs, fmt.Errorf("server.port is required"))
	}

	if err := h.Session.Validate(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
