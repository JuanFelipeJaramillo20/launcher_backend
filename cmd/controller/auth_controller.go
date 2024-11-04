package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"venecraft-back/cmd/service"
)

type AuthController struct {
	AuthService service.AuthService
}

// Request model for login credentials
// swagger:model LoginRequest
type LoginRequest struct {
	// The email address for login
	// required: true
	// example: user@example.com
	Email string `json:"email"`

	// The password for login
	// required: true
	// min length: 8
	Password string `json:"password"`
}

// swagger:parameters loginUser
type LoginParams struct {
	// Login credentials
	// in: body
	// required: true
	Body LoginRequest
}

// swagger:response StringResponse
type StringResponse struct {
	// JWT token returned after successful authentication
	Token string `json:"token"`
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService}
}

// swagger:route POST /auth/login auth loginUser
// Logs in a user by email and password, returning a JWT token.
//
// Consumes:
//   - application/json
//
// Responses:
//
//	200: StringResponse
//	400: CommonError
//	401: CommonError
func (ac *AuthController) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := ac.AuthService.Login(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
