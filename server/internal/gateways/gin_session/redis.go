package gin_session

import (
	"net/http"

	"conformitea/server/internal/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func NewStore() (sessions.Store, error) {
	config := config.GetConfig()

	var keyBytes [][]byte
	for _, k := range config.HTTPServer.Session.KeyPairs {
		keyBytes = append(keyBytes, []byte(k))
	}

	store, err := redis.NewStore(10, "tcp", config.Redis.Address, config.Redis.User, config.Redis.Password, keyBytes...)
	if err != nil {
		return nil, err
	}

	store.Options(sessions.Options{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   config.HTTPServer.Session.Timeout,
	})

	return store, nil
}
