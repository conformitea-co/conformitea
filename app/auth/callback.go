package auth

import (
	"context"
	"fmt"

	"conformitea/infrastructure/gateway/hydra"
	"conformitea/infrastructure/gateway/microsoft"
)

// CallbackRequest represents the OAuth callback parameters
type CallbackRequest struct {
	Code                string
	State               string
	Nonce               string
	HydraLoginChallenge string
	Provider            string
}

// CallbackResult represents the result of processing OAuth callback
type CallbackResult struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	Email        string
	Name         string
	Provider     string
}

// ProcessCallback handles the business logic for OAuth2 callback processing
func ProcessCallback(ctx context.Context, req CallbackRequest) (CallbackResult, error) {
	// Route to appropriate IdP handler
	switch req.Provider {
	case "microsoft":
		return processMicrosoftCallback(ctx, req)
	default:
		return CallbackResult{}, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

// Processes Microsoft OAuth2 callback and completes Hydra flow
func processMicrosoftCallback(ctx context.Context, req CallbackRequest) (CallbackResult, error) {
	microsoftClient, err := microsoft.GetOAuthClient()
	if err != nil {
		return CallbackResult{}, fmt.Errorf("failed to get Microsoft OAuth client: %w", err)
	}

	// Exchange authorization code for access token
	token, err := microsoftClient.ExchangeCodeForToken(ctx, req.Code)
	if err != nil {
		return CallbackResult{}, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user profile from Microsoft Graph
	userProfile, err := microsoftClient.GetUserProfile(ctx, token)
	if err != nil {
		return CallbackResult{}, fmt.Errorf("failed to get user profile: %w", err)
	}

	hydraClient, err := hydra.GetHydraClient()
	if err != nil {
		return CallbackResult{}, fmt.Errorf("failed to get Hydra client: %w", err)
	}

	hydraTokens, err := hydraClient.AcceptLoginSession(req.HydraLoginChallenge, userProfile.ID)
	if err != nil {
		return CallbackResult{}, fmt.Errorf("failed to accept Hydra login session: %w", err)
	}

	// Extract email from profile with fallback
	email := userProfile.Mail
	if email == "" {
		email = userProfile.UserPrincipalName
	}

	return CallbackResult{
		AccessToken:  hydraTokens.AccessToken,
		RefreshToken: hydraTokens.RefreshToken,
		UserID:       userProfile.ID,
		Email:        email,
		Name:         userProfile.DisplayName,
		Provider:     "microsoft",
	}, nil
}
