package server

import (
	"conformitea/server/internal"
	"conformitea/server/types"

	"go.uber.org/zap"
)

func Initialize(c types.Config, l *zap.Logger, appAuth types.AppAuth) (types.Server, error) {
	return internal.Initialize(c, l, appAuth)
}
