package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"housing-api/api/routes"
	"housing-api/internal/config"
	"housing-api/internal/middleware/logging"
	"housing-api/pkg/logger"
)

func Run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Worksquare Housing API",
		ServerHeader: "Fiber",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    code,
					"message": err.Error(),
				},
				"data": nil,
			})
		},
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(logging.RequestLogger())

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Worksquare Housing API is running!",
			"data": fiber.Map{
				"status":  "healthy",
				"version": "1.0.0",
			},
		})
	})

	// Setup routes
	routes.Setup(app, cfg)

	// Start server
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	log.Printf("🚀 Server started on port %s", cfg.Port)
	log.Printf("📚 Swagger documentation available at http://localhost:%s/swagger/", cfg.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	if err := app.Shutdown(); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("✅ Server exited")
	return nil
}