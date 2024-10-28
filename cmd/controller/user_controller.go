package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strconv"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/middlewares"
	"venecraft-back/cmd/service"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService}
}

// swagger:route POST /api/users users createUser
// Creates a new user.
//
// Responses:
//
//	201: CommonSuccess
//	400: CommonError
//	409: CommonError
//	500: CommonError
func (uc *UserController) CreateUser(c *gin.Context) {
	_, roles, loggedIn := middlewares.GetLoggedInUser(c)
	if !loggedIn || !slices.Contains(roles, "ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	role := c.Query("role")
	err := uc.UserService.CreateUser(&user, role)
	if err != nil {
		switch err.Error() {
		case "invalid email format":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "invalid nickname (only letters, numbers, and underscores allowed; must be 3-30 characters)":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "password must be at least 8 characters long":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "password must contain at least one uppercase letter":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "password must contain at least one digit":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "password must contain at least one special character":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "user with this email already exists":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// swagger:route GET /api/users users getAllUsers
// Retrieves all users.
//
// Responses:
//
//	200: []User
//	500: CommonError
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// swagger:route GET /api/users/{id} users getUserByID
// Retrieves a user by their ID.
//
// Responses:
//
//	200: User
//	400: CommonError
//	404: CommonError
func (uc *UserController) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uc.UserService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// swagger:route PUT /api/users/{id} users updateUser
// Updates user information.
//
// Responses:
//
//	200: CommonSuccess
//	400: CommonError
//	500: CommonError
func (uc *UserController) UpdateUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := uc.UserService.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// swagger:route DELETE /api/users/{id} users deleteUser
// Deletes a user by ID.
//
// Responses:
//
//	200: CommonSuccess
//	400: CommonError
//	500: CommonError
func (uc *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = uc.UserService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// PasswordResetRequest Initiates a password reset request.
// swagger:route POST /api/password-reset-request users passwordResetRequest
func (uc *UserController) PasswordResetRequest(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	err := uc.UserService.RequestPasswordReset(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate password reset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

// ResetPassword Resets the password using a token.
// swagger:route POST /api/reset-password users resetPassword
func (uc *UserController) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := uc.UserService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
