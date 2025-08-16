package users

import (
	"conformitea/server/config"
)

type UsersHandlers struct {
	config config.Config
}

func Initialize(cfg config.Config) *UsersHandlers {
	return &UsersHandlers{
		config: cfg,
	}
}
