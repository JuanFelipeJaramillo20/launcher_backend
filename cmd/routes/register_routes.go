package routes

import (
	"github.com/gin-gonic/gin"
	"venecraft-back/cmd/controller"
)

func RegisterRoutes(router *gin.Engine, registerController *controller.RegisterController) {
	router.POST("/api/register", registerController.CreateRegister)
}
