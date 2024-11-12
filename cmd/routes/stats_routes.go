package routes

import (
	"github.com/gin-gonic/gin"
	"venecraft-back/cmd/controller"
)

func ServerStatsRoutes(router *gin.RouterGroup, serverStatsController *controller.ServerStatsController) {
	adminGroup := router.Group("/server")
	{
		adminGroup.GET("/stats", serverStatsController.GetServerStats)
		adminGroup.GET("/stats/pdf", serverStatsController.GeneratePDFReport)
	}
}
