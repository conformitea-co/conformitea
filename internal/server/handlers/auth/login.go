package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	cerror "github.com/conformitea-co/conformitea/internal/error"
	"github.com/conformitea-co/conformitea/internal/gateway/hydra"
	"github.com/conformitea-co/conformitea/internal/gateway/microsoft"
	"github.com/conformitea-co/conformitea/internal/server/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// generateNonce creates a secure random nonce for OAuth2 state.
func generateNonce() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Login handles the initial login request from Hydra and routes to appropriate IdP.
func Login(c *gin.Context) {
	logger := middlewares.GetLogger(c)

	// Extract login_challenge from Hydra
	loginChallenge := c.Query("login_challenge")

	if loginChallenge == "" {
		authErr := cerror.NewAuthError(cerror.AuthInvalidState, map[string]interface{}{
			"parameter": "login_challenge",
			"reason":    "missing",
		})
		logger.Warn("Login attempt without challenge",
			zap.String("error_code", string(authErr.Code)),
		)
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	logger.Info("Login initiated",
		zap.String("login_challenge", loginChallenge),
	)

	// Create Hydra client and get login session details
	hydraClient := hydra.NewClient(viper.GetString("hydra.admin_url"))
	loginSession, err := hydraClient.GetLoginSession(loginChallenge)
	if err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthSessionNotFound, err.Error(), map[string]interface{}{
			"login_challenge": loginChallenge,
		})
		logger.Error("Failed to get Hydra login session",
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
		authErr := cerror.NewAuthError(cerror.AuthProviderNotSupported, map[string]interface{}{
			"provider": provider,
		})
		logger.Warn("Unsupported provider requested",
			zap.String("provider", provider),
			zap.String("error_code", string(authErr.Code)),
		)
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Generate nonce for security
	nonce, err := generateNonce()
	if err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthSessionCreateFailed, err.Error(), nil)
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Store auth info in session for callback handler
	session := sessions.Default(c)
	session.Set("hydra_login_challenge", loginChallenge)
	session.Set("idp_provider", provider)
	session.Set("auth_nonce", nonce)
	if err := session.Save(); err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthSessionCreateFailed, err.Error(), nil)
		logger.Error("Failed to save session",
			zap.Error(err),
			zap.String("error_code", string(authErr.Code)),
		)
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Create Microsoft OAuth2 client and generate auth URL
	microsoftConfig := microsoft.Config{
		ClientID:     viper.GetString("oauth.microsoft.client_id"),
		ClientSecret: viper.GetString("oauth.microsoft.client_secret"),
		RedirectURL:  viper.GetString("oauth.microsoft.redirect_url"),
		Scopes:       viper.GetStringSlice("oauth.microsoft.scopes"),
	}
	microsoftClient := microsoft.NewClient(microsoftConfig)
	authURL := microsoftClient.GenerateAuthURL(loginChallenge, nonce)

	// Redirect user to Microsoft OAuth2
	logger.Info("Redirecting to OAuth2 provider",
		zap.String("provider", provider),
		zap.String("login_challenge", loginChallenge),
	)
	c.Redirect(http.StatusFound, authURL)
}
