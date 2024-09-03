package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	zap "go.uber.org/zap"
	"golang.org/x/time/rate"
)

func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen int64 //lint:ignore U1000 This field is currently unused but may be used in future implementations
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	return func(c *gin.Context) {
		// Use only IP if Authorization header is empty
		key := c.ClientIP()
		if auth := c.GetHeader("Authorization"); auth != "" {
			key += ":" + auth
		}

		mu.Lock()
		if _, found := clients[key]; !found {
			clients[key] = &client{limiter: rate.NewLimiter(r, b)}
		}
		if !clients[key].limiter.Allow() {
			mu.Unlock()
			// Log the rate limit exceeded event
			authHeader := c.GetHeader("Authorization")
			maskedAuth := "No Authorization"
			if authHeader != "" {
				maskedAuth = "Bearer *****"
				if len(authHeader) > 15 {
					maskedAuth = authHeader[:15] + "*****"
				}
			}
			logger.Info("Rate limit exceeded",
				zap.String("client_ip", c.ClientIP()),
				zap.String("authorization", maskedAuth))
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		mu.Unlock()
		c.Next()
	}
}
