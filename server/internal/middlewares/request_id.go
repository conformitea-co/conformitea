package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestId string

		if cRequestId, exists := c.Get("request_id"); exists {
			requestId = cRequestId.(string)
		} else {
			requestId = uuid.New().String()
			c.Set("request_id", requestId)
		}

		if logger, exists := c.Get("logger"); exists {
			c.Set("logger", logger.(*zap.Logger).With(zap.String("request_id", requestId)))
		}

		// Add to response header for debugging
		c.Header("X-Request-ID", requestId)

		c.Next()
	}
}
