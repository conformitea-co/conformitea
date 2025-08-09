package routes

import (
	"conformitea/server/internal/handlers"
	"conformitea/server/internal/handlers/auth"
	"conformitea/server/internal/handlers/users"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, auth *auth.AuthHandlers, users *users.UsersHandlers) {
	// Authentication routes
	router.GET("/auth/callback", auth.Callback)
	router.GET("/auth/consent", auth.Consent)
	router.GET("/auth/login", auth.Login)
	router.POST("/auth/logout", auth.Logout)

	// User routes
	router.GET("/users/me", users.Me)

	// Health check
	router.GET("/ping", handlers.Ping)
}
