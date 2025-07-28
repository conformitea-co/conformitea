package users

import (
	"net/http"

	"conformitea/server/internal/cerror"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// MeResponse represents the authenticated user data returned by the me endpoint.
type MeResponse struct {
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Provider      string `json:"provider"`
	Authenticated bool   `json:"authenticated"`
}

// Me returns the current user's session information.
func (a *UsersHandlers) Me(c *gin.Context) {
	session := sessions.Default(c)

	// Check if user is authenticated
	authenticated, exists := session.Get("authenticated").(bool)
	if !exists || !authenticated {
		authErr := cerror.NewAuthError(cerror.AuthSessionExpired, nil)
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	// Extract user data from session
	userID, _ := session.Get("user_id").(string)
	email, _ := session.Get("email").(string)
	name, _ := session.Get("name").(string)
	provider, _ := session.Get("provider").(string)

	// Validate required fields
	if userID == "" || email == "" {
		authErr := cerror.NewAuthError(cerror.AuthSessionExpired, map[string]any{
			"reason": "missing_user_data",
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	user := MeResponse{
		UserID:        userID,
		Email:         email,
		Name:          name,
		Provider:      provider,
		Authenticated: true,
	}

	c.JSON(http.StatusOK, user)
}
