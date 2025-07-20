package auth

import (
	"context"
	"net/http"

	"conformitea/server/internal/cerror"
	"conformitea/server/internal/config"
	"conformitea/server/types"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handles OAuth2 callbacks from identity providers and completes the Hydra flow.
func (a *AuthHandlers) Callback(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)
	logger.Info("processing oauth2 callback")

	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		logger.Warn("oauth2 callback without code")

		authErr := cerror.NewAuthError(cerror.AuthInvalidState, map[string]any{
			"parameter": "code",
			"reason":    "missing",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	session := sessions.Default(c)

	hydraLoginChallenge, exists := session.Get("hydra_login_challenge").(string)
	if !exists || hydraLoginChallenge == "" {
		logger.Warn("oauth2 callback without hydra login challenge")

		authErr := cerror.NewAuthError(cerror.AuthSessionNotFound, map[string]any{
			"session_key": "hydra_login_challenge",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	provider, exists := session.Get("idp_provider").(string)
	if !exists || provider == "" {
		logger.Warn("oauth2 callback without identity provider")

		authErr := cerror.NewAuthError(cerror.AuthSessionNotFound, map[string]any{
			"session_key": "idp_provider",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	nonce, exists := session.Get("auth_nonce").(string)
	if !exists || nonce == "" {
		logger.Warn("oauth2 callback without auth nonce")

		authErr := cerror.NewAuthError(cerror.AuthSessionNotFound, map[string]any{
			"session_key": "auth_nonce",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	req := types.CallbackRequest{
		Code:                code,
		State:               state,
		Nonce:               nonce,
		HydraLoginChallenge: hydraLoginChallenge,
		Provider:            provider,
	}

	ctx := context.Background()
	result, err := a.appAuth.ProcessCallback(ctx, req)
	if err != nil {
		logger.Error("failed to process oauth2 callback", zap.Error(err))

		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthMicrosoftExchange, err.Error(), map[string]any{
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
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthSessionCreateFailed, err.Error(), nil)

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	c.Redirect(http.StatusFound, config.GetConfig().General.FrontendURL)
}
