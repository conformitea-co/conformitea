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

type AppAuth interface {
	InitiateLogin(req LoginRequest) (LoginResult, error)
	ProcessCallback(ctx context.Context, req CallbackRequest) (CallbackResult, error)
	ProcessConsent(ctx context.Context, req ConsentRequest) (ConsentResult, error)
}
