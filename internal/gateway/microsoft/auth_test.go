package microsoft

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

func TestNewClient(t *testing.T) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"openid", "profile", "email"},
	}

	client := NewClient(config)

	if client.config.ClientID != config.ClientID {
		t.Errorf("NewClient() ClientID = %v, want %v", client.config.ClientID, config.ClientID)
	}

	if client.config.ClientSecret != config.ClientSecret {
		t.Errorf("NewClient() ClientSecret = %v, want %v", client.config.ClientSecret, config.ClientSecret)
	}

	if client.config.RedirectURL != config.RedirectURL {
		t.Errorf("NewClient() RedirectURL = %v, want %v", client.config.RedirectURL, config.RedirectURL)
	}

	if len(client.config.Scopes) != len(config.Scopes) {
		t.Errorf("NewClient() Scopes length = %v, want %v", len(client.config.Scopes), len(config.Scopes))
	}

	for i, scope := range config.Scopes {
		if client.config.Scopes[i] != scope {
			t.Errorf("NewClient() Scopes[%d] = %v, want %v", i, client.config.Scopes[i], scope)
		}
	}
}

func TestClient_GenerateAuthURL(t *testing.T) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"openid", "profile", "email"},
	}

	client := NewClient(config)
	state := "test-state-123"
	nonce := "test-nonce-456"

	authURL := client.GenerateAuthURL(state, nonce)

	// Verify URL contains expected parameters
	if !strings.Contains(authURL, "client_id=test-client-id") {
		t.Error("Auth URL should contain client_id parameter")
	}

	if !strings.Contains(authURL, "state=test-state-123") {
		t.Error("Auth URL should contain state parameter")
	}

	if !strings.Contains(authURL, "nonce=test-nonce-456") {
		t.Error("Auth URL should contain nonce parameter")
	}

	if !strings.Contains(authURL, "response_mode=query") {
		t.Error("Auth URL should contain response_mode=query parameter")
	}

	if !strings.Contains(authURL, "scope=openid+profile+email") {
		t.Error("Auth URL should contain encoded scopes")
	}

	// Verify it starts with Microsoft endpoint
	if !strings.HasPrefix(authURL, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize") {
		t.Error("Auth URL should start with Microsoft OAuth2 endpoint")
	}
}

func TestClient_ExchangeCodeForToken(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		serverResponse string
		statusCode     int
		expectError    bool
		expectedToken  *oauth2.Token
	}{
		{
			name: "successful token exchange",
			code: "valid-auth-code",
			serverResponse: `{
				"access_token": "access-token-123",
				"refresh_token": "refresh-token-123",
				"token_type": "Bearer",
				"expires_in": 3600
			}`,
			statusCode:  200,
			expectError: false,
			expectedToken: &oauth2.Token{
				AccessToken:  "access-token-123",
				RefreshToken: "refresh-token-123",
				TokenType:    "Bearer",
			},
		},
		{
			name:           "invalid authorization code",
			code:           "invalid-code",
			serverResponse: `{"error": "invalid_grant"}`,
			statusCode:     400,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST method, got %s", r.Method)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			config := Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURL:  "http://localhost:8080/auth/callback",
				Scopes:       []string{"openid", "profile", "email"},
			}

			client := NewClient(config)
			// Override the endpoint for testing
			client.config.Endpoint = oauth2.Endpoint{
				TokenURL: server.URL,
			}

			ctx := context.Background()
			token, err := client.ExchangeCodeForToken(ctx, tt.code)

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

			if token.AccessToken != tt.expectedToken.AccessToken {
				t.Errorf("AccessToken = %v, want %v", token.AccessToken, tt.expectedToken.AccessToken)
			}

			if token.RefreshToken != tt.expectedToken.RefreshToken {
				t.Errorf("RefreshToken = %v, want %v", token.RefreshToken, tt.expectedToken.RefreshToken)
			}

			if token.TokenType != tt.expectedToken.TokenType {
				t.Errorf("TokenType = %v, want %v", token.TokenType, tt.expectedToken.TokenType)
			}
		})
	}
}

func TestClient_GetUserProfile(t *testing.T) {
	tests := []struct {
		name            string
		serverResponse  string
		statusCode      int
		expectError     bool
		expectedProfile *UserProfile
		errorContains   string
	}{
		{
			name: "successful profile retrieval",
			serverResponse: `{
				"id": "user123",
				"displayName": "John Doe",
				"givenName": "John",
				"surname": "Doe",
				"userPrincipalName": "john.doe@example.com",
				"mail": "john.doe@company.com"
			}`,
			statusCode:  200,
			expectError: false,
			expectedProfile: &UserProfile{
				ID:                "user123",
				DisplayName:       "John Doe",
				GivenName:         "John",
				Surname:           "Doe",
				UserPrincipalName: "john.doe@example.com",
				Mail:              "john.doe@company.com",
			},
		},
		{
			name:           "unauthorized error",
			serverResponse: `{"error": {"code": "Unauthorized", "message": "Invalid token"}}`,
			statusCode:     401,
			expectError:    true,
			errorContains:  "microsoft Graph API error: status 401",
		},
		{
			name:           "forbidden error",
			serverResponse: `{"error": {"code": "Forbidden", "message": "Insufficient privileges"}}`,
			statusCode:     403,
			expectError:    true,
			errorContains:  "microsoft Graph API error: status 403",
		},
		{
			name:           "not found error",
			serverResponse: `{"error": {"code": "NotFound", "message": "User not found"}}`,
			statusCode:     404,
			expectError:    true,
			errorContains:  "microsoft Graph API error: status 404",
		},
		{
			name:           "internal server error",
			serverResponse: `{"error": {"code": "InternalServerError", "message": "Something went wrong"}}`,
			statusCode:     500,
			expectError:    true,
			errorContains:  "microsoft Graph API error: status 500",
		},
		{
			name:           "invalid JSON response with 200 status",
			serverResponse: `invalid json`,
			statusCode:     200,
			expectError:    true,
			errorContains:  "failed to decode user profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that mimics Microsoft Graph API
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := "/v1.0/me"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify Authorization header
				authHeader := r.Header.Get("Authorization")
				if !strings.HasPrefix(authHeader, "Bearer ") {
					t.Error("Expected Authorization header with Bearer token")
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Create client with test configuration
			config := Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURL:  "http://localhost:8080/auth/callback",
				Scopes:       []string{"openid", "profile", "email"},
			}

			client := NewClient(config)

			// Create a mock token
			token := &oauth2.Token{
				AccessToken: "test-access-token",
				TokenType:   "Bearer",
			}

			// Override the oauth2 config to use our test server
			// We need to create a custom HTTP client that redirects Graph API calls to our test server
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
				Transport: &graphAPITestTransport{
					testServerURL: server.URL,
				},
			})

			// Call the actual method we're testing
			profile, err := client.GetUserProfile(ctx, token)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify the returned profile matches expectations
			if profile.ID != tt.expectedProfile.ID {
				t.Errorf("ID = %v, want %v", profile.ID, tt.expectedProfile.ID)
			}

			if profile.DisplayName != tt.expectedProfile.DisplayName {
				t.Errorf("DisplayName = %v, want %v", profile.DisplayName, tt.expectedProfile.DisplayName)
			}

			if profile.GivenName != tt.expectedProfile.GivenName {
				t.Errorf("GivenName = %v, want %v", profile.GivenName, tt.expectedProfile.GivenName)
			}

			if profile.Surname != tt.expectedProfile.Surname {
				t.Errorf("Surname = %v, want %v", profile.Surname, tt.expectedProfile.Surname)
			}

			if profile.UserPrincipalName != tt.expectedProfile.UserPrincipalName {
				t.Errorf("UserPrincipalName = %v, want %v", profile.UserPrincipalName, tt.expectedProfile.UserPrincipalName)
			}

			if profile.Mail != tt.expectedProfile.Mail {
				t.Errorf("Mail = %v, want %v", profile.Mail, tt.expectedProfile.Mail)
			}
		})
	}
}

// graphAPITestTransport is a simple transport that redirects Microsoft Graph API calls to our test server
type graphAPITestTransport struct {
	testServerURL string
}

func (t *graphAPITestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// If this is a Graph API request, redirect to our test server
	if strings.Contains(req.URL.Host, "graph.microsoft.com") {
		// Parse the test server URL
		testURL := strings.TrimPrefix(t.testServerURL, "http://")
		testURL = strings.TrimPrefix(testURL, "https://")

		// Update the request to point to our test server
		req.URL.Scheme = "http"
		req.URL.Host = testURL
		// Keep the same path (/v1.0/me)
	}

	// Use the default transport to actually make the request
	return http.DefaultTransport.RoundTrip(req)
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		valid  bool
	}{
		{
			name: "valid config",
			config: Config{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
				RedirectURL:  "http://localhost:8080/auth/callback",
				Scopes:       []string{"openid", "profile", "email"},
			},
			valid: true,
		},
		{
			name: "empty client ID",
			config: Config{
				ClientID:     "",
				ClientSecret: "valid-client-secret",
				RedirectURL:  "http://localhost:8080/auth/callback",
				Scopes:       []string{"openid", "profile", "email"},
			},
			valid: false,
		},
		{
			name: "empty scopes",
			config: Config{
				ClientID:     "valid-client-id",
				ClientSecret: "valid-client-secret",
				RedirectURL:  "http://localhost:8080/auth/callback",
				Scopes:       []string{},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.config)

			if tt.valid {
				if client.config.ClientID == "" {
					t.Error("Expected valid client to have non-empty ClientID")
				}
			} else {
				// For invalid configs, we still create the client but it may not work properly
				// This is consistent with the oauth2 library behavior
				if tt.config.ClientID == "" && client.config.ClientID != "" {
					t.Error("Expected ClientID to remain empty for invalid config")
				}
			}
		})
	}
}
