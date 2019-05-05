package identity

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("todo-create-managed-key")

// jwtClaims contains the 'claim' details for authorizing via JWT token. 
type jwtClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

//-------------------------------------------------------------------------------------------------
// Public Functions

// JWTMiddleware is a Gin middleware handler function that handles the extraction and validation of 
// JWT token cookies.
func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Retrieve the cookie token if it is present.
		tokenStr, err := ctx.Cookie("token")
		if err != nil || tokenStr == "" {
			if err == http.ErrNoCookie {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Parse the token string into a token with claims.
		claims := &jwtClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// Check the token is valid (authenticates and has not expired)
		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
}