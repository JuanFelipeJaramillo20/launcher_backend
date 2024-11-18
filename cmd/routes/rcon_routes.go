package routes

import (
	"venecraft-back/cmd/controller"

	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func RCONRoutes(r *gin.RouterGroup, rconClient *rcon.Conn) {
	api := r.Group("/rcon")
	{
		api.POST("/command", controller.SendCommand(rconClient))
		api.GET("/logs", controller.NewConsoleLogController().StreamLogs)
	}
}
