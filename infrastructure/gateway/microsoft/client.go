// Package microsoft provides OAuth2 integration with Microsoft Graph API.
package microsoft

import (
	"conformitea/infrastructure/config"
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func Initialize(msConfigValues config.MicrosoftOAuthConfig) (*OAuthClient, error) {
	if err := msConfigValues.Validate(); err != nil {
		return nil, fmt.Errorf("invalid Microsoft OAuth configuration: %w", err)
	}

	client := &OAuthClient{
		config: oauth2.Config{
			ClientID:     msConfigValues.ClientID,
			ClientSecret: msConfigValues.ClientSecret,
			RedirectURL:  msConfigValues.RedirectURL,
			Scopes:       msConfigValues.Scopes,
			Endpoint:     microsoft.AzureADEndpoint("common"),
		},
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
