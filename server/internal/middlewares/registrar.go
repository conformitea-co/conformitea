package middlewares

import (
	"fmt"

	"conformitea/server/internal/types"

	"github.com/gin-gonic/gin"
)

func RegisterMiddlewares(server types.Server) error {
	router := server.GetRouter()
	sessionMiddleware, err := SessionMiddleware()
	if err != nil {
		return fmt.Errorf("failed to initialize session middleware: %w", err)
	}

	// Most of the time, the order of middlewares is important.
	router.Use(LogRequestDetails())
	router.Use(ContextLoggerMiddleware())
	router.Use(RequestIdMiddleware())
	router.Use(CORSMiddleware())
	router.Use(sessionMiddleware)
	router.Use(gin.Recovery())

	return nil
}
