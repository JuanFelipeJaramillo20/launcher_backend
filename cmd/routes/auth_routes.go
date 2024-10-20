package routes

import (
	"github.com/gin-gonic/gin"
	"venecraft-back/cmd/controller"
)

func AuthRoutes(router *gin.Engine, authController *controller.AuthController) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
	}
}
