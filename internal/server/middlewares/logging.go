package middlewares

import (
	"time"

	"github.com/conformitea-co/conformitea/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggingMiddleware creates a middleware for HTTP request logging
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		latencyMS := latency.Milliseconds()

		// Get status code
		status := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get error if any
		var errorMessage string
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		// Build log fields
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.Int("status", status),
			zap.Int64("latency_ms", latencyMS),
			zap.String("ip", clientIP),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// Add user ID if authenticated
		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		// Add session ID if exists
		if sessionID, exists := c.Get("session_id"); exists {
			fields = append(fields, zap.String("session_id", sessionID.(string)))
		}

		// Log based on status code
		msg := "HTTP request"

		if errorMessage != "" {
			fields = append(fields, zap.String("error", errorMessage))
		}

		switch {
		case status >= 500:
			logger.Error(msg, fields...)
		case status >= 400:
			logger.Warn(msg, fields...)
		case status >= 300:
			logger.Info(msg, fields...)
		default:
			logger.Info(msg, fields...)
		}
	}
}

// RequestIDMiddleware ensures a request ID is available in the context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists (from logging middleware)
		requestID, exists := c.Get("request_id")
		if !exists {
			requestID = uuid.New().String()
			c.Set("request_id", requestID)
		}

		// Add to response header for debugging
		c.Header("X-Request-ID", requestID.(string))

		c.Next()
	}
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}

// GetLogger returns a logger with request context
func GetLogger(c *gin.Context) *zap.Logger {
	fields := []zap.Field{
		zap.String("request_id", GetRequestID(c)),
	}

	// Add user ID if authenticated
	if userID, exists := c.Get("user_id"); exists {
		fields = append(fields, zap.String("user_id", userID.(string)))
	}

	// Add session ID if exists
	if sessionID, exists := c.Get("session_id"); exists {
		fields = append(fields, zap.String("session_id", sessionID.(string)))
	}

	return logger.WithContext(fields...)
}
