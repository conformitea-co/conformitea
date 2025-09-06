package hydra

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// JWTParser fetches and caches JWKS from Hydra and provides a jwt.Keyfunc compatible lookup.
type JWTParser struct {
	jwksURL         string
	set             jwk.Set
	mu              sync.RWMutex
	client          *http.Client
	refreshInterval time.Duration
	stopCh          chan struct{}
	wg              sync.WaitGroup
}

// NewJWTParserFromURL creates a new JWT parser that fetches JWKS from the provided URL and starts background refresh.
func NewJWTParserFromURL(jwksURL string, httpClient *http.Client) (*JWTParser, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	parser := &JWTParser{
		jwksURL:         jwksURL,
		client:          httpClient,
		refreshInterval: 1 * time.Hour,
		stopCh:          make(chan struct{}),
	}

	// initial fetch
	if err := parser.refresh(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS from %s: %w", jwksURL, err)
	}

	// start background refresher
	parser.wg.Add(1)
	go func() {
		defer parser.wg.Done()
		ticker := time.NewTicker(parser.refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = parser.refresh(context.Background())
			case <-parser.stopCh:
				return
			}
		}
	}()

	return parser, nil
}

// Close stops the background refresh goroutine.
func (j *JWTParser) Close() {
	close(j.stopCh)
	j.wg.Wait()
}

// refresh fetches the latest JWKS and replaces the in-memory set.
func (j *JWTParser) refresh(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, j.jwksURL, nil)
	if err != nil {
		return err
	}
	if j.client == nil {
		j.client = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := j.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	set, err := jwk.ParseReader(resp.Body)
	if err != nil {
		return err
	}

	j.mu.Lock()
	j.set = set
	j.mu.Unlock()
	return nil
}

// keyfuncForJWT provides a jwt.Keyfunc compatible function for github.com/golang-jwt/jwt/v5.
func (j *JWTParser) keyfuncForJWT(token *jwt.Token) (any, error) {
	// extract kid
	kidVal, ok := token.Header["kid"]
	if !ok {
		return nil, fmt.Errorf("token header missing kid")
	}
	kid := fmt.Sprint(kidVal)

	// attempt to get key from current set
	j.mu.RLock()
	var key jwk.Key
	if j.set != nil {
		key, _ = j.set.LookupKeyID(kid)
	}
	j.mu.RUnlock()

	if key == nil {
		// try a single immediate refresh then lookup again
		if err := j.refresh(context.Background()); err == nil {
			j.mu.RLock()
			if j.set != nil {
				key, _ = j.set.LookupKeyID(kid)
			}
			j.mu.RUnlock()
		}
	}

	if key == nil {
		return nil, fmt.Errorf("no key found for kid %s", kid)
	}

	// use the key
	var pubkey rsa.PublicKey
	if err := key.Raw(&pubkey); err != nil {
		// try to decode into interface (some keys might be returned as *rsa.PublicKey)
		var anyKey any
		if err2 := key.Raw(&anyKey); err2 == nil {
			switch k := anyKey.(type) {
			case *rsa.PublicKey:
				return k, nil
			case rsa.PublicKey:
				return &k, nil
			}
		}
		return nil, fmt.Errorf("failed to get raw key for kid %s: %w", kid, err)
	}

	return &pubkey, nil
}

// ParseAndValidateJWT parses and validates a JWT token using the JWKS-backed keyfunc.
func (j *JWTParser) ParseAndValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse the token using the JWKS Keyfunc
	token, err := jwt.Parse(tokenString, j.keyfuncForJWT)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from JWT")
	}

	// Validate expiration claim (exp)
	if expVal, ok := claims["exp"]; ok {
		switch exp := expVal.(type) {
		case float64:
			if int64(exp) < time.Now().Unix() {
				return nil, fmt.Errorf("JWT token expired")
			}
		case int64:
			if exp < time.Now().Unix() {
				return nil, fmt.Errorf("JWT token expired")
			}
		case int:
			if int64(exp) < time.Now().Unix() {
				return nil, fmt.Errorf("JWT token expired")
			}
		default:
			// Unsupported exp type — skip strict validation
		}
	}

	return claims, nil
}

// ParseJWTWithoutValidation parses a JWT without validating the signature (for extracting claims only)
func (j *JWTParser) ParseJWTWithoutValidation(tokenString string) (jwt.MapClaims, error) {
	// Parse without validation
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from JWT")
	}

	return claims, nil
}
