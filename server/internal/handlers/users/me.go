package users

import (
	"net/http"

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
	// TODO: retrieve user details from db

	c.Status(http.StatusNoContent)
}
