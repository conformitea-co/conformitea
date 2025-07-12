package middlewares

import (
	"fmt"

	"conformitea/server/internal/config"
	"conformitea/server/internal/gateways/gin_session"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware() (gin.HandlerFunc, error) {
	sessionStore, err := gin_session.NewStore()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize session store: %w", err)
	}

	cookieName := config.GetConfig().HTTPServer.Session.CookieName

	return sessions.Sessions(cookieName, sessionStore), nil
}
