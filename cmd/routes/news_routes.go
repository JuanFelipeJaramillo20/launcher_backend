package routes

import (
	"venecraft-back/cmd/controller"

	"github.com/gin-gonic/gin"
)

func NewsRoutes(router *gin.RouterGroup, newsController *controller.NewsController) {
	newsGroup := router.Group("/news")
	{
		newsGroup.POST("/", newsController.CreateNews)
		newsGroup.GET("/", newsController.GetAllNews)
		newsGroup.GET("/latest/", newsController.GetLatestNews)
		newsGroup.GET("/:id", newsController.GetNewsByID)
		newsGroup.PUT("/:id", newsController.UpdateNews)
		newsGroup.DELETE("/:id", newsController.DeleteNews)
		newsGroup.POST("/:id/reaction/:reactionType", newsController.ToggleReaction)
	}
}
