package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"venecraft-back/cmd/service"
)

type AuthController struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService}
}

func (ac *AuthController) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := ac.AuthService.Login(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
