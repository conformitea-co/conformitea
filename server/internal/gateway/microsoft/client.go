// Package microsoft provides OAuth2 integration with Microsoft Graph API.
package microsoft

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"conformitea/server/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	client *OAuthClient
	once   sync.Once
)

func Initialize(cfg config.MicrosoftOAuthConfig) {
	once.Do(func() {
		client = &OAuthClient{
			config: oauth2.Config{
				ClientID:     cfg.ClientID,
				ClientSecret: cfg.ClientSecret,
				RedirectURL:  cfg.RedirectURL,
				Scopes:       cfg.Scopes,
				Endpoint:     microsoft.AzureADEndpoint("common"),
			},
		}
	})
}

// GetOAuthClient returns the initialized Microsoft OAuth2 client.
func GetOAuthClient() (*OAuthClient, error) {
	if client == nil {
		return nil, fmt.Errorf("microsoft OAuth client is not initialized")
	}

	return client, nil
}

// Creates a Microsoft OAuth2 authorization URL with state and nonce parameters.
func (c *OAuthClient) GenerateAuthURL(state, nonce string) (string, error) {
	if state == "" {
		return "", fmt.Errorf("state parameter cannot be empty")
	}

	if nonce == "" {
		return "", fmt.Errorf("nonce parameter cannot be empty")
	}

	return c.config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("nonce", nonce),
		oauth2.SetAuthURLParam("response_mode", "query"),
	), nil
}

// Exchanges an authorization code for an OAuth2 token.
func (c *OAuthClient) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}

	return token, nil
}

// Retrieves the user's profile information from Microsoft Graph API.
func (c *OAuthClient) GetUserProfile(ctx context.Context, token *oauth2.Token) (*MicrosoftUserProfile, error) {
	client := c.config.Client(ctx, token)

	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("microsoft Graph API error: status %d", resp.StatusCode)
	}

	var profile MicrosoftUserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode user profile: %w", err)
	}

	return &profile, nil
}
