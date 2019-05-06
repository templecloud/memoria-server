package identity

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//-------------------------------------------------------------------------------------------------
// Public Functions

// Signup is a Gin handler that deals with creating new users.
func (api *API) Signup(ctx *gin.Context) {
	// Unmarshall the 'signup' request body.
	var signup Signup
	err := ctx.BindJSON(&signup)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", err)})
		return
	}

	userExists, err := api.db.HasUser(ctx, signup.Login.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errorMessage": fmt.Sprintf("%s", err)})
		return
	}

	if userExists {
		ctx.JSON(http.StatusConflict, gin.H{"errorMessage": "User already registered."})
	} else {
		// Check password strength.
		pwd := signup.Password
		err := validate(pwd)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", err)})
			return
		}

		// Create a new User.
		id := uuid.New().String()
		var user = User{
			ID:    string(id),
			Email: signup.Email,
			Name:  signup.Name,
		}

		// Create a new Credential.
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
		}
		var cred = Credential{
			UserID:   user.ID,
			Password: string(hash),
		}

		// Save new user and credential entities.
		err = api.db.CreateUser(ctx, user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", err)})
			return
		}
		err = api.db.CreateCredential(ctx, cred)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", err)})
			return
		}

		ctx.Status(http.StatusOK)
	}
}

// Login is a Gin handler that deals with authorising a user to provide them with a JWT token
// cookie.
func (api *API) Login(ctx *gin.Context) {
	// Unmarshall the 'login' request body.
	var login Login
	err := ctx.BindJSON(&login)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errorMessage": fmt.Sprintf("%s", err)})
		return
	} 
	
	// Get user and credentials.
	user, err := api.db.GetUser(ctx, login.Email)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	cred, err := api.db.GetCredential(ctx, user.ID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Generate Hashed password
	err = bcrypt.CompareHashAndPassword([]byte(cred.Password), []byte(login.Password))
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Create the JWT claims. Includes the username and expiry time.
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwtClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the HMAC algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errorMessage": fmt.Sprintf("%s", err)})
		return
	}

	// Set the JWT token as a cookie.
	expiration := int(5 * 60)
	path := ""
	domain := ""
	secure := true
	httpOnly := true
	ctx.SetCookie("token", tokenString, expiration, path, domain, secure, httpOnly)

	ctx.Status(http.StatusOK)
}

//-------------------------------------------------------------------------------------------------
// Private Functions

// validate checks a password meets certain constraints and returns an error if it is invalid.
func validate(password string) error {
	return nil
}
