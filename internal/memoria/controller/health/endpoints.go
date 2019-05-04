package health

import (
	"github.com/gin-gonic/gin"
)

func (api *API) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"health": "ALIVE",
	})
}