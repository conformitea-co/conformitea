package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterMiddlewares(router *gin.Engine) {
	corsMiddleware, _ := CORSMiddleware()
	router.Use(corsMiddleware)

	sessionMiddleware, err := SessionMiddleware()
	if err != nil {
		log.Fatalf("failed to initialize session middleware: %v", err)
	}
	router.Use(sessionMiddleware)
}
