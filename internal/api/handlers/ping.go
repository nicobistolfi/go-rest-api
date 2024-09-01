package handlers

import (
	"github.com/gin-gonic/gin"
)

// PingHandler handles the GET request for the /ping route
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
