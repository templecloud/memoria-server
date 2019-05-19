package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/templecloud/memoria-server/internal/memoria/controller/health"
	"github.com/templecloud/memoria-server/internal/memoria/controller/identity"
)

const (
	// APIv1 denotes the v1 endpoint.
	APIv1 = "/api/v1"

	// HealthEndpoint denotes the 'health' endpoint.
	HealthEndpoint = "health"
	// HealthRoute denotes the 'health' endpoint.
	HealthRoute = "/" + HealthEndpoint

	// LoginEndpoint denotes the 'login' endpoint.
	LoginEndpoint = "login"
	// LoginRoute denotes the 'login' endpoint.
	LoginRoute = "/" + LoginEndpoint

	// SignupEndpoint denotes the 'signup' endpoint.
	SignupEndpoint = "signup"
	// SignupRoute denotes the 'signup' endpoint.
	SignupRoute = "/" + SignupEndpoint
)

// ConfigureEndpoints configures a new Gin Engine server.
func ConfigureEndpoints(
	engine *gin.Engine,
	mongoClient *mongo.Client,
	config *Config,
) *gin.Engine {
	// Ensure Config.
	if config == nil {
		config = newDefaultConfig()
	}

	// Create Endpoint handlers.
	healthAPI := &health.API{}
	identityAPI := identity.NewAPI(mongoClient)

	// Non-authenticated routes.
	public := engine.Group(APIv1)
	if config.isEnabled(LoginEndpoint) {
		public.POST(LoginRoute, identityAPI.Login)
	}
	if config.isEnabled(SignupEndpoint) {
		public.POST(SignupRoute, identityAPI.Signup)
	}

	// Authenticated routes.
	private := engine.Group(APIv1)
	private.Use(identity.JWTMiddleware())
	if config.isEnabled(HealthEndpoint) {
		private.GET(HealthRoute, healthAPI.Health)
	}

	return engine
}
