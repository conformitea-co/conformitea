package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"conformitea/infrastructure/gateway/hydra"
	"conformitea/infrastructure/gateway/microsoft"
)

// LoginRequest represents the input for login flow initiation
type LoginRequest struct {
	LoginChallenge string
}

// LoginResult represents the complete result of login flow initiation
type LoginResult struct {
	// URL to redirect user to for authentication
	AuthURL string
	// Session data to store for callback processing
	HydraLoginChallenge string
	IDPProvider         string
	AuthNonce           string
}

// Creates a secure random nonce for OAuth2 state.
func generateNonce() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// InitiateLogin handles the business logic for starting the OAuth2 login flow
func InitiateLogin(req LoginRequest) (LoginResult, error) {
	if req.LoginChallenge == "" {
		return LoginResult{}, fmt.Errorf("login_challenge parameter is required")
	}

	// Get Hydra client and validate login session
	hydraClient, err := hydra.GetHydraClient()
	if err != nil {
		return LoginResult{}, fmt.Errorf("failed to get Hydra client: %w", err)
	}

	loginSession, err := hydraClient.GetLoginSession(req.LoginChallenge)
	if err != nil {
		return LoginResult{}, fmt.Errorf("failed to get Hydra login session: %w", err)
	}

	// Determine IdP based on client_id
	provider := loginSession.Client.ClientId
	if provider != "microsoft" {
		return LoginResult{}, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Generate nonce for security
	nonce, err := generateNonce()
	if err != nil {
		return LoginResult{}, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Get Microsoft OAuth client
	microsoftClient, err := microsoft.GetOAuthClient()
	if err != nil {
		return LoginResult{}, fmt.Errorf("failed to get Microsoft OAuth client: %w", err)
	}

	// Generate OAuth URL
	authURL, err := microsoftClient.GenerateAuthURL(req.LoginChallenge, nonce)
	if err != nil {
		return LoginResult{}, fmt.Errorf("failed to generate Microsoft OAuth URL: %w", err)
	}

	return LoginResult{
		AuthURL:             authURL,
		HydraLoginChallenge: req.LoginChallenge,
		IDPProvider:         provider,
		AuthNonce:           nonce,
	}, nil
}
