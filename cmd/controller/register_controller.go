package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/service"
)

type RegisterController struct {
	RegisterService service.RegisterService
}

func NewRegisterController(registerService service.RegisterService) *RegisterController {
	return &RegisterController{registerService}
}

// swagger:route POST /api/register register createRegister
// Creates a registration request for a new user.
//
// Responses:
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

// swagger:route PUT /api/register/approve/{id} register approveRegister
// Approves a user registration request by ID.
//
// Responses:
//
//	200: User
//	400: CommonError
//	500: CommonError
func (rc *RegisterController) ApproveRegister(c *gin.Context) {
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
