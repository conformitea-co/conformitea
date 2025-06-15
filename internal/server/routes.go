package server

import (
	"github.com/conformitea-co/conformitea/internal/server/handlers"
	"github.com/conformitea-co/conformitea/internal/server/handlers/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Authentication routes
	router.GET("/auth/login", auth.Login)
	router.GET("/auth/callback", auth.Callback)
	router.GET("/auth/me", auth.Me)
	router.POST("/auth/logout", auth.Logout)

	// Health check
	router.GET("/ping", handlers.Ping)
}
