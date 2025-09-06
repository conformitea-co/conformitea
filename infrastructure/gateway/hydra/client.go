// Package hydra provides a client for interacting with Ory Hydra's admin API.
package hydra

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"conformitea/infrastructure/config"
)

var client *HydraClient

func Initialize(hydraConfigValues config.HydraConfig) (*HydraClient, error) {
	if err := hydraConfigValues.Validate(); err != nil {
		return nil, fmt.Errorf("invalid hydra configuration: %w", err)
	}

	// Create HTTP client
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Derive JWKS URL from the public URL (e.g., https://hydra.example/.well-known/jwks.json)
	jwksURL := strings.TrimRight(hydraConfigValues.PublicURL, "/") + "/.well-known/jwks.json"

	// Create JWKS-backed JWT parser
	jwtParser, err := NewJWTParserFromURL(jwksURL, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWKS parser: %w", err)
	}

	client = &HydraClient{
		adminURL:              hydraConfigValues.AdminURL,
		publicURL:             hydraConfigValues.PublicURL,
		redirectURL:           hydraConfigValues.RedirectURL,
		httpClient:            httpClient,
		jwtParser:             jwtParser,
		microsoftClientID:     hydraConfigValues.Microsoft.ClientID,
		microsoftClientSecret: hydraConfigValues.Microsoft.ClientSecret,
	}

	return client, nil
}

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
func (c *HydraClient) AcceptLoginSession(loginChallenge, userID string) (*AcceptLoginResponse, error) {
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

	var result AcceptLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode accept login response: %w", err)
	}

	return &result, nil
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

// ExchangeAuthorizationCode exchanges an OAuth2 authorization code for tokens
func (c *HydraClient) ExchangeAuthorizationCode(code string) (*TokenExchangeResponse, error) {
	// Token exchange should happen with the public endpoint
	tokenURL := fmt.Sprintf("%s/oauth2/token", c.publicURL)

	// Prepare form data for authorization code grant
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	// include client_id in the body for clarity; primary auth is via HTTP Basic
	data.Set("client_id", c.microsoftClientID)
	data.Set("redirect_uri", c.redirectURL)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add HTTP Basic auth header using client credentials
	creds := c.microsoftClientID + ":" + c.microsoftClientSecret
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(creds)))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hydra token exchange error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// ParseJWT parses and validates a JWT token using the loaded public key
func (c *HydraClient) ParseJWT(tokenString string) (map[string]interface{}, error) {
	if c.jwtParser == nil {
		return nil, fmt.Errorf("JWT parser not initialized")
	}

	claims, err := c.jwtParser.ParseAndValidateJWT(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	return claims, nil
}

// ParseJWTWithoutValidation parses a JWT without validating the signature
func (c *HydraClient) ParseJWTWithoutValidation(tokenString string) (map[string]interface{}, error) {
	if c.jwtParser == nil {
		return nil, fmt.Errorf("JWT parser not initialized")
	}

	claims, err := c.jwtParser.ParseJWTWithoutValidation(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	return claims, nil
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

// GetConsentSession retrieves consent session details from Hydra using the provided consent challenge.
func (c *HydraClient) GetConsentSession(consentChallenge string) (*HydraGetConsentResponse, error) {
	url := fmt.Sprintf("%s/admin/oauth2/auth/requests/consent?consent_challenge=%s", c.adminURL, consentChallenge)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get consent session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hydra API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var session HydraGetConsentResponse
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("failed to decode consent session: %w", err)
	}

	return &session, nil
}

// AcceptConsentSession accepts a Hydra consent session with the provided claims and scopes.
func (c *HydraClient) AcceptConsentSession(consentChallenge string, acceptReq HydraPutAcceptConsentRequest) (*HydraPutAcceptConsentResponse, error) {
	jsonData, err := json.Marshal(acceptReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal accept consent request: %w", err)
	}

	url := fmt.Sprintf("%s/admin/oauth2/auth/requests/consent/accept?consent_challenge=%s", c.adminURL, consentChallenge)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create accept consent request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to accept consent session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hydra accept consent API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result HydraPutAcceptConsentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode accept consent response: %w", err)
	}

	return &result, nil
}

// RejectConsentSession rejects a Hydra consent session with the specified error.
func (c *HydraClient) RejectConsentSession(consentChallenge string, errorCode string, errorDescription string) error {
	rejectReq := map[string]any{
		"error":             errorCode,
		"error_description": errorDescription,
	}

	jsonData, err := json.Marshal(rejectReq)
	if err != nil {
		return fmt.Errorf("failed to marshal reject consent request: %w", err)
	}

	url := fmt.Sprintf("%s/admin/oauth2/auth/requests/consent/reject?consent_challenge=%s", c.adminURL, consentChallenge)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create reject consent request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to reject consent session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("hydra reject consent API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
