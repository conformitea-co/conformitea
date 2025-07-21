package config

import (
	"errors"
	"fmt"
)

type SessionConfig struct {
	CookieName string   `mapstructure:"cookie_name"`
	KeyPairs   []string `mapstructure:"key_pairs"`
	Timeout    int      `mapstructure:"timeout"`
}

func (s *SessionConfig) Validate() error {
	var errs []error
	if s.CookieName == "" {
		errs = append(errs, fmt.Errorf("server.session.cookie_name is required"))
	}

	if len(s.KeyPairs) == 0 {
		errs = append(errs, fmt.Errorf("server.session.key_pairs is required and must not be empty"))
	}

	if s.Timeout <= 0 {
		errs = append(errs, fmt.Errorf("server.session.timeout must be positive"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
