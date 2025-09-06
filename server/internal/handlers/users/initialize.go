package users

import (
	"conformitea/server/config"

	"github.com/gin-gonic/gin"
)

type UsersHandlers struct {
	config config.Config
}

func Initialize(routerGroup *gin.RouterGroup, config config.Config) {
	handlers := &UsersHandlers{
		config: config,
	}

	routerGroup.GET("/users/me", handlers.Me)
}
