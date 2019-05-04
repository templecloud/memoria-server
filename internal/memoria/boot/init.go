package boot

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/templecloud/memoria-server/internal/memoria/controller/health"
	"github.com/templecloud/memoria-server/internal/memoria/controller/identity"
	"github.com/templecloud/memoria-server/internal/memoria/persistence"
	
)

func Start() {

	fmt.Println("Starting memoria-server...")

	db := persistence.MongoClient()

	healthAPI := &health.API{}
	identityAPI := &identity.API{DB: db}

	router := gin.Default()

	public := router.Group("/api/v1")
	public.POST("/signup", identityAPI.Signup)
	public.POST("/login", identityAPI.Login)

	private := router.Group("/api/v1")
	private.Use(identity.JWTMiddleware())
	private.GET("/health", healthAPI.Health)

	router.Run()

}