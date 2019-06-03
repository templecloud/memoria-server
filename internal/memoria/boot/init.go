package boot

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/templecloud/memoria-server/internal/memoria/boot/logging"
	"github.com/templecloud/memoria-server/internal/memoria/boot/persistence"
	"github.com/templecloud/memoria-server/internal/memoria/boot/server"
	"github.com/templecloud/memoria-server/internal/memoria/controller"
)

// NewDefaultServer creates a new default configured server.
func NewDefaultServer() *gin.Engine {
	return NewServer(NewDefaultConfig())
}

// NewServer creates a new configured server.
func NewServer(config *Config) *gin.Engine {
	// Initialise MongoDB client.
	mongo := persistence.NewMongoClient(config.Persistence)
	// Initialise Gin server.
	server := server.NewGinServer()
	server = controller.ConfigureEndpoints(server, mongo, config.Controller)
	return server
}

// Start initialises the Memoria API webserver.
func Start(config *Config) {
	logging.Configure(config.Logging)
	server := NewServer(config)
	log.Info("Starting memoria-server...")
	server.Run()
}
