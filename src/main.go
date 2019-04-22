package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting memoria-server...")

	r := gin.Default();
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{ 
			"health": "ALIVE", 
		})
	})

	r.Run()
}
