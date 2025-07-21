package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ContextLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", logger)

		c.Next()
	}
}
