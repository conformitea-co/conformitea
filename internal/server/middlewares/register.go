package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterMiddlewares(router *gin.Engine) {
	sessionMiddleware, err := SessionMiddleware()
	if err != nil {
		log.Fatalf("failed to initialize session middleware: %v", err)
	}
	router.Use(sessionMiddleware)
}
