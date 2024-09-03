package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User represents a  user for token generation
type User struct {
	ID       string
	Username string
	Email    string
}

// GenerateJWT creates a JWT token with  user data
func GenerateJWT(secretKey []byte) (string, error) {
	// Create  user data
	user := User{
		ID:       "123456",
		Username: "user",
		Email:    "@example.com",
	}

	// Create claims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Minute * time.Duration(mustAtoi(os.Getenv("JWT_EXPIRATION_MINUTES")))).Unix(), // Token expires in 24 hours
		"iat":      time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// mustAtoi converts a string to an integer and panics if it fails
func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
