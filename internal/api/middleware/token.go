package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

type Profile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		tokenURL := os.Getenv("TOKEN_URL")

		fmt.Println("tokenURL", tokenURL)
		fmt.Println("token", token)
		if tokenURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "TOKEN_URL not set"})
			c.Abort()
			return
		}

		// Validate token using the TOKEN_URL
		req, err := http.NewRequest("GET", tokenURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			c.Abort()
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)

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

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			c.Abort()
			return
		}

		fmt.Println("body", string(body))

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

		c.Set("user", profile)
		c.Next()
	}
}
