package hydra

import (
	"fmt"
	"net/http"
	"time"

	"conformitea/server/config"
)

// ProvideHydraClient creates a new HydraClient instance based on the provided configuration.
// This replaces the singleton Initialize pattern with proper dependency injection.
func ProvideHydraClient(config config.HydraConfig) (*HydraClient, error) {
	if config.AdminURL == "" {
		return nil, fmt.Errorf("hydra admin URL is not configured")
	}

	client := &HydraClient{
		adminURL: config.AdminURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return client, nil
}
