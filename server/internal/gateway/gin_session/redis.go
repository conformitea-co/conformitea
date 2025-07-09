package gin_session

import (
	"fmt"

	"conformitea/server/config"

	"github.com/gin-contrib/sessions"
)

// Deprecated: Use ProvideRedisStore instead. This function is kept for backward compatibility.
func Initialize(redisConfig config.RedisConfig, httpServerConfig config.HTTPServerConfig) error {
	// This function is deprecated and should not be used.
	// Use the Wire-based ProvideRedisStore function instead.
	return nil
}

// Deprecated: Use injected redis store instead. This function is kept for backward compatibility.
func GetRedisStore() (sessions.Store, error) {
	// This function is deprecated and should not be used.
	// Use dependency injection to get the redis store instance.
	return nil, fmt.Errorf("deprecated function: use dependency injection to get redis store")
}
