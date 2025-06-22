package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cerror "github.com/conformitea-co/conformitea/internal/error"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup session middleware for testing
	store := cookie.NewStore([]byte("test-secret-key"))
	router.Use(sessions.Sessions("test_session", store))

	return router
}

func setupTestConfig() {
	viper.Set("hydra.admin_url", "http://localhost:4445")
	viper.Set("oauth.microsoft.client_id", "test-client-id")
	viper.Set("oauth.microsoft.client_secret", "test-client-secret")
	viper.Set("oauth.microsoft.scopes", []string{"openid", "profile", "email"})
}

func TestLogin_MissingChallenge(t *testing.T) {
	setupTestConfig()
	router := setupTestRouter()
	router.GET("/auth/login", Login)

	req, _ := http.NewRequest("GET", "/auth/login", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response cerror.AuthError
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Code != cerror.AuthInvalidState {
		t.Errorf("Expected error code %s, got %s", cerror.AuthInvalidState, response.Code)
	}

	details, ok := response.Details["parameter"]
	if !ok || details != "login_challenge" {
		t.Error("Expected details to contain parameter: login_challenge")
	}
}

func TestLogin_InvalidChallenge(t *testing.T) {
	setupTestConfig()
	router := setupTestRouter()
	router.GET("/auth/login", Login)

	// Mock Hydra server that returns error for invalid challenge
	hydraServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/admin/oauth2/auth/requests/login") {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "not_found"}`))
		}
	}))
	defer hydraServer.Close()

	viper.Set("hydra.admin_url", hydraServer.URL)

	req, _ := http.NewRequest("GET", "/auth/login?login_challenge=invalid-challenge", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response cerror.AuthError
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Code != cerror.AuthSessionNotFound {
		t.Errorf("Expected error code %s, got %s", cerror.AuthSessionNotFound, response.Code)
	}
}

func TestLogin_UnsupportedProvider(t *testing.T) {
	setupTestConfig()
	router := setupTestRouter()
	router.GET("/auth/login", Login)

	// Mock Hydra server that returns a session with unsupported client
	hydraServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/admin/oauth2/auth/requests/login") {
			response := `{
				"challenge": "test-challenge",
				"client": {
					"client_id": "google"
				},
				"requested_at": "2025-01-01T00:00:00Z",
				"request_url": "http://localhost:4444/oauth2/auth",
				"skip": false,
				"subject": "",
				"requested_scope": ["openid", "profile", "email"]
			}`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		}
	}))
	defer hydraServer.Close()

	viper.Set("hydra.admin_url", hydraServer.URL)

	req, _ := http.NewRequest("GET", "/auth/login?login_challenge=test-challenge", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response cerror.AuthError
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Code != cerror.AuthProviderNotSupported {
		t.Errorf("Expected error code %s, got %s", cerror.AuthProviderNotSupported, response.Code)
	}

	provider, ok := response.Details["provider"]
	if !ok || provider != "google" {
		t.Error("Expected details to contain provider: google")
	}
}

func TestLogin_SuccessfulMicrosoftRedirect(t *testing.T) {
	setupTestConfig()
	router := setupTestRouter()
	router.GET("/auth/login", Login)

	// Mock Hydra server that returns a valid Microsoft session
	hydraServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/admin/oauth2/auth/requests/login") {
			challenge := r.URL.Query().Get("login_challenge")
			if challenge != "microsoft-challenge" {
				t.Errorf("Expected challenge 'microsoft-challenge', got '%s'", challenge)
			}

			response := `{
				"challenge": "microsoft-challenge",
				"client": {
					"client_id": "microsoft"
				},
				"requested_at": "2025-01-01T00:00:00Z",
				"request_url": "http://localhost:4444/oauth2/auth",
				"skip": false,
				"subject": "",
				"requested_scope": ["openid", "profile", "email"]
			}`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		}
	}))
	defer hydraServer.Close()

	viper.Set("hydra.admin_url", hydraServer.URL)

	req, _ := http.NewRequest("GET", "/auth/login?login_challenge=microsoft-challenge", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status %d (redirect), got %d", http.StatusFound, w.Code)
	}

	location := w.Header().Get("Location")
	if !strings.Contains(location, "login.microsoftonline.com") {
		t.Error("Expected redirect to Microsoft OAuth2 endpoint")
	}

	if !strings.Contains(location, "client_id=test-client-id") {
		t.Error("Expected client_id in redirect URL")
	}

	if !strings.Contains(location, "state=microsoft-challenge") {
		t.Error("Expected state parameter in redirect URL")
	}

	if !strings.Contains(location, "scope=openid+profile+email") {
		t.Error("Expected scopes in redirect URL")
	}
}

func TestGenerateNonce(t *testing.T) {
	nonce1, err1 := generateNonce()
	if err1 != nil {
		t.Errorf("generateNonce() error = %v", err1)
	}

	nonce2, err2 := generateNonce()
	if err2 != nil {
		t.Errorf("generateNonce() error = %v", err2)
	}

	if len(nonce1) != 32 { // 16 bytes hex encoded = 32 characters
		t.Errorf("Expected nonce length 32, got %d", len(nonce1))
	}

	if nonce1 == nonce2 {
		t.Error("Expected different nonces on multiple calls")
	}

	// Verify it's valid hex
	for _, char := range nonce1 {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Errorf("Nonce contains invalid hex character: %c", char)
		}
	}
}

func TestLogin_SessionStorage(t *testing.T) {
	setupTestConfig()
	router := setupTestRouter()
	router.GET("/auth/login", Login)

	// Mock Hydra server
	hydraServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/admin/oauth2/auth/requests/login") {
			response := `{
				"challenge": "test-challenge",
				"client": {
					"client_id": "microsoft"
				},
				"requested_at": "2025-01-01T00:00:00Z",
				"request_url": "http://localhost:4444/oauth2/auth",
				"skip": false,
				"subject": "",
				"requested_scope": ["openid", "profile", "email"]
			}`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		}
	}))
	defer hydraServer.Close()

	viper.Set("hydra.admin_url", hydraServer.URL)

	// Create a request with session support
	req, _ := http.NewRequest("GET", "/auth/login?login_challenge=test-challenge", nil)
	w := httptest.NewRecorder()

	// Setup router with session middleware for this test
	testRouter := setupTestRouter()
	testRouter.GET("/auth/login", Login)

	// Add a route to check session values
	testRouter.GET("/check-session", func(c *gin.Context) {
		session := sessions.Default(c)
		challenge := session.Get("hydra_login_challenge")
		provider := session.Get("idp_provider")
		nonce := session.Get("auth_nonce")

		c.JSON(http.StatusOK, gin.H{
			"challenge": challenge,
			"provider":  provider,
			"nonce":     nonce,
		})
	})

	// First, make the login request
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status %d, got %d", http.StatusFound, w.Code)
	}

	// Extract session cookie
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("Expected session cookie to be set")
	}

	// Make a second request to check session values
	checkReq, _ := http.NewRequest("GET", "/check-session", nil)
	for _, cookie := range cookies {
		checkReq.AddCookie(cookie)
	}

	checkW := httptest.NewRecorder()
	testRouter.ServeHTTP(checkW, checkReq)

	if checkW.Code != http.StatusOK {
		t.Errorf("Expected status %d for session check, got %d", http.StatusOK, checkW.Code)
	}

	var sessionData map[string]interface{}
	if err := json.Unmarshal(checkW.Body.Bytes(), &sessionData); err != nil {
		t.Errorf("Failed to unmarshal session data: %v", err)
	}

	if sessionData["challenge"] != "test-challenge" {
		t.Errorf("Expected challenge 'test-challenge', got %v", sessionData["challenge"])
	}

	if sessionData["provider"] != "microsoft" {
		t.Errorf("Expected provider 'microsoft', got %v", sessionData["provider"])
	}

	if sessionData["nonce"] == nil || sessionData["nonce"] == "" {
		t.Error("Expected nonce to be set in session")
	}
}
