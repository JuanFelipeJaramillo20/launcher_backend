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

// Request model for registration details
// swagger:model RegisterRequest
type RegisterRequest struct {
	// Full name for registration
	// required: true
	FullName string `json:"full_name"`

	// Email address for registration
	// required: true
	// example: user@example.com
	Email string `json:"email"`

	// Nickname chosen for registration
	// required: true
	Nickname string `json:"nickname"`

	// Password for registration
	// required: true
	// min length: 8
	Password string `json:"password"`
}

// Parameters for creating a registration request
// swagger:parameters createRegister
type RegisterParams struct {
	// Registration details
	// in: body
	// required: true
	Body RegisterRequest
}

// Parameters for approving or denying a registration request
// swagger:parameters approveRegister denyRegister
type RegisterActionParams struct {
	// ID of the registration request
	// in: path
	// required: true
	ID uint64 `json:"id"`
}

type RegisterController struct {
	RegisterService service.RegisterService
}

func NewRegisterController(registerService service.RegisterService) *RegisterController {
	return &RegisterController{registerService}
}

// swagger:route POST /api/register register createRegister
// Creates a registration request for a new user.
//
// responses:
//
//	201: CommonSuccess
//	400: CommonError
//	500: CommonError
func (rc *RegisterController) CreateRegister(c *gin.Context) {
	var register entity.Register
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := rc.RegisterService.CreateRegister(&register)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration request created successfully"})
}

// swagger:route GET /api/register register getAllRegisters
// Retrieves all registers.
//
// Security:
//   - BearerAuth: []
//
// Responses:
//
//	200: []Register
//	403: CommonError
//	500: CommonError
func (rc *RegisterController) GetAllRegisters(c *gin.Context) {
	_, roles, loggedIn := middlewares.GetLoggedInUser(c)
	if !loggedIn || !slices.Contains(roles, "ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	registers, err := rc.RegisterService.GetAllRegisters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, registers)
}

// swagger:route PUT /api/register/approve/{id} register approveRegister
// Approves a user registration request by ID.
//
// Security:
//   - BearerAuth: []
//
// responses:
//
//	200: User
//	400: CommonError
//	500: CommonError
func (rc *RegisterController) ApproveRegister(c *gin.Context) {
	_, roles, loggedIn := middlewares.GetLoggedInUser(c)
	if !loggedIn || !slices.Contains(roles, "ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration ID"})
		return
	}

	user, err := rc.RegisterService.ApproveRegister(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// swagger:route PUT /api/register/deny/{id} register denyRegister
// Denies a user registration request by ID.
//
// Security:
//   - BearerAuth: []
//
// responses:
//
//	200: CommonSuccess
//	400: CommonError
//	500: CommonError
func (rc *RegisterController) DenyRegister(c *gin.Context) {
	_, roles, loggedIn := middlewares.GetLoggedInUser(c)
	if !loggedIn || !slices.Contains(roles, "ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration ID"})
		return
	}

	err = rc.RegisterService.DenyRegister(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Register request denied successfully")
}
