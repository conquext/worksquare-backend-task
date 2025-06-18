package routes

import (
	"housing-api/internal/config"
	"housing-api/internal/controllers"
	"housing-api/internal/middleware/ratelimit"

	"github.com/gofiber/fiber/v2"
)

// Setup configures all application routes
func Setup(app *fiber.App, cfg *config.Config) {
	// Apply rate limiting to all routes
	app.Use(ratelimit.RateLimiter(cfg))

	// API prefix
	api := app.Group(cfg.APIPrefix + "/" + cfg.APIVersion)

	// Demo endpoints
	demoRoutes := api.Group("/demo")
	demoRoutes.Get("/credentials", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Demo credentials for testing",
			"data": fiber.Map{
				"email":    cfg.DemoUserEmail,
				"password": cfg.DemoUserPassword,
				"note":     "Use these credentials to test authentication endpoints",
			},
		})
	})
}