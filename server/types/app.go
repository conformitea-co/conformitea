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
	AccessToken  string
	RefreshToken string
	UserID       string
	Email        string
	Name         string
	Provider     string
}

type AppAuth interface {
	InitiateLogin(req LoginRequest) (LoginResult, error)
	ProcessCallback(ctx context.Context, req CallbackRequest) (CallbackResult, error)
}
