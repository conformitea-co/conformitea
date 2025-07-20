package auth

import "conformitea/server/types"

type AuthHandlers struct {
	appAuth types.AppAuth
}

func Initialize(appAuth types.AppAuth) *AuthHandlers {
	return &AuthHandlers{
		appAuth: appAuth,
	}
}
