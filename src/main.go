package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// API is the main entry point into the API.
type API struct {
	db *mongo.Client
}

func (api *API) health(c *gin.Context) {
	c.JSON(200, gin.H{
		"health": "ALIVE",
	})
}

// Signup denotes the minimum requires details for logging in.
type Signup struct {
	Name string `form:"name" json:"name" binding:"required"`
	Login
}

func (api *API) signup(ctx *gin.Context) {
	var signup Signup
	unmarshallErr := ctx.BindJSON(&signup)
	if unmarshallErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", unmarshallErr)})
	} else {
		db := api.db.Database("memoria")
		col := db.Collection("users")
		filter := bson.M{"name": bson.M{"$eq": signup.Login.Email}}

		var existingUser Signup
		if readErr := col.FindOne(ctx, filter).Decode(&existingUser); readErr != nil {
			_, createErr := col.InsertOne(ctx, signup)
			if createErr != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", createErr)})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
			}
		} else {
			ctx.JSON(http.StatusConflict, gin.H{"errorMessage": "User already registered."})
		}
	}
}

// Login denotes the minimum requires details for registering.
type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (api *API) login(ctx *gin.Context) {
	var login Login
	unmarshallErr := ctx.BindJSON(&login)
	if unmarshallErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", unmarshallErr)})
	} else {
		db := api.db.Database("memoria")
		col := db.Collection("users")
		filter := bson.M{"name": bson.M{"$eq": login.Email}}

		var user Signup
		if readErr := col.FindOne(ctx, filter).Decode(&user); readErr != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"errorMessage": fmt.Sprintf("%s", readErr)})

		} else {
			if login.Email == user.Email && login.Password == user.Password {
				ctx.JSON(http.StatusOK, gin.H{"status": "Authorised"})
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorised"})
			}
		}
	}
}

func mongoClient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Failed to create client.")
		log.Fatal(err)
	}

	if err := client.Connect(context.TODO()); err != nil {
		fmt.Println("Failed to connect to mongodb.")
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Failed to ping mongodb.")
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

func main() {
	fmt.Println("Starting memoria-server...")

	db := mongoClient()
	api := &API{db: db}
	router := gin.Default()

	router.GET("/health", api.health)
	router.POST("/signup", api.signup)
	router.POST("/login", api.login)
	router.Run()
}
