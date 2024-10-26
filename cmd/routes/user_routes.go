package routes

import (
	"github.com/gin-gonic/gin"
	"venecraft-back/cmd/controller"
)

func UserRoutes(router *gin.RouterGroup, userController *controller.UserController) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", userController.CreateUser)
		userGroup.GET("/", userController.GetAllUsers)
		userGroup.GET("/:id", userController.GetUserByID)
		userGroup.PUT("/:id", userController.UpdateUser)
		userGroup.DELETE("/:id", userController.DeleteUser)
	}
}
