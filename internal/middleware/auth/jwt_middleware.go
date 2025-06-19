package auth

import (
	"strings"

	"housing-api/internal/config"
	"housing-api/internal/services"
	"housing-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware validates JWT tokens
func JWTMiddleware(cfg *config.Config) fiber.Handler {
	authService := services.NewAuthService(cfg)

	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Authorization header is required", nil)
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Unauthorized(c, "Invalid authorization header format", nil)
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return response.Unauthorized(c, "Token is required", nil)
		}

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			return response.Unauthorized(c, "Invalid token", err)
		}

		// Set user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)

		return c.Next()
	}
}
