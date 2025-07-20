package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogRequestDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		latencyMS := latency.Milliseconds()

		status := c.Writer.Status()
		clientIP := c.ClientIP()

		requestDetailFields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.Int("status", status),
			zap.Int64("latency_ms", latencyMS),
			zap.String("ip", clientIP),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// Get error if any
		var errorMessage string
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
			requestDetailFields = append(requestDetailFields, zap.String("error", errorMessage))
		}

		msg := "http request details"
		logger := c.MustGet("logger").(*zap.Logger)

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
