package controller

import (
	"net/http"
	"venecraft-back/cmd/dto"
	"venecraft-back/cmd/service"

	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func SendCommand(client *rcon.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var commandRequest dto.CommandRequest

		if err := c.ShouldBindJSON(&commandRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Comando inválido"})
			return
		}

		response, err := service.SendCommand(client, commandRequest.Command)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al ejecutar el comando"})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

type ConsoleLogController struct {
	consoleLogService service.ConsoleLogService
}

func NewConsoleLogController() *ConsoleLogController {
	return &ConsoleLogController{
		consoleLogService: *service.NewConsoleLogService(),
	}
}

// Handler para streaming de logs
func (clc *ConsoleLogController) StreamLogs(c *gin.Context) {
	if err := clc.consoleLogService.StreamLogs(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
