package auth

import (
	"net/http"

	"conformitea/server/internal/cerror"
	"conformitea/server/types"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Session handles the exchange of authorization code for JWT tokens
func (a *AuthHandlers) Session(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)
	logger.Info("processing session creation")

	var req types.SessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("invalid session request", zap.Error(err))

		authErr := cerror.NewAuthError(cerror.AuthInvalidRequest, map[string]any{
			"error": err.Error(),
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	if req.Code == "" {
		logger.Warn("session request without authorization code")

		authErr := cerror.NewAuthError(cerror.AuthInvalidState, map[string]any{
			"parameter": "code",
			"reason":    "missing",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Exchange authorization code for tokens
	result, err := a.appAuth.ExchangeCodeForTokens(c.Request.Context(), req)
	if err != nil {
		logger.Error("failed to exchange code for tokens",
			zap.Error(err),
			zap.String("code", req.Code[:min(10, len(req.Code))]+"...")) // Log partial code for debugging

		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthTokenExchangeFailed, err.Error(), map[string]any{
			"operation": "token_exchange",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Store refresh token and metadata in session
	session := sessions.Default(c)
	session.Set("refresh_token", result.RefreshToken)
	session.Set("jti", result.JTI)
	session.Set("user_id", result.UserID)
	session.Set("expires_at", result.ExpiresAt)
	session.Set("authenticated", true)

	if err := session.Save(); err != nil {
		logger.Error("failed to save session", zap.Error(err))

		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthSessionCreateFailed, err.Error(), nil)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Set JWT as HTTP-only cookie
	a.setAccessTokenCookie(c, result.AccessToken)

	logger.Info("session created successfully",
		zap.String("user_id", result.UserID),
		zap.String("jti", result.JTI))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Session created successfully",
	})
}

// setAccessTokenCookie sets the JWT access token as an HTTP-only cookie
func (a *AuthHandlers) setAccessTokenCookie(c *gin.Context, token string) {
	// Determine if we're in development mode
	secure := a.config.General.Environment == "production"

	// Set cookie with appropriate security settings
	cookie := &http.Cookie{
		Name:     "cft_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure, // Only HTTPS in production
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600, // 1 hour
	}

	http.SetCookie(c.Writer, cookie)
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
