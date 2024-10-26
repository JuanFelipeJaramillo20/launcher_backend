package dto

import "github.com/dgrijalva/jwt-go"

type JWTCustomClaims struct {
	UserID uint64   `json:"user_id"`
	Role   []string `json:"roles"`
	jwt.StandardClaims
}
