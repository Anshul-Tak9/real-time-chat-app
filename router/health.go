package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler returns the health status of the application
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Application is running",
	})
}
