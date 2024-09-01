package api

import (
	"go-boilerplate/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and returns a Gin router with all routes defined
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", handlers.PingHandler)

	// Add more routes here as needed

	return r
}
