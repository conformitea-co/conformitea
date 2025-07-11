package middlewares

import (
	"conformitea/server/internal/logger"

	"github.com/gin-gonic/gin"
)

func ContextLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", logger.GetLogger())

		c.Next()
	}
}
