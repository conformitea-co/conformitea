package auth

import (
	"net/http"

	"conformitea/server/internal/cerror"
	"conformitea/server/types"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handles the initial login request from Hydra and routes to appropriate IdP.
func (a *AuthHandlers) Login(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)

	// Extract login_challenge from Hydra
	loginChallenge := c.Query("login_challenge")
	if loginChallenge == "" {
		authErr := cerror.NewAuthError(cerror.AuthInvalidState, map[string]any{
			"parameter": "login_challenge",
			"reason":    "missing",
		})

		logger.Warn("login attempt without challenge",
			zap.String("error_code", string(authErr.Code)),
		)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	logger.Info("login initiated", zap.String("login_challenge", loginChallenge))

	result, err := a.appAuth.InitiateLogin(types.LoginRequest{LoginChallenge: loginChallenge})
	if err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthSessionCreateFailed, err.Error(), map[string]any{
			"login_challenge": loginChallenge,
		})

		logger.Error("failed to initiate login",
			zap.String("login_challenge", loginChallenge),
			zap.Error(err),
			zap.String("error_code", string(authErr.Code)),
		)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Store auth info in session for callback handler
	session := sessions.Default(c)
	session.Set("hydra_login_challenge", result.HydraLoginChallenge)
	session.Set("idp_provider", result.IDPProvider)
	session.Set("auth_nonce", result.AuthNonce)

	if err := session.Save(); err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthSessionCreateFailed, err.Error(), nil)

		logger.Error("failed to save session",
			zap.Error(err),
			zap.String("error_code", string(authErr.Code)),
		)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	logger.Info("redirecting to OAuth2 provider",
		zap.String("provider", result.IDPProvider),
		zap.String("login_challenge", loginChallenge),
	)

	c.Redirect(http.StatusFound, result.AuthURL)
}
