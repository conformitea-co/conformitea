package middlewares

import (
	"net/http"
	"strings"

	"conformitea/server/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(appAuth types.AppAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := c.MustGet("logger").(*zap.Logger)

		// Try to get JWT from cookie first
		token := ""
		if cookie, err := c.Request.Cookie("cft_token"); err == nil && cookie != nil {
			token = cookie.Value
		}

		// If no cookie, try Authorization header (for API clients)
		if token == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					token = parts[1]
				}
			}
		}

		if token == "" {
			logger.Warn("no JWT token found in request")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		claims, err := appAuth.ValidateJWT(c.Request.Context(), token)
		if err != nil {
			logger.Warn("invalid JWT token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Add claims to context
		c.Set("user_id", claims.UserID)
		c.Set("jti", claims.JTI)
		c.Set("jwt_claims", claims.RawClaims)
		c.Set("authenticated", true)

		c.Next()
	}
}
