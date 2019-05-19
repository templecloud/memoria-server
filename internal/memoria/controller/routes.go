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
	HealthEndpoint = "/health"
	// LoginEndpoint denotes the 'login' endpoint.
	LoginEndpoint = "/login"
	// SignupEndpoint denotes the 'signup' endpoint.
	SignupEndpoint = "/signup"
)

// ConfigureEndpoints configures a new Gin Engine server.
func ConfigureEndpoints(
	engine *gin.Engine, 
	mongoClient *mongo.Client,
) *gin.Engine {

	// Create Endpoint handlers.
	healthAPI := &health.API{}
	identityAPI := identity.NewAPI(mongoClient)

	// Non-authenticated routes.
	public := engine.Group(APIv1)
	public.POST(LoginEndpoint, identityAPI.Login)
	public.POST(SignupEndpoint, identityAPI.Signup)

	// Authenticated routes.
	private := engine.Group(APIv1)
	private.Use(identity.JWTMiddleware())
	private.GET(HealthEndpoint, healthAPI.Health)

	return engine
}