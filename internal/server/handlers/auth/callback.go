package auth

import (
	"context"
	"net/http"

	cerror "github.com/conformitea-co/conformitea/internal/error"
	"github.com/conformitea-co/conformitea/internal/gateway/hydra"
	"github.com/conformitea-co/conformitea/internal/gateway/microsoft"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Callback handles OAuth2 callbacks from identity providers and completes the Hydra flow.
func Callback(c *gin.Context) {
	// Extract authorization code and state from query parameters
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		authErr := cerror.NewAuthError(cerror.AuthInvalidState, map[string]interface{}{
			"parameter": "code",
			"reason":    "missing",
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Retrieve auth info from session
	session := sessions.Default(c)
	hydraLoginChallenge, exists := session.Get("hydra_login_challenge").(string)
	if !exists || hydraLoginChallenge == "" {
		authErr := cerror.NewAuthError(cerror.AuthSessionNotFound, map[string]interface{}{
			"session_key": "hydra_login_challenge",
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	provider, exists := session.Get("idp_provider").(string)
	if !exists || provider == "" {
		authErr := cerror.NewAuthError(cerror.AuthSessionNotFound, map[string]interface{}{
			"session_key": "idp_provider",
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	nonce, exists := session.Get("auth_nonce").(string)
	if !exists || nonce == "" {
		authErr := cerror.NewAuthError(cerror.AuthSessionNotFound, map[string]interface{}{
			"session_key": "auth_nonce",
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Route to appropriate IdP handler
	switch provider {
	case "microsoft":
		c.Set("auth_context", map[string]interface{}{
			"code":                  code,
			"state":                 state,
			"nonce":                 nonce,
			"hydra_login_challenge": hydraLoginChallenge,
		})
		handleMicrosoftCallback(c)
	default:
		authErr := cerror.NewAuthError(cerror.AuthProviderNotSupported, map[string]interface{}{
			"provider": provider,
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}
}

// handleMicrosoftCallback processes Microsoft OAuth2 callback and completes Hydra flow.
func handleMicrosoftCallback(c *gin.Context) {
	authContext := c.MustGet("auth_context").(map[string]interface{})
	code := authContext["code"].(string)
	hydraLoginChallenge := authContext["hydra_login_challenge"].(string)

	// Create Microsoft OAuth2 client
	microsoftConfig := microsoft.Config{
		ClientID:     viper.GetString("oauth.microsoft.client_id"),
		ClientSecret: viper.GetString("oauth.microsoft.client_secret"),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       viper.GetStringSlice("oauth.microsoft.scopes"),
	}
	microsoftClient := microsoft.NewClient(microsoftConfig)

	// Exchange authorization code for access token
	ctx := context.Background()
	token, err := microsoftClient.ExchangeCodeForToken(ctx, code)
	if err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthMicrosoftExchange, err.Error(), map[string]interface{}{
			"provider": "microsoft",
			"step":     "token_exchange",
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Get user profile from Microsoft Graph
	userProfile, err := microsoftClient.GetUserProfile(ctx, token)
	if err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthMicrosoftProfile, err.Error(), map[string]interface{}{
			"provider": "microsoft",
			"step":     "profile_fetch",
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Accept Hydra login session
	hydraClient := hydra.NewClient(viper.GetString("hydra.admin_url"))
	hydraTokens, err := hydraClient.AcceptLoginSession(hydraLoginChallenge, userProfile.ID)
	if err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthHydraAcceptFailed, err.Error(), map[string]interface{}{
			"login_challenge": hydraLoginChallenge,
			"user_id":         userProfile.ID,
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Store Hydra tokens and user data in session
	session := sessions.Default(c)
	session.Set("access_token", hydraTokens.AccessToken)
	session.Set("refresh_token", hydraTokens.RefreshToken)
	session.Set("user_id", userProfile.ID)
	session.Set("email", getEmailFromProfile(userProfile))
	session.Set("name", userProfile.DisplayName)
	session.Set("provider", "microsoft")
	session.Set("authenticated", true)

	// Clear temporary auth data
	session.Delete("hydra_login_challenge")
	session.Delete("idp_provider")
	session.Delete("auth_nonce")

	if err := session.Save(); err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthSessionCreateFailed, err.Error(), nil)
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Redirect to frontend
	c.Redirect(http.StatusFound, viper.GetString("general.frontend_url"))
}

// getEmailFromProfile extracts email from Microsoft user profile with fallback.
func getEmailFromProfile(profile *microsoft.UserProfile) string {
	if profile.Mail != "" {
		return profile.Mail
	}
	return profile.UserPrincipalName
}
