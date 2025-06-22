package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	cerror "github.com/conformitea-co/conformitea/internal/error"
	"github.com/conformitea-co/conformitea/internal/gateway/hydra"
	"github.com/conformitea-co/conformitea/internal/gateway/microsoft"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	// Extract login_challenge from Hydra
	loginChallenge := c.Query("login_challenge")

	if loginChallenge == "" {
		authErr := cerror.NewAuthError(cerror.AuthInvalidState, map[string]interface{}{
			"parameter": "login_challenge",
			"reason":    "missing",
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Create Hydra client and get login session details
	hydraClient := hydra.NewClient(viper.GetString("hydra.admin_url"))
	loginSession, err := hydraClient.GetLoginSession(loginChallenge)
	if err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthSessionNotFound, err.Error(), map[string]interface{}{
			"login_challenge": loginChallenge,
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Determine IdP based on client_id
	provider := loginSession.Client.ClientId
	if provider != "microsoft" {
		authErr := cerror.NewAuthError(cerror.AuthProviderNotSupported, map[string]interface{}{
			"provider": provider,
		})
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
	c.Redirect(http.StatusFound, authURL)
}
