package server

import (
	"conformitea/server/internal"
	"conformitea/server/types"
)

func Initialize(c types.Config) (types.Server, error) {
	return internal.Initialize(c)
}
