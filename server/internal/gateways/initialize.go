package gateways

import (
	"conformitea/server/internal/gateways/hydra"
	"conformitea/server/internal/gateways/microsoft"
)

func Initialize() {
	hydra.Initialize()
	microsoft.Initialize()
}
