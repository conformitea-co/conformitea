package types

import (
	"github.com/gin-gonic/gin"
)

type Server interface {
	GetRouter() *gin.Engine
}
