package server

import (
	"github.com/conformitea-co/conformitea/internal/server/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/ping", handlers.Ping)
}
