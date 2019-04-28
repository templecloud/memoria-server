package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

// User is a `user` of the system.
type User struct {
	ID    string `json:"id" bson:"id" binding:"required"`
	Email string `json:"email" bson:"email" binding:"required"`
	Name  string `json:"name" bson:"name" binding:"optional"`
}

// Credential is the `encrypted password` of a `user`.
type Credential struct {
	UserID   string `json:"userId" bson:"userId" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

// The signup Gin handler deals with creating new users.
func (api *API) signup(ctx *gin.Context) {
	var signup Signup
	unmarshallErr := ctx.BindJSON(&signup)
	if unmarshallErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", unmarshallErr)})
	} else {

		db := api.db.Database("memoria")
		users := db.Collection("users")
		creds := db.Collection("credentials")
		userQuery := bson.M{"email": bson.M{"$eq": signup.Login.Email}}

		var user User
		// check if the email exists...
		if readErr := users.FindOne(ctx, userQuery).Decode(&user); readErr != nil {

			id := uuid.New().String()
			var newUser = User{
				ID:    string(id),
				Email: signup.Email,
				Name:  signup.Name,
			}

			pwd := signup.Password
			pwe := validate(pwd)
			if pwe != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", pwe)})
				return
			}


			hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
			if err != nil {
				log.Println(err)
			}

			var cred = Credential{
				UserID:   newUser.ID,
				Password: string(hash),
			}

			// Create new user and credential entities.
			_, cue := users.InsertOne(ctx, newUser)
			if cue != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", cue)})
			}
			_, iue := creds.InsertOne(ctx, cred)
			if iue != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", iue)})
			}

			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})

		} else {
			ctx.JSON(http.StatusConflict, gin.H{"errorMessage": "User already registered."})
		}
	}
}

// Validate checks a password and returns an error if it is invalid.
func validate(password string) error {
	return nil
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
		users := db.Collection("users")
		creds := db.Collection("credentials")
		userQuery := bson.M{"email": bson.M{"$eq": login.Email}}

		var user User
		if ue := users.FindOne(ctx, userQuery).Decode(&user); ue == nil {
			credQuery := bson.M{"userId": bson.M{"$eq": user.ID}}
			var cred Credential
			if ce := creds.FindOne(ctx, credQuery).Decode(&cred); ce == nil {
				if login.Email == user.Email {
					bcrypt.CompareHashAndPassword([]byte(cred.Password), []byte(login.Password))
					ctx.JSON(http.StatusOK, gin.H{"status": "Authorised"})
					return
				}
			}
		}

		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorised"})
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
