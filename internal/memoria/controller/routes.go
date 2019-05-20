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
	// HealthPath denotes the 'health' endpoint.
	HealthPath = "/" + HealthEndpoint

	// LoginEndpoint denotes the 'login' endpoint.
	LoginEndpoint = "login"
	// LoginPath denotes the 'login' endpoint.
	LoginPath = "/" + LoginEndpoint

	// SignupEndpoint denotes the 'signup' endpoint.
	SignupEndpoint = "signup"
	// SignupPath denotes the 'signup' endpoint.
	SignupPath= "/" + SignupEndpoint
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
		public.POST(LoginPath, identityAPI.Login)
	}
	if config.isEnabled(SignupEndpoint) {
		public.POST(SignupPath, identityAPI.Signup)
	}

	// Authenticated routes.
	private := engine.Group(APIv1)
	private.Use(identity.JWTMiddleware())
	if config.isEnabled(HealthEndpoint) {
		private.GET(HealthPath, healthAPI.Health)
	}

	return engine
}
