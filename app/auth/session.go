package auth

import (
	"context"
	"fmt"
	"strings"

	"conformitea/server/types"
)

// ExchangeCodeForTokens exchanges an authorization code for JWT tokens from Hydra
func (a *Auth) ExchangeCodeForTokens(ctx context.Context, req types.SessionRequest) (types.SessionResult, error) {
	// Exchange authorization code with Hydra
	tokenResponse, err := a.hydraClient.ExchangeAuthorizationCode(req.Code)
	if err != nil {
		return types.SessionResult{}, fmt.Errorf("failed to exchange authorization code: %w", err)
	}

	// Parse the JWT to extract claims
	claims, err := a.hydraClient.ParseJWT(tokenResponse.AccessToken)
	if err != nil {
		return types.SessionResult{}, fmt.Errorf("failed to parse JWT: %w", err)
	}

	// Extract JTI (JWT ID) from claims
	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return types.SessionResult{}, fmt.Errorf("JWT missing jti claim")
	}

	// Extract subject (user ID) from claims
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return types.SessionResult{}, fmt.Errorf("JWT missing sub claim")
	}

	// Extract expiration from claims
	exp, ok := claims["exp"].(float64)
	if !ok {
		return types.SessionResult{}, fmt.Errorf("JWT missing exp claim")
	}

	// If ID token is present, it might have additional user info
	// For now, we'll use the access token claims
	userID := sub

	// Clean up the user ID if it contains provider prefix
	if idx := strings.LastIndex(userID, ":"); idx > 0 {
		userID = userID[idx+1:]
	}

	return types.SessionResult{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		JTI:          jti,
		UserID:       userID,
		ExpiresAt:    int64(exp),
	}, nil
}
