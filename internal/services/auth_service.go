package services

import (
	"fmt"
	"strings"
	"time"

	"housing-api/internal/config"
	"housing-api/internal/models"
	"housing-api/internal/utils"
	"housing-api/pkg/jwt"
)

// AuthService handles authentication business logic
type AuthService struct {
	config *config.Config
	users  []models.User // In-memory user store for demo
}

func NewAuthService(cfg *config.Config) *AuthService {
	service := &AuthService{
		config: cfg,
		users:  []models.User{},
	}

	service.createDemoUser()

	return service
}

func (s *AuthService) createDemoUser() {
	hashedPassword, err := utils.HashPassword(s.config.DemoUserPassword)
	if err != nil {
		return
	}

	demoUser := models.User{
		ID:        1,
		Email:     s.config.DemoUserEmail,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.users = append(s.users, demoUser)
}

// Login authenticates a user and returns JWT tokens
func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	
	// Find user by email
	user := s.findUserByEmail(email)
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT tokens
	accessToken, err := jwt.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.JWTExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.JWTRefreshExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.AuthResponse{
		User:         user.ToUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWTExpiresIn.Seconds()),
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	
	// Check if email is empty
	if email == "" {
		return nil, fmt.Errorf("email must not be empty")
	}
	
	// Check if user already exists
	if s.findUserByEmail(email) != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	newUser := models.User{
		ID:        len(s.users) + 1,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.users = append(s.users, newUser)

	// Generate JWT tokens
	accessToken, err := jwt.GenerateToken(newUser.ID, newUser.Email, s.config.JWTSecret, s.config.JWTExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateToken(newUser.ID, newUser.Email, s.config.JWTSecret, s.config.JWTRefreshExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.AuthResponse{
		User:         newUser.ToUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWTExpiresIn.Seconds()),
	}, nil
}

// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	// Validate refresh token
	claims, err := jwt.ValidateToken(refreshToken, s.config.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Find user
	user := s.findUserByID(claims.UserID)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new access token
	accessToken, err := jwt.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.JWTExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := jwt.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.JWTRefreshExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.AuthResponse{
		User:         user.ToUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.config.JWTExpiresIn.Seconds()),
	}, nil
}

// ValidateToken validates JWT token and returns user claims
func (s *AuthService) ValidateToken(token string) (*jwt.Claims, error) {
	return jwt.ValidateToken(token, s.config.JWTSecret)
}

// GetUserByID returns user by ID
func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	user := s.findUserByID(id)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// findUserByEmail finds user by email
func (s *AuthService) findUserByEmail(email string) *models.User {
	for _, user := range s.users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

// findUserByID finds user by ID
func (s *AuthService) findUserByID(id int) *models.User {
	for _, user := range s.users {
		if user.ID == id {
			return &user
		}
	}
	return nil
}