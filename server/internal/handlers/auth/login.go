package auth

import (
	"net/http"

	appAuth "conformitea/app/auth"
	cftError "conformitea/server/internal/error"
	"conformitea/server/internal/handlers/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handles the initial login request from Hydra and routes to appropriate IdP.
func Login(c *gin.Context) {
	logger := utils.GetLogger(c)

	// Extract login_challenge from Hydra
	loginChallenge := c.Query("login_challenge")
	if loginChallenge == "" {
		authErr := cftError.NewAuthError(cftError.AuthInvalidState, map[string]any{
			"parameter": "login_challenge",
			"reason":    "missing",
		})

		logger.Warn("login attempt without challenge",
			zap.String("error_code", string(authErr.Code)),
		)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	logger.Info("login initiated",
		zap.String("login_challenge", loginChallenge),
	)

	// Use application layer to handle login logic
	req := appAuth.LoginRequest{
		LoginChallenge: loginChallenge,
	}

	result, err := appAuth.InitiateLogin(req)
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthSessionCreateFailed, err.Error(), map[string]interface{}{
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
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthSessionCreateFailed, err.Error(), nil)

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
