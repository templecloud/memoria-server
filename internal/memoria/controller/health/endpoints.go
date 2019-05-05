package health

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

// Health is a Gin handler that reports the health of the API.
func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "ALIVE",
	})
}