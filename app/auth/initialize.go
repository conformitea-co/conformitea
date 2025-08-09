package auth

import (
	"conformitea/domain/user"
	"conformitea/infrastructure/gateway/hydra"
	"conformitea/infrastructure/gateway/microsoft"

	"gorm.io/gorm"
)

type Auth struct {
	db          *gorm.DB
	userService *user.UserService
	msClient    *microsoft.OAuthClient
	hydraClient *hydra.HydraClient
}

func Initialize(db *gorm.DB, us *user.UserService, mc *microsoft.OAuthClient, hc *hydra.HydraClient) *Auth {
	return &Auth{
		db:          db,
		userService: us,
		msClient:    mc,
		hydraClient: hc,
	}
}
