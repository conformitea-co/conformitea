// Package hydra provides a client for interacting with Ory Hydra's admin API.
package hydra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a Hydra admin API client for managing OAuth2 sessions and tokens.
type Client struct {
	adminURL   string
	httpClient *http.Client
}

// NewClient creates a new Hydra admin API client with the specified admin URL.
func NewClient(adminURL string) *Client {
	return &Client{
		adminURL: adminURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// LoginSession represents a Hydra OAuth2 login session request.
type LoginSession struct {
	Challenge      string   `json:"challenge"`
	ClientID       string   `json:"client"`
	RequestedAt    string   `json:"requested_at"`
	RequestURL     string   `json:"request_url"`
	Skip           bool     `json:"skip"`
	Subject        string   `json:"subject"`
	RequestedScope []string `json:"requested_scope"`
}

// TokenResponse represents Hydra's OAuth2 token response containing access and refresh tokens.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// AcceptLoginRequest represents the request payload for accepting a Hydra login session.
type AcceptLoginRequest struct {
	Subject     string `json:"subject"`
	Remember    bool   `json:"remember"`
	RememberFor int    `json:"remember_for"`
}

// GetLoginSession retrieves login session details from Hydra using the provided login challenge.
func (c *Client) GetLoginSession(loginChallenge string) (*LoginSession, error) {
	url := fmt.Sprintf("%s/admin/oauth2/auth/requests/login?login_challenge=%s", c.adminURL, loginChallenge)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get login session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hydra API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var session LoginSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("failed to decode login session: %w", err)
	}

	return &session, nil
}

// AcceptLoginSession accepts a Hydra login session with the provided user ID and returns OAuth2 tokens.
func (c *Client) AcceptLoginSession(loginChallenge, userID string) (*TokenResponse, error) {
	acceptReq := AcceptLoginRequest{
		Subject:     userID,
		Remember:    true,
		RememberFor: 3600,
	}

	jsonData, err := json.Marshal(acceptReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal accept request: %w", err)
	}

	url := fmt.Sprintf("%s/admin/oauth2/auth/requests/login/accept?login_challenge=%s", c.adminURL, loginChallenge)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create accept request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to accept login session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hydra accept API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// RejectLoginSession rejects a Hydra login session with the specified error code.
func (c *Client) RejectLoginSession(loginChallenge, errorCode string) error {
	rejectReq := map[string]interface{}{
		"error":             errorCode,
		"error_description": "Authentication failed",
	}

	jsonData, err := json.Marshal(rejectReq)
	if err != nil {
		return fmt.Errorf("failed to marshal reject request: %w", err)
	}

	url := fmt.Sprintf("%s/admin/oauth2/auth/requests/login/reject?login_challenge=%s", c.adminURL, loginChallenge)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create reject request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to reject login session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("hydra reject API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// TokenInfo represents the response from Hydra's token introspection endpoint.
type TokenInfo struct {
	Active bool   `json:"active"`
	Sub    string `json:"sub"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
	Scope  string `json:"scope"`
}

// IntrospectToken validates an OAuth2 token with Hydra's introspection endpoint.
func (c *Client) IntrospectToken(token string) (*TokenInfo, error) {
	data := fmt.Sprintf("token=%s", token)

	url := fmt.Sprintf("%s/admin/oauth2/introspect", c.adminURL)

	resp, err := c.httpClient.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(data))
	if err != nil {
		return nil, fmt.Errorf("failed to introspect token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hydra introspect API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenInfo TokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode token info: %w", err)
	}

	return &tokenInfo, nil
}
