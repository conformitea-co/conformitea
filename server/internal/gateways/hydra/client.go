// Package hydra provides a client for interacting with Ory Hydra's admin API.
package hydra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"conformitea/server/internal/config"
)

var client *HydraClient

func Initialize() {
	client = &HydraClient{
		adminURL: config.GetConfig().Hydra.AdminURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func GetHydraClient() (*HydraClient, error) {
	if client == nil {
		return nil, fmt.Errorf("hydra client was not initialized")
	}

	return client, nil
}

// GetLoginSession retrieves login session details from Hydra using the provided login challenge.
func (c *HydraClient) GetLoginSession(loginChallenge string) (*HydraLoginSession, error) {
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

	var session HydraLoginSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("failed to decode login session: %w", err)
	}

	return &session, nil
}

// AcceptLoginSession accepts a Hydra login session with the provided user ID and returns OAuth2 tokens.
func (c *HydraClient) AcceptLoginSession(loginChallenge, userID string) (*TokenResponse, error) {
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
func (c *HydraClient) RejectLoginSession(loginChallenge, errorCode string) error {
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

// IntrospectToken validates an OAuth2 token with Hydra's introspection endpoint.
func (c *HydraClient) IntrospectToken(token string) (*TokenInfo, error) {
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
