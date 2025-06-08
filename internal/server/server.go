package server

import (
	"github.com/conformitea-co/conformitea/internal/server/middlewares"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	middlewares.RegisterMiddlewares(router)
	RegisterRoutes(router)

	return router
}

func Start() error {
	router := NewRouter()

	return router.Run()
}
