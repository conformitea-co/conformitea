package middlewares

import (
	"conformitea/server/internal/types"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogRequestDetails(server types.InternalServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		var logger *zap.Logger
		if ctxLogger, exists := c.Get("logger"); exists {
			logger = ctxLogger.(*zap.Logger)
		} else {
			logger = server.GetLogger()
		}

		// Calculate latency
		latency := time.Since(start)
		latencyMS := latency.Milliseconds()

		status := c.Writer.Status()
		clientIP := c.ClientIP()

		// Get error if any
		var errorMessage string
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		requestDetailFields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.Int("status", status),
			zap.Int64("latency_ms", latencyMS),
			zap.String("ip", clientIP),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		msg := "HTTP request"

		if errorMessage != "" {
			requestDetailFields = append(requestDetailFields, zap.String("error", errorMessage))
		}

		switch {
		case status >= 500:
			logger.Error(msg, requestDetailFields...)
		case status >= 400:
			logger.Warn(msg, requestDetailFields...)
		default:
			logger.Info(msg, requestDetailFields...)
		}
	}
}
