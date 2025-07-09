package routes

import (
	"conformitea/server/internal/handlers"
	"conformitea/server/internal/handlers/auth"
	"conformitea/server/internal/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(server types.InternalServer) {
	router := server.GetRouter()

	// Authentication routes
	router.GET("/auth/login", handlerWrapper(auth.Login, server))
	router.GET("/auth/callback", handlerWrapper(auth.Callback, server))
	router.GET("/auth/me", handlerWrapper(auth.Me, server))
	router.POST("/auth/logout", handlerWrapper(auth.Logout, server))

	// Health check
	router.GET("/ping", handlers.Ping)
}

func handlerWrapper(handler func(c *gin.Context, logger *zap.Logger), server types.InternalServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logger *zap.Logger

		// The following serves two purposes:
		// - To avoid the casting of the logger in every handler;
		// - To ensure that the logger is always available for the handler (although it should be set in the context by the middleware);
		if ctxLogger, exists := c.Get("logger"); exists {
			logger = ctxLogger.(*zap.Logger)
		} else {
			logger = server.GetLogger()
		}

		handler(c, logger)
	}
}
