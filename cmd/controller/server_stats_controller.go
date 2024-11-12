package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"venecraft-back/cmd/service"
)

type ServerStatsController struct {
	ServerStatsService service.ServerStatsService
}

func NewServerStatsController(serverStatsService service.ServerStatsService) *ServerStatsController {
	return &ServerStatsController{ServerStatsService: serverStatsService}
}

// GetServerStats returns server metrics in JSON format
// @Summary Get server statistics
// @Description Fetches server metrics including CPU usage, active users, and transactions
// @Tags admin
// @Produce json
// @Success 200 {object} service.ServerMetrics
// @Failure 500 {object} gin.H{"error": "Unable to fetch server metrics"}
// @Router /admin/stats [get]
func (a *ServerStatsController) GetServerStats(c *gin.Context) {
	metrics, err := a.ServerStatsService.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch server metrics"})
		return
	}
	c.JSON(http.StatusOK, metrics)
}

// GeneratePDFReport returns server metrics as a downloadable PDF
// @Summary Generate server statistics report
// @Description Generates a PDF report of server statistics including CPU usage, active users, and transactions
// @Tags admin
// @Produce application/pdf
// @Success 200 {file} server_stats.pdf
// @Failure 500 {object} gin.H{"error": "Failed to generate PDF"}
// @Router /admin/stats/pdf [get]
func (a *ServerStatsController) GeneratePDFReport(c *gin.Context) {
	metrics, err := a.ServerStatsService.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch server metrics"})
		return
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Server Statistics Report")

	// Add metrics to the PDF
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("CPU Usage: %.2f%%", metrics.CPUUsage))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Memory Usage: %.2f%%", metrics.MemoryUsage))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Disk Usage: %.2f%%", metrics.DiskUsage))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Active Users: %d", metrics.ActiveUsers))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Transactions: %d", metrics.Transactions))

	var buffer bytes.Buffer
	err = pdf.Output(&buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=server_stats.pdf")
	c.Data(http.StatusOK, "application/pdf", buffer.Bytes())
}
