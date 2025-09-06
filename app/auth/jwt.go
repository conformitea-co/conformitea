package auth

import (
	"context"
	"fmt"
	"time"

	"conformitea/server/types"
)

// ValidateJWT validates and parses a JWT token using the Hydra client (which performs signature validation)
// and returns a types.JWTClaims value so the app implements the shared interface.
func (a *Auth) ValidateJWT(ctx context.Context, token string) (*types.JWTClaims, error) {
	claimsMap, err := a.hydraClient.ParseJWT(token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate JWT: %w", err)
	}

	// Extract standard claims
	sub, ok := claimsMap["sub"].(string)
	if !ok || sub == "" {
		return nil, fmt.Errorf("JWT missing sub claim")
	}

	jti, ok := claimsMap["jti"].(string)
	if !ok || jti == "" {
		return nil, fmt.Errorf("JWT missing jti claim")
	}

	expFloat, ok := claimsMap["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("JWT missing exp claim")
	}
	exp := int64(expFloat)

	// Check expiration
	if time.Now().Unix() > exp {
		return nil, fmt.Errorf("JWT token expired")
	}

	// Extract roles
	var roles []string
	if rolesRaw, ok := claimsMap["roles"].([]interface{}); ok {
		for _, r := range rolesRaw {
			if rs, ok := r.(string); ok {
				roles = append(roles, rs)
			}
		}
	}

	// Extract permissions
	var permissions []string
	if permsRaw, ok := claimsMap["permissions"].([]interface{}); ok {
		for _, p := range permsRaw {
			if ps, ok := p.(string); ok {
				permissions = append(permissions, ps)
			}
		}
	}

	// Convert map claims to map[string]interface{} for RawClaims
	raw := make(map[string]interface{}, len(claimsMap))
	for k, v := range claimsMap {
		raw[k] = v
	}

	return &types.JWTClaims{
		UserID:      sub,
		JTI:         jti,
		Exp:         exp,
		Roles:       roles,
		Permissions: permissions,
		RawClaims:   raw,
	}, nil
}
