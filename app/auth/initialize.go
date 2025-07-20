package auth

import (
	"conformitea/domain/user"
	"conformitea/infrastructure/gateway/hydra"
	"conformitea/infrastructure/gateway/microsoft"
)

type Auth struct {
	userService *user.UserService
	msClient    *microsoft.OAuthClient
	hydraClient *hydra.HydraClient
}

func NewAuth(us *user.UserService, mc *microsoft.OAuthClient, hc *hydra.HydraClient) *Auth {
	return &Auth{
		userService: us,
		msClient:    mc,
		hydraClient: hc,
	}
}
