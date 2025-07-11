package hydra

import "net/http"

// HydraClient represents a Hydra admin API client for managing OAuth2 sessions and tokens.
type HydraClient struct {
	adminURL   string
	httpClient *http.Client
}

// HydraLoginSession represents a Hydra OAuth2 login session request.
type HydraLoginSession struct {
	Challenge string `json:"challenge"`
	Client    struct {
		ClientId string `json:"client_id"`
	} `json:"client"`
	RequestURL     string   `json:"request_url"`
	Skip           bool     `json:"skip"`
	Subject        string   `json:"subject"`
	RequestedScope []string `json:"requested_scope"`
}

// TokenResponse represents Hydra's OAuth2 token response containing access and refresh tokens.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// AcceptLoginRequest represents the request payload for accepting a Hydra login session.
type AcceptLoginRequest struct {
	Subject     string `json:"subject"`
	Remember    bool   `json:"remember"`
	RememberFor int    `json:"remember_for"`
}

// TokenInfo represents the response from Hydra's token introspection endpoint.
type TokenInfo struct {
	Active bool   `json:"active"`
	Sub    string `json:"sub"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
	Scope  string `json:"scope"`
}
