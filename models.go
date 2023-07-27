package auth

import "github.com/golang-jwt/jwt"

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type User struct {
	Username string `bson:"username" binding:"required"`
	Password string `bson:"password" binding:"required"`
}
