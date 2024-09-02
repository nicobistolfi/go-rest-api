package api

import (
	"time"

	"go-boilerplate/internal/api/middleware"
	"go-boilerplate/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func SetupRouter(r *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// Add global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.RateLimiter(rate.Every(time.Second), 10)) // 10 requests per second

	// Public routes
	public := r.Group("/api/v1")
	{
		public.GET("/health", HealthCheck)
		public.GET("/ping", Ping)
		// Add other public routes
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.VerifyToken())
	{
		protected.GET("/profile", GetProfile)
		// Add other protected routes
	}
}
