package gin_session

import (
	"net/http"

	"conformitea/server/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func NewStore(cfg config.Config) (sessions.Store, error) {
	var keyBytes [][]byte
	for _, k := range cfg.HTTPServer.Session.KeyPairs {
		keyBytes = append(keyBytes, []byte(k))
	}

	store, err := redis.NewStore(10, "tcp", cfg.Redis.Address, cfg.Redis.User, cfg.Redis.Password, keyBytes...)
	if err != nil {
		return nil, err
	}

	store.Options(sessions.Options{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   cfg.HTTPServer.Session.Timeout,
	})

	return store, nil
}
