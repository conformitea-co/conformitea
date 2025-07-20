package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type LoginRequest struct {
	LoginChallenge string
}

type LoginResult struct {
	// URL to redirect user to for authentication
	AuthURL string
	// Session data to store for callback processing
	HydraLoginChallenge string
	IDPProvider         string
	AuthNonce           string
}

// Creates a secure random nonce for OAuth2 state.
func (a *Auth) generateNonce() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// Initiate OAuth2 login flow
func (a *Auth) InitiateLogin(req LoginRequest) (LoginResult, error) {
	if req.LoginChallenge == "" {
		return LoginResult{}, fmt.Errorf("login_challenge parameter is required")
	}

	loginSession, err := a.hydraClient.GetLoginSession(req.LoginChallenge)
	if err != nil {
		return LoginResult{}, fmt.Errorf("failed to get Hydra login session: %w", err)
	}

	// Determine IdP based on client_id
	provider := loginSession.Client.ClientId
	if provider != "microsoft" {
		return LoginResult{}, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Generate nonce for security
	nonce, err := a.generateNonce()
	if err != nil {
		return LoginResult{}, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Generate OAuth URL
	authURL, err := a.msClient.GenerateAuthURL(req.LoginChallenge, nonce)
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
