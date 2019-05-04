package identity

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Retrieve the cookie token if it is present.
		tokenStr, err := ctx.Cookie("token")
		fmt.Printf("tokenStr: %s\n", tokenStr)
		if err != nil || tokenStr == "" {
			if err == http.ErrNoCookie {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Parse the token string into a token with claims.
		// var claims Claims
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		fmt.Printf("token: %+v\n", token)
		fmt.Printf("claims: %+v\n", claims)

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