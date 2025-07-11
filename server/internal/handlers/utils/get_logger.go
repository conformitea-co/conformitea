package utils

import (
	"conformitea/server/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// This function serves two purposes:
//
// - To avoid the casting of the logger in every handler;
// - To ensure that the logger is always available for the handler
// (although it should be set in the context by the middleware);
func GetLogger(c *gin.Context) *zap.Logger {
	if ctxLogger, exists := c.Get("logger"); exists {
		return ctxLogger.(*zap.Logger)
	}
	return logger.GetLogger()
}
