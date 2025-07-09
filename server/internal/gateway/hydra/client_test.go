package hydra

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"conformitea/server/config"
)

// resetSingleton resets the singleton for testing purposes
func resetSingleton() {
	client = nil
	once = sync.Once{}
	initErr = nil
}

func TestInitialize(t *testing.T) {
	cfg := config.HydraConfig{
		AdminURL:  "http://localhost:4445",
		PublicURL: "http://localhost:4444",
	}

	// Reset the singleton for testing
	resetSingleton()

	err := Initialize(cfg)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	hydraClient, err := GetHydraClient()
	if err != nil {
		t.Fatalf("GetHydraClient() error = %v", err)
	}

	if hydraClient.adminURL != cfg.AdminURL {
		t.Errorf("Initialize() adminURL = %v, want %v", hydraClient.adminURL, cfg.AdminURL)
	}

	if hydraClient.httpClient == nil {
		t.Error("Initialize() httpClient is nil")
	}

	if hydraClient.httpClient.Timeout != 30*time.Second {
		t.Errorf("Initialize() httpClient.Timeout = %v, want %v", hydraClient.httpClient.Timeout, 30*time.Second)
	}
}

func TestInitialize_SingletonPattern(t *testing.T) {
	cfg := config.HydraConfig{
		AdminURL:  "http://localhost:4445",
		PublicURL: "http://localhost:4444",
	}

	// Reset the singleton for testing
	resetSingleton()

	// Initialize multiple times
	err1 := Initialize(cfg)
	if err1 != nil {
		t.Fatalf("First Initialize() error = %v", err1)
	}

	err2 := Initialize(cfg)
	if err2 != nil {
		t.Fatalf("Second Initialize() error = %v", err2)
	}

	err3 := Initialize(cfg)
	if err3 != nil {
		t.Fatalf("Third Initialize() error = %v", err3)
	}

	// Should only initialize once - get the same instance
	hydraClient1, err1 := GetHydraClient()
	if err1 != nil {
		t.Fatalf("GetHydraClient() error = %v", err1)
	}

	hydraClient2, err2 := GetHydraClient()
	if err2 != nil {
		t.Fatalf("GetHydraClient() error = %v", err2)
	}

	// Should be the same instance
	if hydraClient1 != hydraClient2 {
		t.Error("Multiple Initialize() calls should return the same client instance")
	}
}

func TestInitialize_EmptyAdminURL(t *testing.T) {
	cfg := config.HydraConfig{
		AdminURL:  "", // Empty admin URL should cause error
		PublicURL: "http://localhost:4444",
	}

	// Reset the singleton for testing
	resetSingleton()

	err := Initialize(cfg)
	if err == nil {
		t.Error("Initialize() should return error for empty AdminURL")
	}

	expectedError := "hydra admin URL is not configured"
	if err.Error() != expectedError {
		t.Errorf("Initialize() error = %v, want %v", err.Error(), expectedError)
	}

	// GetHydraClient should also return error
	_, err = GetHydraClient()
	if err == nil {
		t.Error("GetHydraClient() should return error when initialization failed")
	}
}

func TestGetHydraClient_NotInitialized(t *testing.T) {
	// Reset the singleton for testing
	resetSingleton()

	_, err := GetHydraClient()
	if err == nil {
		t.Error("GetHydraClient() should return error when not initialized")
	}

	expectedError := "hydra client was not initialized"
	if err.Error() != expectedError {
		t.Errorf("GetHydraClient() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestClient_GetLoginSession(t *testing.T) {
	tests := []struct {
		name            string
		challenge       string
		serverResponse  string
		statusCode      int
		expectError     bool
		expectedSession *HydraLoginSession
	}{
		{
			name:      "successful login session retrieval",
			challenge: "test-challenge-123",
			serverResponse: `{
				"challenge": "test-challenge-123",
				"client": { "client_id": "microsoft" },
				"requested_at": "2025-01-01T00:00:00Z",
				"request_url": "http://localhost:4444/oauth2/auth",
				"skip": false,
				"subject": "user123",
				"requested_scope": ["openid", "profile", "email"]
			}`,
			statusCode:  200,
			expectError: false,
			expectedSession: &HydraLoginSession{
				Challenge: "test-challenge-123",
				Client: struct {
					ClientId string `json:"client_id"`
				}{
					ClientId: "microsoft",
				},
				RequestURL:     "http://localhost:4444/oauth2/auth",
				Skip:           false,
				Subject:        "user123",
				RequestedScope: []string{"openid", "profile", "email"},
			},
		},
		{
			name:           "server error",
			challenge:      "invalid-challenge",
			serverResponse: `{"error": "not found"}`,
			statusCode:     404,
			expectError:    true,
		},
		{
			name:           "invalid JSON response",
			challenge:      "test-challenge",
			serverResponse: `invalid json`,
			statusCode:     200,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/admin/oauth2/auth/requests/login"
				if !strings.HasPrefix(r.URL.Path, expectedPath) {
					t.Errorf("Expected path to start with %s, got %s", expectedPath, r.URL.Path)
				}

				challenge := r.URL.Query().Get("login_challenge")
				if challenge != tt.challenge {
					t.Errorf("Expected challenge %s, got %s", tt.challenge, challenge)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Reset and initialize with test server
			resetSingleton()
			cfg := config.HydraConfig{
				AdminURL:  server.URL,
				PublicURL: "http://localhost:4444",
			}

			err := Initialize(cfg)
			if err != nil {
				t.Fatalf("Initialize() error = %v", err)
			}

			hydraClient, err := GetHydraClient()
			if err != nil {
				t.Fatalf("GetHydraClient() error = %v", err)
			}

			session, err := hydraClient.GetLoginSession(tt.challenge)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if session.Challenge != tt.expectedSession.Challenge {
				t.Errorf("Challenge = %v, want %v", session.Challenge, tt.expectedSession.Challenge)
			}

			if session.Client.ClientId != tt.expectedSession.Client.ClientId {
				t.Errorf("ClientID = %v, want %v", session.Client.ClientId, tt.expectedSession.Client.ClientId)
			}

			if len(session.RequestedScope) != len(tt.expectedSession.RequestedScope) {
				t.Errorf("RequestedScope length = %v, want %v", len(session.RequestedScope), len(tt.expectedSession.RequestedScope))
			}
		})
	}
}

func TestClient_AcceptLoginSession(t *testing.T) {
	tests := []struct {
		name           string
		challenge      string
		userID         string
		serverResponse string
		statusCode     int
		expectError    bool
		expectedTokens *TokenResponse
	}{
		{
			name:      "successful login acceptance",
			challenge: "test-challenge-123",
			userID:    "user123",
			serverResponse: `{
				"access_token": "access-token-123",
				"refresh_token": "refresh-token-123",
				"token_type": "Bearer",
				"expires_in": 3600
			}`,
			statusCode:  200,
			expectError: false,
			expectedTokens: &TokenResponse{
				AccessToken:  "access-token-123",
				RefreshToken: "refresh-token-123",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			},
		},
		{
			name:           "server error",
			challenge:      "invalid-challenge",
			userID:         "user123",
			serverResponse: `{"error": "invalid_challenge"}`,
			statusCode:     400,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("Expected PUT method, got %s", r.Method)
				}

				expectedPath := "/admin/oauth2/auth/requests/login/accept"
				if !strings.HasPrefix(r.URL.Path, expectedPath) {
					t.Errorf("Expected path to start with %s, got %s", expectedPath, r.URL.Path)
				}

				challenge := r.URL.Query().Get("login_challenge")
				if challenge != tt.challenge {
					t.Errorf("Expected challenge %s, got %s", tt.challenge, challenge)
				}

				// Verify request body
				var acceptReq AcceptLoginRequest
				if err := json.NewDecoder(r.Body).Decode(&acceptReq); err != nil {
					t.Errorf("Failed to decode request body: %v", err)
				}

				if acceptReq.Subject != tt.userID {
					t.Errorf("Expected subject %s, got %s", tt.userID, acceptReq.Subject)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Reset and initialize with test server
			resetSingleton()
			cfg := config.HydraConfig{
				AdminURL:  server.URL,
				PublicURL: "http://localhost:4444",
			}

			err := Initialize(cfg)
			if err != nil {
				t.Fatalf("Initialize() error = %v", err)
			}

			hydraClient, err := GetHydraClient()
			if err != nil {
				t.Fatalf("GetHydraClient() error = %v", err)
			}

			tokens, err := hydraClient.AcceptLoginSession(tt.challenge, tt.userID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tokens.AccessToken != tt.expectedTokens.AccessToken {
				t.Errorf("AccessToken = %v, want %v", tokens.AccessToken, tt.expectedTokens.AccessToken)
			}

			if tokens.TokenType != tt.expectedTokens.TokenType {
				t.Errorf("TokenType = %v, want %v", tokens.TokenType, tt.expectedTokens.TokenType)
			}

			if tokens.ExpiresIn != tt.expectedTokens.ExpiresIn {
				t.Errorf("ExpiresIn = %v, want %v", tokens.ExpiresIn, tt.expectedTokens.ExpiresIn)
			}
		})
	}
}

func TestClient_RejectLoginSession(t *testing.T) {
	tests := []struct {
		name        string
		challenge   string
		errorCode   string
		statusCode  int
		expectError bool
	}{
		{
			name:        "successful login rejection",
			challenge:   "test-challenge-123",
			errorCode:   "access_denied",
			statusCode:  200,
			expectError: false,
		},
		{
			name:        "server error",
			challenge:   "invalid-challenge",
			errorCode:   "access_denied",
			statusCode:  400,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("Expected PUT method, got %s", r.Method)
				}

				expectedPath := "/admin/oauth2/auth/requests/login/reject"
				if !strings.HasPrefix(r.URL.Path, expectedPath) {
					t.Errorf("Expected path to start with %s, got %s", expectedPath, r.URL.Path)
				}

				challenge := r.URL.Query().Get("login_challenge")
				if challenge != tt.challenge {
					t.Errorf("Expected challenge %s, got %s", tt.challenge, challenge)
				}

				// Verify request body
				var rejectReq map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&rejectReq); err != nil {
					t.Errorf("Failed to decode request body: %v", err)
				}

				if rejectReq["error"] != tt.errorCode {
					t.Errorf("Expected error %s, got %s", tt.errorCode, rejectReq["error"])
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(`{}`))
			}))
			defer server.Close()

			// Reset and initialize with test server
			resetSingleton()
			cfg := config.HydraConfig{
				AdminURL:  server.URL,
				PublicURL: "http://localhost:4444",
			}

			err := Initialize(cfg)
			if err != nil {
				t.Fatalf("Initialize() error = %v", err)
			}

			hydraClient, err := GetHydraClient()
			if err != nil {
				t.Fatalf("GetHydraClient() error = %v", err)
			}

			err = hydraClient.RejectLoginSession(tt.challenge, tt.errorCode)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestClient_IntrospectToken(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		serverResponse string
		statusCode     int
		expectError    bool
		expectedInfo   *TokenInfo
	}{
		{
			name:  "valid active token",
			token: "valid-token-123",
			serverResponse: `{
				"active": true,
				"sub": "user123",
				"exp": 1640995200,
				"iat": 1640991600,
				"scope": "openid profile email"
			}`,
			statusCode:  200,
			expectError: false,
			expectedInfo: &TokenInfo{
				Active: true,
				Sub:    "user123",
				Exp:    1640995200,
				Iat:    1640991600,
				Scope:  "openid profile email",
			},
		},
		{
			name:  "inactive token",
			token: "invalid-token",
			serverResponse: `{
				"active": false
			}`,
			statusCode:  200,
			expectError: false,
			expectedInfo: &TokenInfo{
				Active: false,
			},
		},
		{
			name:           "server error",
			token:          "test-token",
			serverResponse: `{"error": "server_error"}`,
			statusCode:     500,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST method, got %s", r.Method)
				}

				expectedPath := "/admin/oauth2/introspect"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				contentType := r.Header.Get("Content-Type")
				if contentType != "application/x-www-form-urlencoded" {
					t.Errorf("Expected Content-Type application/x-www-form-urlencoded, got %s", contentType)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Reset and initialize with test server
			resetSingleton()
			cfg := config.HydraConfig{
				AdminURL:  server.URL,
				PublicURL: "http://localhost:4444",
			}

			err := Initialize(cfg)
			if err != nil {
				t.Fatalf("Initialize() error = %v", err)
			}

			hydraClient, err := GetHydraClient()
			if err != nil {
				t.Fatalf("GetHydraClient() error = %v", err)
			}

			tokenInfo, err := hydraClient.IntrospectToken(tt.token)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tokenInfo.Active != tt.expectedInfo.Active {
				t.Errorf("Active = %v, want %v", tokenInfo.Active, tt.expectedInfo.Active)
			}

			if tokenInfo.Sub != tt.expectedInfo.Sub {
				t.Errorf("Sub = %v, want %v", tokenInfo.Sub, tt.expectedInfo.Sub)
			}

			if tokenInfo.Scope != tt.expectedInfo.Scope {
				t.Errorf("Scope = %v, want %v", tokenInfo.Scope, tt.expectedInfo.Scope)
			}
		})
	}
}
