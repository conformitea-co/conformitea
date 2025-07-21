package config

import (
	"errors"
	"fmt"
)

type OAuthConfig struct {
	Microsoft MicrosoftOAuthConfig `mapstructure:"microsoft"`
}

func (o *OAuthConfig) Validate() error {
	return o.Microsoft.Validate()
}

type MicrosoftOAuthConfig struct {
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	RedirectURL  string   `mapstructure:"redirect_url"`
	Scopes       []string `mapstructure:"scopes"`
}

func (m *MicrosoftOAuthConfig) Validate() error {
	var errs []error
	if m.ClientID == "" {
		errs = append(errs, fmt.Errorf("oauth.microsoft.client_id is required"))
	}

	if m.ClientSecret == "" {
		errs = append(errs, fmt.Errorf("oauth.microsoft.client_secret is required"))
	}

	if m.RedirectURL == "" {
		errs = append(errs, fmt.Errorf("oauth.microsoft.redirect_url is required"))
	}

	if len(m.Scopes) == 0 {
		errs = append(errs, fmt.Errorf("oauth.microsoft.scopes is required and must not be empty"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
