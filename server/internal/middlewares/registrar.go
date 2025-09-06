package middlewares

import (
	"fmt"

	"conformitea/server/config"
	"conformitea/server/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterCommonMiddlewares(r *gin.Engine, l *zap.Logger, c config.Config) error {
	sessionMiddleware, err := SessionMiddleware(c)
	if err != nil {
		return fmt.Errorf("failed to initialize session middleware: %w", err)
	}

	// Most of the time, the order of middlewares is important.
	r.Use(LogRequestDetails())
	r.Use(ContextLoggerMiddleware(l))
	r.Use(RequestIdMiddleware())
	r.Use(CORSMiddleware())
	r.Use(sessionMiddleware)
	r.Use(gin.Recovery())

	return nil
}

func RegisterAPIGroupMiddlewares(g *gin.RouterGroup, appAuth types.AppAuth) {
	g.Use(AuthMiddleware(appAuth))
}
