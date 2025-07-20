package auth

import (
	"context"
	"fmt"
)

type CallbackRequest struct {
	Code                string
	State               string
	Nonce               string
	HydraLoginChallenge string
	Provider            string
}

type CallbackResult struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	Email        string
	Name         string
	Provider     string
}

// Process OAuth2 callback
func (a *Auth) ProcessCallback(ctx context.Context, req CallbackRequest) (CallbackResult, error) {
	// Route to appropriate IdP handler
	switch req.Provider {
	case "microsoft":
		return a.processMicrosoftCallback(ctx, req)
	default:
		return CallbackResult{}, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

// Processes Microsoft OAuth2 callback and completes Hydra flow
func (a *Auth) processMicrosoftCallback(ctx context.Context, req CallbackRequest) (CallbackResult, error) {
	// Exchange authorization code for access token
	token, err := a.msClient.ExchangeCodeForToken(ctx, req.Code)
	if err != nil {
		return CallbackResult{}, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user profile from Microsoft Graph
	userProfile, err := a.msClient.GetUserProfile(ctx, token)
	if err != nil {
		return CallbackResult{}, fmt.Errorf("failed to get user profile: %w", err)
	}

	hydraTokens, err := a.hydraClient.AcceptLoginSession(req.HydraLoginChallenge, userProfile.ID)
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
