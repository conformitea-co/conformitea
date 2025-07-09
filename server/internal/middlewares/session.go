package middlewares

import (
	"conformitea/server/internal/types"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware(server types.InternalServer) gin.HandlerFunc {
	redisStore := server.GetSessionStore()
	cookieName := server.GetConfig().HTTPServer.Session.CookieName

	return sessions.Sessions(cookieName, redisStore)
}
