package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SessionMiddleware() (gin.HandlerFunc, error) {
	address := viper.GetString("redis.address")
	username := viper.GetString("redis.username")
	password := viper.GetString("redis.password")
	keyPairs := viper.GetStringSlice("session.key_pairs")
	timeout := viper.GetInt("session.timeout")
	cookieName := viper.GetString("session.cookie_name")

	var keyBytes [][]byte
	for _, k := range keyPairs {
		keyBytes = append(keyBytes, []byte(k))
	}

	store, err := redis.NewStore(10, "tcp", address, username, password, keyBytes...)

	if err != nil {
		return nil, err
	}

	store.Options(sessions.Options{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   timeout,
	})

	return sessions.Sessions(cookieName, store), nil
}
