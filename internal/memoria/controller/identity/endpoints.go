package identity

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Signup is a Gin handler that deals with creating new users.
func (api *API) Signup(ctx *gin.Context) {
	var signup Signup
	unmarshallErr := ctx.BindJSON(&signup)
	if unmarshallErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", unmarshallErr)})
	} else {

		db := api.DB.Database("memoria")
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

// Login is a Gin handler that deals with authorising a user to provide them with a JWT token 
// cookie.
func (api *API) Login(ctx *gin.Context) {
	var login Login
	unmarshallErr := ctx.BindJSON(&login)
	if unmarshallErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", unmarshallErr)})
	} else {
		db := api.DB.Database("memoria")
		users := db.Collection("users")
		creds := db.Collection("credentials")
		userQuery := bson.M{"email": bson.M{"$eq": login.Email}}

		var user User
		if ue := users.FindOne(ctx, userQuery).Decode(&user); ue == nil {
			credQuery := bson.M{"userId": bson.M{"$eq": user.ID}}
			var cred Credential
			if ce := creds.FindOne(ctx, credQuery).Decode(&cred); ce == nil {
				if login.Email == user.Email {
					// Generate Hashed password
					loginErr := bcrypt.CompareHashAndPassword([]byte(cred.Password), []byte(login.Password))
					if loginErr == nil {
						// Create the JWT claims, which includes the username and expiry time
						expirationTime := time.Now().Add(5 * time.Minute)
						claims := &JWTClaims{
							UserID: user.ID,
							StandardClaims: jwt.StandardClaims{
								// In JWT, the expiry time is expressed as unix milliseconds
								ExpiresAt: expirationTime.Unix(),
							},
						}
						// Declare the token with the HMAC algorithm used for signing, and the claims
						token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
						// Create the JWT string
						tokenString, err := token.SignedString(jwtKey)
						if err != nil {
							// If there is an error in creating the JWT return an internal server error
							ctx.JSON(http.StatusInternalServerError, gin.H{"errorMessage": fmt.Sprintf("%s", unmarshallErr)})
							return
						}
						// Finally, we set the client cookie for "token" as the JWT we just generated
						// we also set an expiry time which is the same as the token itself
						expiration := int(5 * 60)
						path := ""
						domain := ""
						secure := true
						httpOnly := true
						ctx.SetCookie("token", tokenString, expiration, path, domain, secure, httpOnly)

						ctx.JSON(http.StatusOK, gin.H{"status": "Authorised"})
						return
					}
				}
			}
		}

		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorised"})
	}
}


// validate checks a password meets certain constraints and returns an error if it is invalid.
func validate(password string) error {
	return nil
}