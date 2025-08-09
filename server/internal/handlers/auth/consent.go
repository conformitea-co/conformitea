package auth

import (
	"net/http"

	"conformitea/server/internal/cerror"
	"conformitea/server/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (a *AuthHandlers) Consent(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)
	logger.Info("processing oauth2 consent")

	consentChallenge := c.Query("consent_challenge")
	if consentChallenge == "" {
		logger.Warn("consent request without consent_challenge")

		authErr := cerror.NewAuthError(cerror.AuthInvalidState, map[string]any{
			"parameter": "consent_challenge",
			"reason":    "missing",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	req := types.ConsentRequest{
		ConsentChallenge: consentChallenge,
	}

	result, err := a.appAuth.ProcessConsent(c.Request.Context(), req)
	if err != nil {
		logger.Error("failed to process oauth2 consent",
			zap.Error(err),
			zap.String("consent_challenge", consentChallenge))

		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthHydraAcceptFailed, err.Error(), map[string]any{
			"consent_challenge": consentChallenge,
			"operation":         "consent",
		})

		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	logger.Info("consent accepted, redirecting")

	c.Redirect(http.StatusFound, result.RedirectTo)
}
