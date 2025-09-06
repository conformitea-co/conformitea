package auth

import (
	"conformitea/server/config"
	"conformitea/server/types"

	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	appAuth types.AppAuth
	config  config.Config
}

func Initialize(routerGroup *gin.RouterGroup, appAuth types.AppAuth, config config.Config) {
	handlers := &AuthHandlers{
		appAuth: appAuth,
		config:  config,
	}

	routerGroup.GET("/callback", handlers.Callback)
	routerGroup.GET("/consent", handlers.Consent)
	routerGroup.GET("/login", handlers.Login)
	routerGroup.POST("/logout", handlers.Logout)
	routerGroup.POST("/session", handlers.Session)
}
