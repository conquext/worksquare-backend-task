package ratelimit

import (
	"housing-api/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter creates a rate limiting middleware
func RateLimiter(cfg *config.Config) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cfg.RateLimitMaxRequests,
		Expiration: cfg.RateLimitWindowMS,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    fiber.StatusTooManyRequests,
					"message": "Rate limit exceeded",
				},
				"data": nil,
			})
		},
	})
}
