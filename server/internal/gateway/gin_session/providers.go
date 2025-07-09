package gin_session

import (
	"net/http"

	"conformitea/server/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func ProvideRedisStore(redisConfig config.RedisConfig, httpServerConfig config.HTTPServerConfig) (sessions.Store, error) {
	var keyBytes [][]byte
	for _, k := range httpServerConfig.Session.KeyPairs {
		keyBytes = append(keyBytes, []byte(k))
	}

	store, err := redis.NewStore(10, "tcp", redisConfig.Address, redisConfig.User, redisConfig.Password, keyBytes...)
	if err != nil {
		return nil, err
	}

	store.Options(sessions.Options{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   httpServerConfig.Session.Timeout,
	})

	return store, nil
}
