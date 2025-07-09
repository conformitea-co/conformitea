package config

import (
	"errors"
	"fmt"
	"strings"
)

func (c *Config) Validate() error {
	var errs []error
	if err := c.General.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.HTTPServer.Validate(); err != nil {
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

	return nil
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

func (s *SessionConfig) Validate() error {
	var errs []string
	if s.CookieName == "" {
		errs = append(errs, "server.session.cookie_name is required")
	}

	if len(s.KeyPairs) == 0 {
		errs = append(errs, "server.session.key_pairs is required and must not be empty")
	}

	if s.Timeout <= 0 {
		errs = append(errs, "server.session.timeout must be positive")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func (l *LoggerConfig) Validate() error {
	var errs []string
	if l.Level == "" {
		errs = append(errs, "logger.level is required")
	}

	if l.Format == "" {
		errs = append(errs, "logger.format is required")
	}

	if l.Output == "" {
		errs = append(errs, "logger.output is required")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func (r *RedisConfig) Validate() error {
	if r.Address == "" {
		return fmt.Errorf("redis.address is required")
	}

	return nil
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

func (o *OAuthConfig) Validate() error {
	return o.Microsoft.Validate()
}

func (m *MicrosoftOAuthConfig) Validate() error {
	var errs []string
	if m.ClientID == "" {
		errs = append(errs, "oauth.microsoft.client_id is required")
	}

	if m.ClientSecret == "" {
		errs = append(errs, "oauth.microsoft.client_secret is required")
	}

	if m.RedirectURL == "" {
		errs = append(errs, "oauth.microsoft.redirect_url is required")
	}

	if len(m.Scopes) == 0 {
		errs = append(errs, "oauth.microsoft.scopes is required and must not be empty")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
