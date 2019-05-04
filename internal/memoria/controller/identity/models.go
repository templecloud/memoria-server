package identity

import (
	"github.com/dgrijalva/jwt-go"
)

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

// Login denotes the minimum requires details for registering.
type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Claims are the JWT token claims.
type JWTClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

// Create the JWT key used to create the signature
var jwtKey = []byte("todo-create-managed-key")