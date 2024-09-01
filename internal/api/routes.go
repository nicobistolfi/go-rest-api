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
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.RateLimiter(rate.Every(time.Second), 10)) // 10 requests per second

	// Public routes
	public := r.Group("/api/v1")
	{
		public.GET("/health", HealthCheck)
		public.GET("/ping", Ping)
		// Add other public routes
	}

	// OAuth protected routes
	oauth := r.Group("/api/v1/oauth")
	oauth.Use(middleware.OAuthMiddleware())
	{
		oauth.GET("/profile", GetProfile)
		// Add other OAuth protected routes
	}

	// JWT protected routes
	jwt := r.Group("/api/v1/jwt")
	jwt.Use(middleware.JWTMiddleware())
	{
		jwt.GET("/profile", GetProfile)
		// Add other JWT protected routes
	}

	// API Key protected routes
	apiKey := r.Group("/api/v1/apikey")
	apiKey.Use(middleware.APIKeyMiddleware())
	{
		apiKey.GET("/profile", GetProfile)
		// Add other API Key protected routes
	}
}
