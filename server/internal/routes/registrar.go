package routes

import (
	"conformitea/server/internal/handlers"
	"conformitea/server/internal/handlers/auth"
	"conformitea/server/internal/types"
)

func RegisterRoutes(s types.Server) {
	router := s.GetRouter()

	// Authentication routes
	router.GET("/auth/callback", auth.Callback)
	router.GET("/auth/login", auth.Login)
	router.GET("/auth/me", auth.Me)
	router.POST("/auth/logout", auth.Logout)

	// Health check
	router.GET("/ping", handlers.Ping)
}
