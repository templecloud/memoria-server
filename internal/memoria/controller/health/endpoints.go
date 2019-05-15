package health

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Health is a Gin handler that reports the health of the API.
func (api *API) Health(c *gin.Context) {
	log.WithFields(log.Fields{
		"clientIP": c.ClientIP(),
		"url": c.Request.URL,
	  }).Info("Health Check")
	
	c.JSON(http.StatusOK, gin.H{
		"health": "ALIVE",
	})
}