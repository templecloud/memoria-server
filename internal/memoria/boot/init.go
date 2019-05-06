package boot

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/templecloud/memoria-server/internal/memoria/controller/health"
	"github.com/templecloud/memoria-server/internal/memoria/controller/identity"
	"github.com/templecloud/memoria-server/internal/memoria/persistence"
)

// Start initialises the Memoria API webserver.
func Start() {
	fmt.Println("Starting memoria-server...")
	router := NewServer()
	router.Run()
}

// NewServer creates a new Gin Engine server.
func NewServer() *gin.Engine {
	// Initialise resources.
	mongo := persistence.MongoClient()
	healthAPI := &health.API{}
	identityAPI := identity.NewAPI(mongo)

	// Define the routes and middleware
	//
	router := gin.Default()
	// Non-authenticated routes.
	public := router.Group("/api/v1")
	public.POST("/signup", identityAPI.Signup)
	public.POST("/login", identityAPI.Login)
	// Authenticated routes.
	private := router.Group("/api/v1")
	private.Use(identity.JWTMiddleware())
	private.GET("/health", healthAPI.Health)

	return router
}

