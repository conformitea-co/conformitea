package error

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAuthError_Error(t *testing.T) {
	tests := []struct {
		name     string
		authErr  *AuthError
		expected string
	}{
		{
			name: "error with message",
			authErr: &AuthError{
				Code:    AuthInvalidState,
				Message: "Custom error message",
			},
			expected: "Custom error message",
		},
		{
			name: "error without message",
			authErr: &AuthError{
				Code: AuthSessionNotFound,
			},
			expected: AuthSessionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.authErr.Error(); got != tt.expected {
				t.Errorf("AuthError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewAuthError(t *testing.T) {
	details := map[string]interface{}{
		"parameter": "test_param",
		"reason":    "missing",
	}

	authErr := NewAuthError(AuthInvalidState, details)

	if authErr.Code != AuthInvalidState {
		t.Errorf("NewAuthError() code = %v, want %v", authErr.Code, AuthInvalidState)
	}

	if authErr.Details["parameter"] != "test_param" {
		t.Errorf("NewAuthError() details[parameter] = %v, want %v", authErr.Details["parameter"], "test_param")
	}
}

func TestNewAuthErrorWithMessage(t *testing.T) {
	message := "Test error message"
	details := map[string]interface{}{
		"context": "test",
	}

	authErr := NewAuthErrorWithMessage(AuthMicrosoftExchange, message, details)

	if authErr.Code != AuthMicrosoftExchange {
		t.Errorf("NewAuthErrorWithMessage() code = %v, want %v", authErr.Code, AuthMicrosoftExchange)
	}

	if authErr.Message != message {
		t.Errorf("NewAuthErrorWithMessage() message = %v, want %v", authErr.Message, message)
	}

	if authErr.Details["context"] != "test" {
		t.Errorf("NewAuthErrorWithMessage() details[context] = %v, want %v", authErr.Details["context"], "test")
	}
}

func TestAuthError_HTTPStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "bad request error",
			code:     AuthInvalidState,
			expected: http.StatusBadRequest,
		},
		{
			name:     "unauthorized error",
			code:     AuthSessionNotFound,
			expected: http.StatusUnauthorized,
		},
		{
			name:     "bad gateway error",
			code:     AuthMicrosoftExchange,
			expected: http.StatusBadGateway,
		},
		{
			name:     "internal server error",
			code:     AuthSessionCreateFailed,
			expected: http.StatusInternalServerError,
		},
		{
			name:     "unknown error defaults to internal server error",
			code:     "UNKNOWN_ERROR",
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authErr := &AuthError{Code: tt.code}
			if got := authErr.HTTPStatusCode(); got != tt.expected {
				t.Errorf("AuthError.HTTPStatusCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAuthError_ToJSON(t *testing.T) {
	authErr := &AuthError{
		Code:    AuthInvalidToken,
		Message: "Token validation failed",
		Details: map[string]interface{}{
			"token_type": "access_token",
		},
	}

	jsonBytes, err := authErr.ToJSON()
	if err != nil {
		t.Fatalf("AuthError.ToJSON() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["code"] != AuthInvalidToken {
		t.Errorf("JSON code = %v, want %v", result["code"], AuthInvalidToken)
	}

	if result["message"] != "Token validation failed" {
		t.Errorf("JSON message = %v, want %v", result["message"], "Token validation failed")
	}

	details, ok := result["details"].(map[string]interface{})
	if !ok {
		t.Fatal("JSON details is not a map")
	}

	if details["token_type"] != "access_token" {
		t.Errorf("JSON details[token_type] = %v, want %v", details["token_type"], "access_token")
	}
}

func TestErrorCodes(t *testing.T) {
	// Test that all error codes are properly defined
	expectedCodes := map[string]string{
		"AuthInvalidState":          "CT_AUTH_001",
		"AuthSessionNotFound":       "CT_AUTH_002",
		"AuthProviderNotSupported":  "CT_AUTH_003",
		"AuthMicrosoftExchange":     "CT_AUTH_004",
		"AuthMicrosoftProfile":      "CT_AUTH_005",
		"AuthHydraAcceptFailed":     "CT_AUTH_006",
		"AuthSessionCreateFailed":   "CT_AUTH_007",
		"AuthTokenIntrospectFailed": "CT_AUTH_008",
		"AuthSessionExpired":        "CT_AUTH_009",
		"AuthInvalidToken":          "CT_AUTH_010",
	}

	actualCodes := map[string]string{
		"AuthInvalidState":          AuthInvalidState,
		"AuthSessionNotFound":       AuthSessionNotFound,
		"AuthProviderNotSupported":  AuthProviderNotSupported,
		"AuthMicrosoftExchange":     AuthMicrosoftExchange,
		"AuthMicrosoftProfile":      AuthMicrosoftProfile,
		"AuthHydraAcceptFailed":     AuthHydraAcceptFailed,
		"AuthSessionCreateFailed":   AuthSessionCreateFailed,
		"AuthTokenIntrospectFailed": AuthTokenIntrospectFailed,
		"AuthSessionExpired":        AuthSessionExpired,
		"AuthInvalidToken":          AuthInvalidToken,
	}

	for name, expected := range expectedCodes {
		if actual, exists := actualCodes[name]; !exists || actual != expected {
			t.Errorf("Error code %s = %v, want %v", name, actual, expected)
		}
	}
}
