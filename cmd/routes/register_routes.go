package routes

import (
	"github.com/gin-gonic/gin"
	"venecraft-back/cmd/controller"
)

func RegisterRoutes(router *gin.Engine, registerController *controller.RegisterController) {
	registerGroup := router.Group("/register")
	{
		registerGroup.POST("/", registerController.CreateRegister)

		registerGroup.PUT("/approve/:id", registerController.ApproveRegister)
	}
}
