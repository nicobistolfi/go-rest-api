package api

import (
	"time"

	"go-rest-api/internal/api/middleware"
	"go-rest-api/internal/config"

	logger "go-rest-api/pkg"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RouterOption func(*routerOptions)

type routerOptions struct {
	skipRateLimiting bool
}

func WithoutRateLimiting() RouterOption {
	return func(ro *routerOptions) {
		ro.skipRateLimiting = true
	}
}

func SetupRouter(router *gin.Engine, cfg *config.Config, logger *logger.Logger, opts ...RouterOption) {
	options := &routerOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Add global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware(logger))

	if options.skipRateLimiting {
		logger.Info("Rate limiting is disabled")
	} else {
		// Apply rate limiting middleware
		router.Use(middleware.RateLimiter(rate.Every(time.Second), 10)) // 10 requests per second
	}

	// Public routes
	public := router.Group("/api/v1")
	{
		public.GET("/health", HealthCheck)
		public.GET("/ping", Ping)
		public.POST("/register", Register)
		// Add other public routes
	}

	// Auth routes
	auth := router.Group("/api/v1")
	{
		auth.POST("/token", GetToken)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.VerifyToken())
	{
		protected.GET("/profile", GetProfile)
	}
}
