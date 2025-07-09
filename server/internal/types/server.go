package types

import (
	"conformitea/server/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InternalServer interface {
	GetConfig() config.Config
	GetLogger() *zap.Logger
	GetRouter() *gin.Engine
	GetSessionStore() sessions.Store
}
