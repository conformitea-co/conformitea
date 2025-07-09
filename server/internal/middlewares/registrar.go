package middlewares

import (
	"conformitea/server/internal/types"

	"github.com/gin-gonic/gin"
)

func RegisterMiddlewares(server types.InternalServer) {
	router := server.GetRouter()
	sessionMiddleware := SessionMiddleware(server)

	// Most of the time, the order of middlewares is important.
	router.Use(LogRequestDetails(server))
	router.Use(ContextLoggerMiddleware(server))
	router.Use(RequestIdMiddleware())
	router.Use(CORSMiddleware())
	router.Use(sessionMiddleware)
	router.Use(gin.Recovery())
}
