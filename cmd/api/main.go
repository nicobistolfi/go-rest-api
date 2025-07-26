package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicobistolfi/go-rest-api/internal/api"
	"github.com/nicobistolfi/go-rest-api/internal/config"
	logger "github.com/nicobistolfi/go-rest-api/pkg"
	"go.uber.org/zap"
)

const (
	readTimeout      = 5 * time.Second
	writeTimeout     = 10 * time.Second
	idleTimeout      = 120 * time.Second
	shutdownTimeout  = 5 * time.Second
)

func main() {
	logger.Init()

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create a new Gin router
	r := gin.Default()

	api.SetupRouter(r, cfg, logger.Log)

	// Create a new server with timeouts
	srv := &http.Server{ //nolint:exhaustruct // Optional fields not needed
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// Graceful shutdown
	go func() {
		if listenErr := srv.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			logger.Fatal("listen: %s\n", zap.Error(listenErr))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(shutdownErr))
	}

	logger.Info("Server exiting")
}
