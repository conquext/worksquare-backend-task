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

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/login [post]
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

// Register godoc
// @Summary User registration
// @Description Register a new user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/register [post]
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

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token" example({"refresh_token": "your_refresh_token_here"})
// @Success 200 {object} models.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/refresh [post]
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

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile (protected route)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/profile [get]
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

// Logout godoc
// @Summary User logout
// @Description Logout user (client should discard tokens)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	// In a real implementation, we will blacklist the token
	// For now, we just return a success message
	return response.Success(ctx, "Logout successful", map[string]string{
		"message": "Please discard your tokens on the client side",
	})
}