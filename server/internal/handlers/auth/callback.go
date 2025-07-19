package auth

import (
	"context"
	"net/http"

	appAuth "conformitea/app/auth"
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
		authErr := cftError.NewAuthError(cftError.AuthInvalidState, map[string]any{
			"parameter": "code",
			"reason":    "missing",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	session := sessions.Default(c)

	hydraLoginChallenge, exists := session.Get("hydra_login_challenge").(string)
	if !exists || hydraLoginChallenge == "" {
		authErr := cftError.NewAuthError(cftError.AuthSessionNotFound, map[string]any{
			"session_key": "hydra_login_challenge",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	provider, exists := session.Get("idp_provider").(string)
	if !exists || provider == "" {
		authErr := cftError.NewAuthError(cftError.AuthSessionNotFound, map[string]any{
			"session_key": "idp_provider",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	nonce, exists := session.Get("auth_nonce").(string)
	if !exists || nonce == "" {
		authErr := cftError.NewAuthError(cftError.AuthSessionNotFound, map[string]any{
			"session_key": "auth_nonce",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	req := appAuth.CallbackRequest{
		Code:                code,
		State:               state,
		Nonce:               nonce,
		HydraLoginChallenge: hydraLoginChallenge,
		Provider:            provider,
	}

	ctx := context.Background()
	result, err := appAuth.ProcessCallback(ctx, req)
	if err != nil {
		authErr := cftError.NewAuthErrorWithMessage(cftError.AuthMicrosoftExchange, err.Error(), map[string]any{
			"provider": provider,
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Store Hydra tokens and user data in session
	session.Set("access_token", result.AccessToken)
	session.Set("refresh_token", result.RefreshToken)
	session.Set("user_id", result.UserID)
	session.Set("email", result.Email)
	session.Set("name", result.Name)
	session.Set("provider", result.Provider)
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
