package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"venecraft-back/cmd/entity"
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
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := uc.UserService.CreateUser(&user)
	if err != nil {
		switch err.Error() {
		case "invalid email format":
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
