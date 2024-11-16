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

// Request model for user creation
// swagger:model CreateUserRequest
type CreateUserRequest struct {
	// Full name of the user
	// required: true
	FullName string `json:"full_name"`

	// Email address of the user
	// required: true
	// example: user@example.com
	Email string `json:"email"`

	// Nickname for the user
	// required: true
	Nickname string `json:"nickname"`

	// Password for the user account
	// required: true
	// min length: 8
	Password string `json:"password"`
}

// Parameters for creating a user
// swagger:parameters createUser
type CreateUserParams struct {
	// User details for account creation
	// in: body
	// required: true
	Body CreateUserRequest
}

// Parameters for retrieving, updating, or deleting a user by ID
// swagger:parameters getUserByID updateUser deleteUser
type UserIDParams struct {
	// ID of the user
	// in: path
	// required: true
	ID uint64 `json:"id"`
}

// Request model for password reset initiation
// swagger:model PasswordResetRequest
type PasswordResetRequest struct {
	// Email address to send the reset link
	// required: true
	// example: user@example.com
	Email string `json:"email"`
}

// Parameters for initiating password reset
// swagger:parameters passwordResetRequest
type PasswordResetParams struct {
	// Password reset details
	// in: body
	// required: true
	Body PasswordResetRequest
}

// Request model for password reset
// swagger:model ResetPasswordRequest
type ResetPasswordRequest struct {
	// Token for password reset verification
	// required: true
	Token string `json:"token"`

	// New password for the user
	// required: true
	// min length: 8
	NewPassword string `json:"newPassword"`
}

// Parameters for resetting the password
// swagger:parameters resetPassword
type ResetPasswordParams struct {
	// Password reset details
	// in: body
	// required: true
	Body ResetPasswordRequest
}

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService}
}

// swagger:route POST /api/users users createUser
// Creates a new user.
//
// Security:
//   - BearerAuth: []
//
// responses:
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

	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user := entity.User{
		FullName: request.FullName,
		Email:    request.Email,
		Nickname: request.Nickname,
		Password: request.Password,
	}

	role := c.Query("role")
	err := uc.UserService.CreateUser(&user, role)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// swagger:route GET /api/users users getAllUsers
// Retrieves all users.
//
// Security:
//   - BearerAuth: []
//
// responses:
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
// Security:
//   - BearerAuth: []
//
// responses:
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
// Security:
//   - BearerAuth: []
//
// responses:
//
//	200: CommonSuccess
//	400: CommonError
//	500: CommonError
func (uc *UserController) UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var userUpdate entity.User
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = uc.UserService.UpdateUser(userID, &userUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// swagger:route DELETE /api/users/{id} users deleteUser
// Deletes a user by ID.
//
// Security:
//   - BearerAuth: []
//
// responses:
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

// swagger:route POST /api/password-reset-request users passwordResetRequest
// Initiates a password reset request.
//
// responses:
//
//	200: CommonSuccess
//	400: CommonError
//	500: CommonError
func (uc *UserController) PasswordResetRequest(c *gin.Context) {
	var req PasswordResetRequest
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

// swagger:route POST /api/reset-password users resetPassword
// Resets the password using a token.
//
// responses:
//
//	200: CommonSuccess
//	400: CommonError
//	500: CommonError
func (uc *UserController) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
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
