package auth

import (
	"net/http"

	"conformitea/server/internal/cerror"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Logout clears the user's session and logs them out.
func (a *AuthHandlers) Logout(c *gin.Context) {
	session := sessions.Default(c)

	// Check if user is authenticated
	authenticated, exists := session.Get("authenticated").(bool)
	if !exists || !authenticated {
		// User is already logged out
		c.JSON(http.StatusOK, gin.H{
			"message":       "User already logged out",
			"authenticated": false,
		})
		return
	}

	// Clear all session data
	session.Clear()
	if err := session.Save(); err != nil {
		authErr := cerror.NewAuthErrorWithMessage(cerror.AuthSessionCreateFailed, err.Error(), map[string]interface{}{
			"operation": "logout",
		})
		c.JSON(authErr.HTTPStatusCode(), authErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Successfully logged out",
		"authenticated": false,
	})
}
