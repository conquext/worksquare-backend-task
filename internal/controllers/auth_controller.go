package controllers

import (
	"housing-api/internal/config"
	"housing-api/internal/models"
	"housing-api/internal/services"
	"housing-api/internal/utils"
	"housing-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		authService: services.NewAuthService(cfg),
	}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req models.LoginRequest

	// Parse request body
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, "Invalid request body", err)
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return response.ValidationError(ctx, "Validation failed", err)
	}

	// Authenticate user
	authResponse, err := c.authService.Login(req)
	if err != nil {
		return response.Unauthorized(ctx, "Authentication failed", err)
	}

	return response.Success(ctx, "Login successful", authResponse)
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req models.RegisterRequest

	// Parse request body
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, "Invalid request body", err)
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return response.ValidationError(ctx, "Validation failed", err)
	}

	// Register user
	authResponse, err := c.authService.Register(req)
	if err != nil {
		if err.Error() == "user with email "+req.Email+" already exists" {
			return response.Conflict(ctx, "User already exists", err)
		}
		return response.InternalServerError(ctx, "Registration failed", err)
	}

	return response.Created(ctx, "Registration successful", authResponse)
}

func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	var req map[string]string

	// Parse request body
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, "Invalid request body", err)
	}

	refreshToken, exists := req["refresh_token"]
	if !exists || refreshToken == "" {
		return response.BadRequest(ctx, "Refresh token is required", nil)
	}

	// Refresh token
	authResponse, err := c.authService.RefreshToken(refreshToken)
	if err != nil {
		return response.Unauthorized(ctx, "Invalid refresh token", err)
	}

	return response.Success(ctx, "Token refreshed successfully", authResponse)
}

func (c *AuthController) GetProfile(ctx *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := ctx.Locals("userID").(int)

	// Get user
	user, err := c.authService.GetUserByID(userID)
	if err != nil {
		return response.NotFound(ctx, "User not found", err)
	}

	return response.Success(ctx, "Profile retrieved successfully", user.ToUserResponse())
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	return response.Success(ctx, "Logout successful", map[string]string{
		"message": "",
	})
}