package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()

	// Set cache expiry from environment variable or use default
	cacheExpiryStr := os.Getenv("TOKEN_CACHE_EXPIRY")
	fmt.Printf("Cache expiry: %v\n", cacheExpiryStr)
	if cacheExpiryStr == "" {
		cacheExpiry = 5 * time.Minute
	} else {
		duration, err := time.ParseDuration(cacheExpiryStr)
		if err != nil {
			logger.Warn("Invalid TOKEN_CACHE_EXPIRY, using default of 5 minutes", zap.Error(err))
			cacheExpiry = 5 * time.Minute
		} else {
			cacheExpiry = duration
		}
	}
}

type Profile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type cacheEntry struct {
	profile Profile
	expiry  time.Time
}

var (
	tokenCache  = make(map[string]cacheEntry)
	cacheMutex  sync.RWMutex
	cacheExpiry time.Duration
)

func VerifyToken(customCacheExpiry ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyCacheExpiry := cacheExpiry
		if len(customCacheExpiry) > 0 {
			verifyCacheExpiryParsed, err := time.ParseDuration(customCacheExpiry[0])
			if err != nil {
				logger.Warn("Invalid cache expiry, using default", zap.Error(err))
			} else {
				verifyCacheExpiry = verifyCacheExpiryParsed
			}
		}
		// Get the token from the context set by AuthMiddleware
		token, exists := c.Get("auth_token")
		authHeader, authHeaderExists := c.Get("auth_header")

		if !authHeaderExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication header is missing"})
			c.Abort()
			return
		}
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token is missing"})
			c.Abort()
			return
		}

		tokenString, ok := token.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenURL := os.Getenv("TOKEN_URL")

		if tokenURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "TOKEN_URL not set"})
			c.Abort()
			return
		}

		// Check cache first
		cacheMutex.RLock()
		entry, found := tokenCache[tokenString]
		cacheMutex.RUnlock()

		// print details of the entry and expiry times as well as the current time
		fmt.Printf("Entry: %+v\n", entry)
		fmt.Printf("Expiry: %v\n", entry.expiry)
		fmt.Printf("Current time: %v\n", time.Now())

		valid := time.Now().Before(entry.expiry)
		fmt.Printf("Valid: %v\n", valid)

		if found && valid {
			fmt.Println("Token is still valid in cache")
			// Token is still valid in cache
			c.Header("X-Token-Cache", "HIT")
			c.Set("user", entry.profile)
			c.Next()
			return
		}

		// Token not found in cache or expired
		c.Header("X-Token-Cache", "MISS")

		// Validate token using the TOKEN_URL
		req, err := http.NewRequest("GET", tokenURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			c.Abort()
			return
		}
		req.Header.Set(authHeader.(string), tokenString)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// use logger to log the error
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			c.Abort()
			return
		}

		var profile Profile
		if err := json.Unmarshal(body, &profile); err != nil {
			// Use a custom unmarshaling to map the fields correctly
			var githubProfile struct {
				ID    uint   `json:"id"`
				Email string `json:"email"`
				Name  string `json:"name"`
				Login string `json:"login"`
			}
			if err := json.Unmarshal(body, &githubProfile); err != nil {
				logger.Error("Failed to parse GitHub profile", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse profile"})
				c.Abort()
				return
			}
			profile = Profile{
				ID:    fmt.Sprintf("%d", githubProfile.ID),
				Email: githubProfile.Email,
				Name:  githubProfile.Name,
			}
			// If Email is empty, use Login as a fallback
			if profile.Email == "" {
				profile.Email = githubProfile.Login
			}
		}

		// Cache the result
		cacheMutex.Lock()
		tokenCache[tokenString] = cacheEntry{
			profile: profile,
			expiry:  time.Now().Add(verifyCacheExpiry),
		}
		cacheMutex.Unlock()

		c.Set("user", profile)
		c.Next()
	}
}
