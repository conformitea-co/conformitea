package microsoft

import (
	"conformitea/server/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

// ProvideMicrosoftClient creates a new Microsoft OAuth client instance based on the provided configuration.
// This replaces the singleton Initialize pattern with proper dependency injection.
func ProvideMicrosoftClient(cfg config.MicrosoftOAuthConfig) *OAuthClient {
	client := &OAuthClient{
		config: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     microsoft.AzureADEndpoint("common"),
		},
	}

	return client
}
