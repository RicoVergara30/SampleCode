package jwttoken

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"sample/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const (
	RandomDataLen = 32
)

// GetTokenKey retrieves the TOKEN_KEY environment variable.
func GetTokenKey() (string, error) {
	tokenKey := os.Getenv("TOKEN_KEY")
	if tokenKey == "" {
		return "", fmt.Errorf("TOKEN_KEY environment variable is not set")
	}
	return tokenKey, nil
}

// GenerateTokens generates a JWT token for the given username.
func GenerateTokens(username string, expirationDuration time.Duration) (string, error) {
	randomData := make([]byte, RandomDataLen)
	if _, err := rand.Read(randomData); err != nil {
		return "", err
	}
	randomString := base64.URLEncoding.EncodeToString(randomData)

	claims := jwt.MapClaims{
		"username":   username,
		"randomData": randomString,
		"exp":        time.Now().Add(expirationDuration).Unix(),
	}

	tokenKey, err := GetTokenKey()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates the JWT token.
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	tokenKey, err := GetTokenKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// TokenAuthMiddleware is a middleware for authenticating JWT tokens.
func TokenAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ResponseWoModel{
			RetCode: "401",
			Message: "Missing Authorization header",
		})
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ResponseWoModel{
			RetCode: "401",
			Message: "Invalid Authorization header format",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

	// Validate and parse the token
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ResponseWoModel{
			RetCode: "401",
			Message: "Invalid or expired JWT token",
		})
	}

	// Get the username from the claims
	username, ok := claims["username"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ResponseWoModel{
			RetCode: "401",
			Message: "Invalid JWT token claims",
		})
	}

	// Optionally, set the username in the request context for later handlers
	c.Locals("username", username)

	return c.Next()
}
