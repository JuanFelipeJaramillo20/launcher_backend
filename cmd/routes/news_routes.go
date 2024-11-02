package routes

import (
	"github.com/gin-gonic/gin"
	"venecraft-back/cmd/controller"
)

func UserRoutes(router *gin.RouterGroup, newsController *controller.NewsController) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", newsController.CreateNews)
		userGroup.GET("/", newsController.GetAllNews)
		userGroup.GET("/:id", newsController.GetNewsByID)
		userGroup.PUT("/:id", newsController.UpdateNews)
		userGroup.DELETE("/:id", newsController.DeleteNews)
	}
}
