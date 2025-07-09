package middlewares

import (
	"conformitea/server/internal/types"

	"github.com/gin-gonic/gin"
)

func ContextLoggerMiddleware(server types.InternalServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", server.GetLogger())

		c.Next()
	}
}
