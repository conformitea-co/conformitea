package auth

import (
	"conformitea/server/config"
	"conformitea/server/types"
)

type AuthHandlers struct {
	appAuth types.AppAuth
	config  config.Config
}

func Initialize(appAuth types.AppAuth, cfg config.Config) *AuthHandlers {
	return &AuthHandlers{
		appAuth: appAuth,
		config:  cfg,
	}
}
