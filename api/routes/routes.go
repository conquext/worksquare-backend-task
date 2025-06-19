package routes

import (
	"housing-api/internal/config"
	"housing-api/internal/controllers"
	"housing-api/internal/middleware/auth"
	"housing-api/internal/middleware/ratelimit"

	"github.com/gofiber/fiber/v2"
)

// Setup configures all application routes
func Setup(app *fiber.App, cfg *config.Config) {
	// Apply rate limiting to all routes
	app.Use(ratelimit.RateLimiter(cfg))

	// API prefix
	api := app.Group(cfg.APIPrefix + "/" + cfg.APIVersion)

	// Initialize controllers
	listingController, err := controllers.NewListingController()
	if err != nil {
		panic("Failed to initialize listing controller: " + err.Error())
	}

	authController := controllers.NewAuthController(cfg)

	// Auth routes (public)
	authRoutes := api.Group("/auth")
	authRoutes.Post("/login", authController.Login)
	authRoutes.Post("/register", authController.Register)
	authRoutes.Post("/refresh", authController.RefreshToken)

	// Protected auth routes
	authRoutes.Get("/profile", auth.JWTMiddleware(cfg), authController.GetProfile)
	authRoutes.Post("/logout", auth.JWTMiddleware(cfg), authController.Logout)

	// Listing routes (public)
	listingRoutes := api.Group("/listings")
	listingRoutes.Get("/", listingController.GetListings)
	listingRoutes.Get("/search", listingController.SearchListings)
	listingRoutes.Get("/filters", listingController.GetFiltersMetadata)
	listingRoutes.Get("/:id", listingController.GetListingByID)

	// Protected listing routes
	listingRoutes.Get("/stats", auth.JWTMiddleware(cfg), listingController.GetListingStats)

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