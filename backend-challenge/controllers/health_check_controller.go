package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var HealthCheckController = &healthCheckController{}

type healthCheckController struct {
}

// HealthCheck returns a health check response with the current time, environment, port, log level, and version.
// The response is returned with a 200 OK status code.
// @Summary Health check
// @Schemes http https
// @Tags health-check
// @Produce json
// @Success 200 {object} interface{}
// @Router /health [get]
func (h *healthCheckController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Kart service is running",
		"time":    time.Now().Format(time.RFC3339),
	})
}
