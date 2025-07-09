package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	cftError "conformitea/server/internal/error"
	"conformitea/server/internal/gateway/hydra"
	cftMicrosoft "conformitea/server/internal/gateway/microsoft"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Creates a secure random nonce for OAuth2 state.
func generateNonce() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// Handles the initial login request from Hydra and routes to appropriate IdP.
func Login(c *gin.Context, logger *zap.Logger) {
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

	hydraClient, err := hydra.GetHydraClient()
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthSessionNotFound, err.Error(), map[string]interface{}{
			"login_challenge": loginChallenge,
		})

		logger.Error("failed to get Hydra client",
			zap.String("login_challenge", loginChallenge),
			zap.Error(err),
			zap.String("error_code", string(authErr.Code)),
		)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	loginSession, err := hydraClient.GetLoginSession(loginChallenge)
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthSessionNotFound, err.Error(), map[string]interface{}{
			"login_challenge": loginChallenge,
		})

		logger.Error("failed to get Hydra login session",
			zap.String("login_challenge", loginChallenge),
			zap.Error(err),
			zap.String("error_code", string(authErr.Code)),
		)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Determine IdP based on client_id
	provider := loginSession.Client.ClientId
	if provider != "microsoft" {
		authErr := cftError.NewAuthError(cftError.AuthProviderNotSupported, map[string]interface{}{
			"provider": provider,
		})

		logger.Warn("unsupported provider requested",
			zap.String("provider", provider),
			zap.String("error_code", string(authErr.Code)),
		)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Generate nonce for security
	nonce, err := generateNonce()
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthSessionCreateFailed, err.Error(), nil)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Store auth info in session for callback handler
	session := sessions.Default(c)
	session.Set("hydra_login_challenge", loginChallenge)
	session.Set("idp_provider", provider)
	session.Set("auth_nonce", nonce)

	if err := session.Save(); err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthSessionCreateFailed, err.Error(), nil)

		logger.Error("failed to save session",
			zap.Error(err),
			zap.String("error_code", string(authErr.Code)),
		)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	microsoftClient, err := cftMicrosoft.GetOAuthClient()
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthSessionCreateFailed, err.Error(), nil)

		logger.Error("failed to get Microsoft OAuth client",
			zap.Error(err),
			zap.String("error_code", string(authErr.Code)),
		)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	authURL, err := microsoftClient.GenerateAuthURL(loginChallenge, nonce)
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthSessionCreateFailed, err.Error(), nil)

		logger.Error("failed to generate Microsoft OAuth URL",
			zap.Error(err),
			zap.String("login_challenge", loginChallenge),
			zap.String("error_code", string(authErr.Code)),
		)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	logger.Info("redirecting to OAuth2 provider",
		zap.String("provider", provider),
		zap.String("login_challenge", loginChallenge),
	)

	c.Redirect(http.StatusFound, authURL)
}
