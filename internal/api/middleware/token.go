package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	defaultCacheExpiry = 5 * time.Minute
	schemeOffset       = 3
)

// TokenValidator handles token validation with caching.
type TokenValidator struct {
	logger      *zap.Logger
	tokenCache  map[string]cacheEntry
	cacheMutex  sync.RWMutex
	cacheExpiry time.Duration
}

// NewTokenValidator creates a new TokenValidator instance.
func NewTokenValidator(logger *zap.Logger) *TokenValidator {
	cacheExpiryStr := os.Getenv("TOKEN_CACHE_EXPIRY")
	var cacheExpiry time.Duration
	if cacheExpiryStr == "" {
		cacheExpiry = defaultCacheExpiry
	} else {
		duration, err := time.ParseDuration(cacheExpiryStr)
		if err != nil {
			logger.Warn("Invalid TOKEN_CACHE_EXPIRY, using default of 5 minutes", zap.Error(err))
			cacheExpiry = defaultCacheExpiry
		} else {
			cacheExpiry = duration
		}
	}

	return &TokenValidator{
		logger:      logger,
		tokenCache:  make(map[string]cacheEntry),
		cacheExpiry: cacheExpiry,
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

// Global default token validator for backward compatibility.
var defaultTokenValidator *TokenValidator

// VerifyToken returns a gin middleware for token verification (backward compatibility).
func VerifyToken(customCacheExpiry ...string) gin.HandlerFunc {
	if defaultTokenValidator == nil {
		logger, _ := zap.NewProduction()
		defaultTokenValidator = NewTokenValidator(logger)
	}
	return defaultTokenValidator.VerifyToken(customCacheExpiry...)
}

// VerifyToken returns a gin middleware for token verification.
func (tv *TokenValidator) VerifyToken(customCacheExpiry ...string) gin.HandlerFunc {
	// Initialize tokenURL if needed.
	initTokenURL()

	return func(c *gin.Context) {
		verifyCacheExpiry := tv.cacheExpiry

		if len(customCacheExpiry) > 0 {
			verifyCacheExpiryParsed, err := time.ParseDuration(customCacheExpiry[0])
			if err != nil {
				tv.logger.Warn("Invalid cache expiry, using default", zap.Error(err))
			} else {
				tv.logger.Info("Using custom cache expiry", zap.Duration("expiry", verifyCacheExpiryParsed))
				verifyCacheExpiry = verifyCacheExpiryParsed
				tv.logger.Info("verifyCacheExpiry", zap.Duration("expiry", verifyCacheExpiry))
			}
		}
		// Get the token from the context set by AuthMiddleware.
		token, exists := c.Get("auth_token")
		authHeader, authHeaderExists := c.Get("auth_header")

		if !authHeaderExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication header is missing"})
			c.Abort()

			return
		}

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()

			return
		}

		tokenString, ok := token.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()

			return
		}

		// Check cache first.
		tv.cacheMutex.RLock()
		entry, found := tv.tokenCache[tokenString]
		tv.cacheMutex.RUnlock()

		valid := time.Now().Before(entry.expiry)
		tv.logger.Info("Token validation result", zap.Bool("valid", valid))

		if found && valid {
			tv.logger.Info("Token is still valid in cache")
			// Token is still valid in cache.
			c.Header("X-Token-Cache", "HIT")
			c.Set("user", entry.profile)
			c.Next()

			return
		}

		// Token not found in cache or expired.
		c.Header("X-Token-Cache", "MISS")

		// Validate token using the TOKEN_URL.
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, tokenURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			c.Abort()

			return
		}

		authHeaderStr, ok := authHeader.(string)
		if !ok {
			tv.logger.Error("authHeader is not a string")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}
		req.Header.Set(authHeaderStr, tokenString)

		client := &http.Client{} //nolint:exhaustruct // Default HTTP client is sufficient.

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
			c.Abort()

			return
		}

		defer func() {
			closeErr := resp.Body.Close()
			if closeErr != nil {
				tv.logger.Error("Failed to close response body", zap.Error(closeErr))
			}
		}()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()

			return
		}

		// Read and parse the response.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			tv.logger.Error("Failed to read response body", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read profile"})
			c.Abort()

			return
		}

		var profile Profile
		if profileErr := json.Unmarshal(body, &profile); profileErr != nil {
			tv.logger.Error("Failed to unmarshal profile", zap.Error(profileErr))

			// Try parsing as a GitHub profile.
			var githubProfile map[string]interface{} //nolint:exhaustruct // Map doesn't need field names.
			if githubErr := json.Unmarshal(body, &githubProfile); githubErr != nil {
				tv.logger.Error("Failed to unmarshal GitHub profile", zap.Error(githubErr))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse profile"})
				c.Abort()

				return
			}

			// Convert GitHub profile to our Profile struct.
			profile = Profile{
				ID:    fmt.Sprintf("%v", githubProfile["id"]),
				Email: fmt.Sprintf("%v", githubProfile["email"]),
				Name:  fmt.Sprintf("%v", githubProfile["name"]),
			}
		}

		// Store in cache.
		tv.cacheMutex.Lock()
		tv.tokenCache[tokenString] = cacheEntry{
			profile: profile,
			expiry:  time.Now().Add(verifyCacheExpiry),
		}
		tv.cacheMutex.Unlock()

		// Set the user profile in the context.
		c.Set("user", profile)
		c.Next()
	}
}

// GetUserFromContext retrieves the user profile from the Gin context.
func GetUserFromContext(c *gin.Context) (Profile, bool) {
	user, exists := c.Get("user")
	if !exists {
		return Profile{}, false //nolint:exhaustruct // Empty struct is intentional.
	}

	profile, ok := user.(Profile)

	return profile, ok
}

var tokenURL = os.Getenv("TOKEN_URL")

func getUserInfoURL(provider string) string {
	userInfoURLs := map[string]string{
		"github": "https://api.github.com/user",
		"google": "https://www.googleapis.com/oauth2/v1/userinfo",
		"okta":   "/oauth2/v1/userinfo", // This will be appended to the issuer URL.
	}

	if url, ok := userInfoURLs[provider]; ok {
		return url
	}

	return ""
}

// initTokenURL initializes the tokenURL if not set.
func initTokenURL() {
	if tokenURL == "" {
		// Check if we're using OIDC.
		issuer := os.Getenv("OIDC_ISSUER")
		if issuer != "" {
			// Determine provider from issuer.
			var provider string
			switch {
			case contains(issuer, "github"):
				provider = "github"
			case contains(issuer, "google"):
				provider = "google"
			case contains(issuer, "okta"):
				provider = "okta"
			default:
				// Default to OIDC userinfo endpoint.
				tokenURL = issuer + "/oauth2/v1/userinfo"

				return
			}

			// Get the appropriate user info URL.
			userInfoURL := getUserInfoURL(provider)
			if provider == "okta" {
				tokenURL = issuer + userInfoURL
			} else {
				tokenURL = userInfoURL
			}
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}

	if len(substr) > len(s) {
		return false
	}

	// Check each possible starting position.
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := range len(substr) {
			if s[i+j] != substr[j] {
				match = false

				break
			}
		}

		if match {
			return true
		}
	}

	return false
}

// EnsureTokenURLisHTTP ensures the token URL starts with http/https.
func EnsureTokenURLisHTTP() {
	if tokenURL != "" && !isHTTP(tokenURL) && !isHTTPS(tokenURL) {
		if isLocalhost(tokenURL) {
			tokenURL = "http://" + tokenURL
		} else {
			tokenURL = "https://" + tokenURL
		}
	}
}

func isHTTP(url string) bool {
	return len(url) >= 7 && url[:7] == "http://"
}

func isHTTPS(url string) bool {
	return len(url) >= 8 && url[:8] == "https://"
}

func isLocalhost(url string) bool {
	return contains(url, "localhost") || contains(url, "127.0.0.1") || contains(url, "host.docker.internal")
}

// isPortNumber checks if a string contains only digits (for port validation).
func isPortNumber(s string) bool {
	if s == "" {
		return false
	}

	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// ExtractPort gets the port from a URL.
func ExtractPort(url string) string {
	// Find the last colon.
	lastColon := -1
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == ':' {
			lastColon = i

			break
		}
	}

	if lastColon == -1 {
		return ""
	}

	// Extract potential port.
	potentialPort := ""
	for i := lastColon + 1; i < len(url); i++ {
		if url[i] == '/' {
			break
		}

		potentialPort += string(url[i])
	}

	// Validate it's a port.
	if isPortNumber(potentialPort) {
		return potentialPort
	}

	return ""
}

// GetBaseURL extracts base URL without path.
func GetBaseURL(url string) string {
	// Find where the path starts (first / after ://).
	schemeEnd := -1
	for i := range len(url) - 2 {
		if url[i] == ':' && url[i+1] == '/' && url[i+2] == '/' {
			schemeEnd = i + schemeOffset

			break
		}
	}

	if schemeEnd == -1 {
		schemeEnd = 0
	}

	// Find the first / after the scheme.
	for i := schemeEnd; i < len(url); i++ {
		if url[i] == '/' {
			return url[:i]
		}
	}

	return url
}

// EnsureTokenURLHasPort ensures local URLs have a port.
func EnsureTokenURLHasPort() {
	if tokenURL == "" {
		return
	}

	baseURL := GetBaseURL(tokenURL)
	port := ExtractPort(baseURL)

	// If it's a local URL without a port, add default port.
	if isLocalhost(baseURL) && port == "" {
		// Find where to insert the port.
		if isHTTP(tokenURL) {
			tokenURL = "http://localhost:3000" + tokenURL[len("http://localhost"):]
		} else if isHTTPS(tokenURL) {
			tokenURL = "https://localhost:3001" + tokenURL[len("https://localhost"):]
		}
	}
}

// GetMaxSizeBytesFromEnv gets the max size from environment variable.
func GetMaxSizeBytesFromEnv() int64 {
	const (
		defaultMaxSize = 10 * 1024 * 1024  // 10MB default
		minSize        = 1024              // 1KB minimum
		maxSize        = 100 * 1024 * 1024 // 100MB maximum
	)

	maxSizeStr := os.Getenv("MAX_BODY_SIZE")
	if maxSizeStr == "" {
		return defaultMaxSize
	}

	maxSizeInt, err := strconv.ParseInt(maxSizeStr, 10, 64)
	if err != nil {
		return defaultMaxSize
	}

	if maxSizeInt < minSize {
		return minSize
	}

	if maxSizeInt > maxSize {
		return maxSize
	}

	return maxSizeInt
}
