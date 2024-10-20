package controller

import (
	"net/http"
	"strconv"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/service"

	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	RegisterService service.RegisterService
}

func NewRegisterController(registerService service.RegisterService) *RegisterController {
	return &RegisterController{registerService}
}

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
