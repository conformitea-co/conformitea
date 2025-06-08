package server

import (
	"github.com/conformitea-co/conformitea/internal/server/handlers"
	"github.com/conformitea-co/conformitea/internal/server/handlers/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/auth/login", auth.Login)

	router.GET("/ping", handlers.Ping)
}
