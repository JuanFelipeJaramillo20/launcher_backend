package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
	"venecraft-back/cmd/dto"
	"venecraft-back/cmd/repository"
)

type AuthService interface {
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

var jwtKey = []byte("secret_key_not_secure_at_all")

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}

	claims := &dto.JWTCustomClaims{
		UserID: user.ID,
		Role:   roleNames,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
