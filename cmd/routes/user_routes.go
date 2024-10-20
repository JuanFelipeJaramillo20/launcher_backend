package routes

import (
	"github.com/gin-gonic/gin"
	"venecraft-back/cmd/controller"
	"venecraft-back/cmd/middleware"
)

func UserRoutes(router *gin.Engine, userController *controller.UserController) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", userController.CreateUser)
		userGroup.GET("/", middleware.JWTMiddleware(), userController.GetAllUsers)
		userGroup.GET("/:id", middleware.JWTMiddleware(), userController.GetUserByID)
		userGroup.PUT("/:id", middleware.JWTMiddleware(), userController.UpdateUser)
		userGroup.DELETE("/:id", middleware.JWTMiddleware(), userController.DeleteUser)
	}
}
