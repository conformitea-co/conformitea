package server

import (
	"conformitea/server/config"
	"conformitea/server/internal"
	"conformitea/server/types"

	"go.uber.org/zap"
)

func Initialize(c config.Config, l *zap.Logger, appAuth types.AppAuth) (types.Server, error) {
	return internal.Initialize(c, l, appAuth)
}
