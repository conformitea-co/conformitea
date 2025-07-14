package microsoft

import "golang.org/x/oauth2"

// OAuth2 configuration.
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// OAuth2 client for authentication flows.
type OAuthClient struct {
	config oauth2.Config
}

// Microsoft user profile from Graph API.
type MicrosoftUserProfile struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	UserPrincipalName string `json:"userPrincipalName"`
	Mail              string `json:"mail"`
}
