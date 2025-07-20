package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterMiddlewares(router *gin.Engine, logger *zap.Logger) error {
	sessionMiddleware, err := SessionMiddleware()
	if err != nil {
		return fmt.Errorf("failed to initialize session middleware: %w", err)
	}

	// Most of the time, the order of middlewares is important.
	router.Use(LogRequestDetails())
	router.Use(ContextLoggerMiddleware(logger))
	router.Use(RequestIdMiddleware())
	router.Use(CORSMiddleware())
	router.Use(sessionMiddleware)
	router.Use(gin.Recovery())

	return nil
}
