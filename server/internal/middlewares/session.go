package middlewares

import (
	"fmt"

	"conformitea/server/config"
	"conformitea/server/internal/gateway/gin_session"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware(cfg config.Config) (gin.HandlerFunc, error) {
	sessionStore, err := gin_session.NewStore(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize session store: %w", err)
	}

	cookieName := cfg.HTTPServer.Session.CookieName

	return sessions.Sessions(cookieName, sessionStore), nil
}
