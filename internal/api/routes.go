package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and returns a Gin router with all routes defined
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Add more routes here as needed

	return r
}
