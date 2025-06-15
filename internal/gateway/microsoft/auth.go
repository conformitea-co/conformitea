// Package microsoft provides OAuth2 integration with Microsoft Graph API.
package microsoft

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

// Config represents Microsoft OAuth2 configuration.
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// Client represents a Microsoft OAuth2 client for authentication flows.
type Client struct {
	config *oauth2.Config
}

// NewClient creates a new Microsoft OAuth2 client with the provided configuration.
func NewClient(cfg Config) *Client {
	return &Client{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     microsoft.AzureADEndpoint("common"),
		},
	}
}

// UserProfile represents a Microsoft user profile from Graph API.
type UserProfile struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	UserPrincipalName string `json:"userPrincipalName"`
	Mail              string `json:"mail"`
}

// GenerateAuthURL creates a Microsoft OAuth2 authorization URL with state and nonce parameters.
func (c *Client) GenerateAuthURL(state, nonce string) string {
	return c.config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("nonce", nonce),
		oauth2.SetAuthURLParam("response_mode", "query"),
	)
}

// ExchangeCodeForToken exchanges an authorization code for an OAuth2 token.
func (c *Client) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}

	return token, nil
}

// GetUserProfile retrieves the user's profile information from Microsoft Graph API.
func (c *Client) GetUserProfile(ctx context.Context, token *oauth2.Token) (*UserProfile, error) {
	client := c.config.Client(ctx, token)

	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("microsoft Graph API error: status %d", resp.StatusCode)
	}

	var profile UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode user profile: %w", err)
	}

	return &profile, nil
}
