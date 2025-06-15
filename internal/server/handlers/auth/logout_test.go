package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TestLogout_NotAuthenticated(t *testing.T) {
	router := setupTestRouter()
	router.POST("/auth/logout", Logout)

	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	message, ok := response["message"]
	if !ok || message != "User already logged out" {
		t.Error("Expected message: User already logged out")
	}

	authenticated, ok := response["authenticated"]
	if !ok || authenticated != false {
		t.Error("Expected authenticated to be false")
	}
}

func TestLogout_SuccessfulLogout(t *testing.T) {
	router := setupTestRouter()
	router.POST("/auth/logout", Logout)

	// Create a request with an authenticated session
	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	w := httptest.NewRecorder()

	// Add a route to set up an authenticated session
	router.GET("/setup-session", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("authenticated", true)
		session.Set("user_id", "user123")
		session.Set("email", "john.doe@example.com")
		session.Set("name", "John Doe")
		session.Set("provider", "microsoft")
		session.Set("access_token", "access-token-123")
		session.Set("refresh_token", "refresh-token-123")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"status": "session_set"})
	})

	// Add a route to check if session is cleared
	router.GET("/check-session", func(c *gin.Context) {
		session := sessions.Default(c)
		authenticated := session.Get("authenticated")
		userID := session.Get("user_id")
		email := session.Get("email")

		c.JSON(http.StatusOK, gin.H{
			"authenticated": authenticated,
			"user_id":       userID,
			"email":         email,
		})
	})

	// First, set up the authenticated session
	setupReq, _ := http.NewRequest("GET", "/setup-session", nil)
	setupW := httptest.NewRecorder()
	router.ServeHTTP(setupW, setupReq)

	// Extract session cookie
	cookies := setupW.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("Expected session cookie to be set")
	}

	// Make the logout request with the session cookie
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	message, ok := response["message"]
	if !ok || message != "Successfully logged out" {
		t.Error("Expected message: Successfully logged out")
	}

	authenticated, ok := response["authenticated"]
	if !ok || authenticated != false {
		t.Error("Expected authenticated to be false")
	}

	// Verify session is cleared by making a check request
	logoutCookies := w.Result().Cookies()
	checkReq, _ := http.NewRequest("GET", "/check-session", nil)
	for _, cookie := range logoutCookies {
		checkReq.AddCookie(cookie)
	}

	checkW := httptest.NewRecorder()
	router.ServeHTTP(checkW, checkReq)

	if checkW.Code != http.StatusOK {
		t.Errorf("Expected status %d for session check, got %d", http.StatusOK, checkW.Code)
	}

	var sessionData map[string]interface{}
	if err := json.Unmarshal(checkW.Body.Bytes(), &sessionData); err != nil {
		t.Errorf("Failed to unmarshal session data: %v", err)
	}

	// All session values should be nil after logout
	if sessionData["authenticated"] != nil {
		t.Error("Expected authenticated to be nil after logout")
	}

	if sessionData["user_id"] != nil {
		t.Error("Expected user_id to be nil after logout")
	}

	if sessionData["email"] != nil {
		t.Error("Expected email to be nil after logout")
	}
}

func TestLogout_FalseAuthenticated(t *testing.T) {
	router := setupTestRouter()
	router.POST("/auth/logout", Logout)

	// Create a request with authenticated=false
	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	w := httptest.NewRecorder()

	// Add a route to set session with authenticated=false
	router.GET("/setup-session", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("authenticated", false)
		session.Set("user_id", "user123")
		session.Set("email", "john.doe@example.com")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"status": "session_set"})
	})

	// First, set up the session
	setupReq, _ := http.NewRequest("GET", "/setup-session", nil)
	setupW := httptest.NewRecorder()
	router.ServeHTTP(setupW, setupReq)

	// Extract session cookie
	cookies := setupW.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("Expected session cookie to be set")
	}

	// Make the logout request with the session cookie
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	message, ok := response["message"]
	if !ok || message != "User already logged out" {
		t.Error("Expected message: User already logged out")
	}

	authenticated, ok := response["authenticated"]
	if !ok || authenticated != false {
		t.Error("Expected authenticated to be false")
	}
}
