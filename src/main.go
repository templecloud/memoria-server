package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login denotes the minimum requires details for logging in.
type Signup struct {
	Name string `form:"name" json:"name" binding:"required"`
	Login
}

func signup(c *gin.Context) {
	fmt.Println("handling signup...")
	var signup Signup
	e := c.BindJSON(&signup)
	if e != nil {
		fmt.Println("signup error: ", e)
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", e)})
	} else {
		fmt.Println("signup Signup: ", signup)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// Login denotes the minimum requires details for registering.
type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{
		"health": "ALIVE",
	})
}

func login(c *gin.Context) {
	fmt.Println("handling login...")
	var login Login
	e := c.BindJSON(&login)
	if e != nil {
		fmt.Println("login error: ", e)
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", e)})
	} else {
		fmt.Printf("signin Login: %+v\n", login)
		if login.Email == "test" && login.Password == "test" {
			c.JSON(http.StatusOK, gin.H{"status": "authorised"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorised"})
		}
	}
}

func main() {
	fmt.Println("Starting memoria-server...")

	r := gin.Default()
	r.GET("/health", health)
	r.POST("/signup", signup)
	r.POST("/login", login)
	r.Run()
}
