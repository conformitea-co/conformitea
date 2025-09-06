package types

import "context"

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

type CallbackRequest struct {
	Code                string
	State               string
	Nonce               string
	HydraLoginChallenge string
	Provider            string
}

type CallbackResult struct {
	RedirectTo string
}

type ConsentRequest struct {
	ConsentChallenge string
}

type ConsentResult struct {
	RedirectTo string
}

type SessionRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state"`
}

type SessionResult struct {
	AccessToken  string
	RefreshToken string
	JTI          string // JWT ID from the token
	UserID       string
	ExpiresAt    int64 // Unix timestamp
}

// JWTClaims represents parsed JWT claims returned by the application layer.
type JWTClaims struct {
	UserID      string                 `json:"user_id"`
	JTI         string                 `json:"jti"`
	Exp         int64                  `json:"exp"`
	Roles       []string               `json:"roles,omitempty"`
	Permissions []string               `json:"permissions,omitempty"`
	RawClaims   map[string]interface{} `json:"raw_claims,omitempty"`
}

type AppAuth interface {
	InitiateLogin(req LoginRequest) (LoginResult, error)
	ProcessCallback(ctx context.Context, req CallbackRequest) (CallbackResult, error)
	ProcessConsent(ctx context.Context, req ConsentRequest) (ConsentResult, error)
	ExchangeCodeForTokens(ctx context.Context, req SessionRequest) (SessionResult, error)
	ValidateJWT(ctx context.Context, token string) (*JWTClaims, error)
}
