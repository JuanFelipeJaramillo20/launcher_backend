package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
	"venecraft-back/cmd/dto"
)

var jwtSecret = []byte("secret_key_not_secure_at_all")

func GenerateToken(userID uint64, role []string) (string, error) {
	claims := dto.JWTCustomClaims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*dto.JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*dto.JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
