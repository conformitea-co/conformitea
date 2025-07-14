package auth

import (
	"context"
	"net/http"

	"conformitea/infrastructure/gateway/hydra"
	"conformitea/infrastructure/gateway/microsoft"
	"conformitea/server/internal/config"
	cftError "conformitea/server/internal/error"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handles OAuth2 callbacks from identity providers and completes the Hydra flow.
func Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		authErr := cftError.NewAuthError(cftError.AuthInvalidState, map[string]interface{}{
			"parameter": "code",
			"reason":    "missing",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	session := sessions.Default(c)

	hydraLoginChallenge, exists := session.Get("hydra_login_challenge").(string)
	if !exists || hydraLoginChallenge == "" {
		authErr := cftError.NewAuthError(cftError.AuthSessionNotFound, map[string]interface{}{
			"session_key": "hydra_login_challenge",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	provider, exists := session.Get("idp_provider").(string)
	if !exists || provider == "" {
		authErr := cftError.NewAuthError(cftError.AuthSessionNotFound, map[string]interface{}{
			"session_key": "idp_provider",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	nonce, exists := session.Get("auth_nonce").(string)
	if !exists || nonce == "" {
		authErr := cftError.NewAuthError(cftError.AuthSessionNotFound, map[string]interface{}{
			"session_key": "auth_nonce",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	authContext := map[string]string{
		"code":                  code,
		"state":                 state,
		"nonce":                 nonce,
		"hydra_login_challenge": hydraLoginChallenge,
	}

	// Route to appropriate IdP handler
	switch provider {
	case "microsoft":
		handleMicrosoftCallback(c, authContext)
	default:
		authErr := cftError.NewAuthError(cftError.AuthProviderNotSupported, map[string]any{
			"provider": provider,
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}
}

// Processes Microsoft OAuth2 callback and completes Hydra flow.
func handleMicrosoftCallback(c *gin.Context, authContext map[string]string) {
	microsoftClient, err := microsoft.GetOAuthClient()
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthHydraClientInit, err.Error(), nil)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Exchange authorization code for access token
	ctx := context.Background()
	token, err := microsoftClient.ExchangeCodeForToken(ctx, authContext["code"])
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthMicrosoftExchange, err.Error(), map[string]any{
			"provider": "microsoft",
			"step":     "token_exchange",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Get user profile from Microsoft Graph
	userProfile, err := microsoftClient.GetUserProfile(ctx, token)
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthMicrosoftProfile, err.Error(), map[string]any{
			"provider": "microsoft",
			"step":     "profile_fetch",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	hydraClient, err := hydra.GetHydraClient()
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthHydraClientInit, err.Error(), nil)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	hydraTokens, err := hydraClient.AcceptLoginSession(authContext["hydra_login_challenge"], userProfile.ID)
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthHydraAcceptFailed, err.Error(), map[string]any{
			"login_challenge": authContext["hydra_login_challenge"],
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
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthSessionCreateFailed, err.Error(), nil)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	c.Redirect(http.StatusFound, config.GetConfig().General.FrontendURL)
}

// Extracts email from Microsoft user profile with fallback.
func getEmailFromProfile(profile *microsoft.MicrosoftUserProfile) string {
	if profile.Mail != "" {
		return profile.Mail
	}

	return profile.UserPrincipalName
}
