// Package error provides ConformiTea-specific error types and codes.
package error

import (
	"encoding/json"
	"net/http"
)

// AuthError represents an authentication-related error with ConformiTea error codes.
type AuthError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *AuthError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}

// ConformiTea authentication error codes
const (
	AuthHydraClientInit       = "CT_AUTH_000"
	AuthInvalidState          = "CT_AUTH_001"
	AuthSessionNotFound       = "CT_AUTH_002"
	AuthProviderNotSupported  = "CT_AUTH_003"
	AuthMicrosoftExchange     = "CT_AUTH_004"
	AuthMicrosoftProfile      = "CT_AUTH_005"
	AuthHydraAcceptFailed     = "CT_AUTH_006"
	AuthSessionCreateFailed   = "CT_AUTH_007"
	AuthTokenIntrospectFailed = "CT_AUTH_008"
	AuthSessionExpired        = "CT_AUTH_009"
	AuthInvalidToken          = "CT_AUTH_010"
)

// NewAuthError creates a new AuthError with the specified code and optional details.
func NewAuthError(code string, details map[string]interface{}) *AuthError {
	return &AuthError{
		Code:    code,
		Details: details,
	}
}

// NewAuthErrorWithMessage creates a new AuthError with code, message, and optional details.
func NewAuthErrorWithMessage(code, message string, details map[string]interface{}) *AuthError {
	return &AuthError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// HTTPStatusCode returns the appropriate HTTP status code for the error.
func (e *AuthError) HTTPStatusCode() int {
	switch e.Code {
	case AuthInvalidState, AuthProviderNotSupported:
		return http.StatusBadRequest
	case AuthSessionNotFound, AuthSessionExpired, AuthInvalidToken:
		return http.StatusUnauthorized
	case AuthMicrosoftExchange, AuthMicrosoftProfile, AuthHydraAcceptFailed:
		return http.StatusBadGateway
	case AuthSessionCreateFailed, AuthTokenIntrospectFailed:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ToJSON serializes the error to JSON format.
func (e *AuthError) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
