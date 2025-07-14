package gateway

import (
	"conformitea/infrastructure/gateway/hydra"
	"conformitea/infrastructure/gateway/microsoft"
)

func Initialize(hydraAdminURL, msClientID, msClientSecret, msRedirectURL string, msScopes []string) {
	hydra.Initialize(hydraAdminURL)
	microsoft.Initialize(msClientID, msClientSecret, msRedirectURL, msScopes)
}
