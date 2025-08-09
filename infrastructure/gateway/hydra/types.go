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

type AcceptLoginRequest struct {
	Subject     string `json:"subject"`
	Remember    bool   `json:"remember"`
	RememberFor int    `json:"remember_for"`
}

type AcceptLoginResponse struct {
	RedirectTo string `json:"redirect_to"`
}

// TokenInfo represents the response from Hydra's token introspection endpoint.
type TokenInfo struct {
	Active bool   `json:"active"`
	Sub    string `json:"sub"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
	Scope  string `json:"scope"`
}

type HydraGetConsentResponse struct {
	Challenge string `json:"challenge"`
	Skip      bool   `json:"skip"`
	Subject   string `json:"subject"`
	Client    struct {
		ClientId   string `json:"client_id"`
		ClientName string `json:"client_name"`
	} `json:"client"`
	RequestURL                   string   `json:"request_url"`
	RequestedScope               []string `json:"requested_scope"`
	RequestedAccessTokenAudience []string `json:"requested_access_token_audience"`
}

type HydraConsentSessionTokens struct {
	AccessToken map[string]any `json:"access_token,omitempty"`
	IDToken     map[string]any `json:"id_token,omitempty"`
}

type HydraPutAcceptConsentRequest struct {
	GrantScope               []string                  `json:"grant_scope"`
	GrantAccessTokenAudience []string                  `json:"grant_access_token_audience,omitempty"`
	Remember                 bool                      `json:"remember"`
	RememberFor              int                       `json:"remember_for"`
	Session                  HydraConsentSessionTokens `json:"session"`
}

type HydraPutAcceptConsentResponse struct {
	RedirectTo string `json:"redirect_to"`
}
