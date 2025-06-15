package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cerror "github.com/conformitea-co/conformitea/internal/error"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TestMe_NotAuthenticated(t *testing.T) {
	router := setupTestRouter()
	router.GET("/auth/me", Me)

	req, _ := http.NewRequest("GET", "/auth/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response cerror.AuthError
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Code != cerror.AuthSessionExpired {
		t.Errorf("Expected error code %s, got %s", cerror.AuthSessionExpired, response.Code)
	}
}

func TestMe_MissingUserData(t *testing.T) {
	router := setupTestRouter()
	router.GET("/auth/me", Me)

	// Create a request with session that has authenticated=true but missing user data
	req, _ := http.NewRequest("GET", "/auth/me", nil)
	w := httptest.NewRecorder()

	// Add a route to set incomplete session data
	router.GET("/setup-session", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("authenticated", true)
		// Missing user_id and email
		session.Set("name", "Test User")
		session.Set("provider", "microsoft")
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

	// Make the /auth/me request with the session cookie
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response cerror.AuthError
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Code != cerror.AuthSessionExpired {
		t.Errorf("Expected error code %s, got %s", cerror.AuthSessionExpired, response.Code)
	}

	reason, ok := response.Details["reason"]
	if !ok || reason != "missing_user_data" {
		t.Error("Expected details to contain reason: missing_user_data")
	}
}

func TestMe_SuccessfulResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/auth/me", Me)

	// Create a request with a complete authenticated session
	req, _ := http.NewRequest("GET", "/auth/me", nil)
	w := httptest.NewRecorder()

	// Add a route to set complete session data
	router.GET("/setup-session", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("authenticated", true)
		session.Set("user_id", "user123")
		session.Set("email", "john.doe@example.com")
		session.Set("name", "John Doe")
		session.Set("provider", "microsoft")
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

	// Make the /auth/me request with the session cookie
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got '%s'", response.UserID)
	}

	if response.Email != "john.doe@example.com" {
		t.Errorf("Expected Email 'john.doe@example.com', got '%s'", response.Email)
	}

	if response.Name != "John Doe" {
		t.Errorf("Expected Name 'John Doe', got '%s'", response.Name)
	}

	if response.Provider != "microsoft" {
		t.Errorf("Expected Provider 'microsoft', got '%s'", response.Provider)
	}

	if !response.Authenticated {
		t.Error("Expected Authenticated to be true")
	}
}

func TestMe_PartialUserData(t *testing.T) {
	router := setupTestRouter()
	router.GET("/auth/me", Me)

	// Test with minimal required data (user_id and email)
	req, _ := http.NewRequest("GET", "/auth/me", nil)
	w := httptest.NewRecorder()

	// Add a route to set minimal session data
	router.GET("/setup-session", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("authenticated", true)
		session.Set("user_id", "user456")
		session.Set("email", "jane.doe@example.com")
		// Missing name and provider
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

	// Make the /auth/me request with the session cookie
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.UserID != "user456" {
		t.Errorf("Expected UserID 'user456', got '%s'", response.UserID)
	}

	if response.Email != "jane.doe@example.com" {
		t.Errorf("Expected Email 'jane.doe@example.com', got '%s'", response.Email)
	}

	// Name and Provider should be empty strings
	if response.Name != "" {
		t.Errorf("Expected Name to be empty, got '%s'", response.Name)
	}

	if response.Provider != "" {
		t.Errorf("Expected Provider to be empty, got '%s'", response.Provider)
	}

	if !response.Authenticated {
		t.Error("Expected Authenticated to be true")
	}
}

func TestMe_FalseAuthenticated(t *testing.T) {
	router := setupTestRouter()
	router.GET("/auth/me", Me)

	// Test with authenticated=false
	req, _ := http.NewRequest("GET", "/auth/me", nil)
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

	// Make the /auth/me request with the session cookie
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response cerror.AuthError
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Code != cerror.AuthSessionExpired {
		t.Errorf("Expected error code %s, got %s", cerror.AuthSessionExpired, response.Code)
	}
}
