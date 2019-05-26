package server

import (
	"github.com/gin-gonic/gin"
)

// NewGinServer creates a new default Gin Engine server.
func NewGinServer() *gin.Engine {
	engine := gin.Default()
	return engine
}
