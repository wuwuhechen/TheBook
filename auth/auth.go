package auth

import "github.com/golang-jwt/jwt/v5"

var jwtSecret []byte

type Claims struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
