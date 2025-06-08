package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SessionManagement(router *gin.Engine) {
	address := viper.GetString("redis.address")
	store, err := redis.NewStore(10, "tcp", address, "", "")
	if err != nil {
		panic(err)
	}

	router.Use(sessions.Sessions("conformitea_session", store))
}
