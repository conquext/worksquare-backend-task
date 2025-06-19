package logging

import (
	"time"

	"housing-api/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// RequestLogger logs HTTP requests
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request
		duration := time.Since(start)
		logger.Info("HTTP Request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration", duration.String(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
		)

		return err
	}
}
