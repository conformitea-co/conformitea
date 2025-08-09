package auth

import (
	"context"
	"fmt"

	"conformitea/server/types"
)

// Process OAuth2 callback
func (a *Auth) ProcessCallback(ctx context.Context, req types.CallbackRequest) (types.CallbackResult, error) {
	// Route to appropriate IdP handler
	switch req.Provider {
	case "microsoft":
		return a.processMicrosoftCallback(ctx, req)
	default:
		return types.CallbackResult{}, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

// Processes Microsoft OAuth2 callback and completes Hydra flow
func (a *Auth) processMicrosoftCallback(ctx context.Context, req types.CallbackRequest) (types.CallbackResult, error) {
	// Exchange authorization code for access token
	token, err := a.msClient.ExchangeCodeForToken(ctx, req.Code)
	if err != nil {
		return types.CallbackResult{}, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user profile from Microsoft Graph
	userProfile, err := a.msClient.GetUserProfile(ctx, token)
	if err != nil {
		return types.CallbackResult{}, fmt.Errorf("failed to get user profile: %w", err)
	}

	result, err := a.hydraClient.AcceptLoginSession(req.HydraLoginChallenge, userProfile.ID)
	if err != nil {
		return types.CallbackResult{}, fmt.Errorf("failed to accept hydra login session: %w", err)
	}

	email := userProfile.Mail
	if email == "" {
		email = userProfile.UserPrincipalName
	}

	return types.CallbackResult{
		RedirectTo: result.RedirectTo,
	}, nil
}
