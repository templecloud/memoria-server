package boot

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/templecloud/memoria-server/internal/memoria/boot/logging"
	"github.com/templecloud/memoria-server/internal/memoria/boot/persistence"
	"github.com/templecloud/memoria-server/internal/memoria/boot/server"
	"github.com/templecloud/memoria-server/internal/memoria/controller"
)

// NewServer creates a new configured server.
func NewServer() *gin.Engine {
	// Initialise MongoDB client.
	mongo := persistence.NewMongoClient()
	// Initialise Gin server.
	server := server.NewGinServer()
	server = controller.ConfigureEndpoints(server, mongo, nil)
	return server
}

// Start initialises the Memoria API webserver.
func Start(config *Config) {
	logging.ConfigureDefault()
	server := NewServer()
	log.Info("Starting memoria-server...")
	server.Run()
}
