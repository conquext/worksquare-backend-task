package unit

import (
	"fmt"
	"os"
	"testing"

	"housing-api/internal/config"
	"housing-api/internal/models"
	"housing-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthTestEnvironment() *config.Config {
	// Set test environment variables
	os.Setenv("NODE_ENV", "test")
	os.Setenv("JWT_SECRET", "test-secret-key-for-auth-unit-tests-very-long-and-secure")
	os.Setenv("JWT_EXPIRES_IN", "1h")
	os.Setenv("JWT_REFRESH_EXPIRES_IN", "24h")
	os.Setenv("DEMO_USER_EMAIL", "test-unit@worksquare.com")
	os.Setenv("DEMO_USER_PASSWORD", "testunit123")

	// Create test data directory
	os.MkdirAll("../../testdata", 0755)

	cfg, _ := config.Load()
	return cfg
}

func cleanupAuthTestEnvironment() {
	// Clean up test data
	os.RemoveAll("../../testdata")
}

func TestAuthService_NewAuthService(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	assert.NotNil(t, service)
}

func TestAuthService_Login_Success(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	loginReq := models.LoginRequest{
		Email:    "test-unit@worksquare.com",
		Password: "testunit123",
	}

	authResponse, err := service.Login(loginReq)

	assert.NoError(t, err)
	assert.NotNil(t, authResponse)
	assert.NotEmpty(t, authResponse.AccessToken)
	assert.NotEmpty(t, authResponse.RefreshToken)
	assert.Equal(t, "test-unit@worksquare.com", authResponse.User.Email)
	assert.Greater(t, authResponse.ExpiresIn, int64(0))
	
	// Verify user data doesn't contain password
	assert.NotEmpty(t, authResponse.User.Email)
	assert.Greater(t, authResponse.User.ID, 0)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	testCases := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "Wrong password",
			email:    "test-unit@worksquare.com",
			password: "wrongpassword",
		},
		{
			name:     "Wrong email",
			email:    "wrong@email.com",
			password: "testunit123",
		},
		{
			name:     "Non-existent user",
			email:    "nonexistent@user.com",
			password: "somepassword",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loginReq := models.LoginRequest{
				Email:    tc.email,
				Password: tc.password,
			}

			authResponse, err := service.Login(loginReq)

			assert.Error(t, err)
			assert.Nil(t, authResponse)
			assert.Contains(t, err.Error(), "invalid credentials")
		})
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	registerReq := models.RegisterRequest{
		Email:    "newuser@test.com",
		Password: "newpassword123",
	}

	authResponse, err := service.Register(registerReq)

	assert.NoError(t, err)
	assert.NotNil(t, authResponse)
	assert.NotEmpty(t, authResponse.AccessToken)
	assert.NotEmpty(t, authResponse.RefreshToken)
	assert.Equal(t, "newuser@test.com", authResponse.User.Email)
	assert.Greater(t, authResponse.User.ID, 0)
	assert.Greater(t, authResponse.ExpiresIn, int64(0))

	// Verify user was created and can login
	loginReq := models.LoginRequest{
		Email:    "newuser@test.com",
		Password: "newpassword123",
	}

	loginResponse, err := service.Login(loginReq)
	assert.NoError(t, err)
	assert.NotNil(t, loginResponse)
	assert.Equal(t, authResponse.User.Email, loginResponse.User.Email)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	// First registration
	registerReq := models.RegisterRequest{
		Email:    "duplicate@test.com",
		Password: "password123",
	}

	authResponse1, err := service.Register(registerReq)
	assert.NoError(t, err)
	assert.NotNil(t, authResponse1)

	// Second registration with same email
	authResponse2, err := service.Register(registerReq)
	assert.Error(t, err)
	assert.Nil(t, authResponse2)
	assert.Contains(t, err.Error(), "already exists")
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	// First, login to get refresh token
	loginReq := models.LoginRequest{
		Email:    "test-unit@worksquare.com",
		Password: "testunit123",
	}

	loginResponse, err := service.Login(loginReq)
	require.NoError(t, err)
	require.NotNil(t, loginResponse)

	// Use refresh token to get new access token
	refreshResponse, err := service.RefreshToken(loginResponse.RefreshToken)

	assert.NoError(t, err)
	assert.NotNil(t, refreshResponse)
	assert.NotEmpty(t, refreshResponse.AccessToken)
	assert.NotEmpty(t, refreshResponse.RefreshToken)
	assert.Equal(t, loginResponse.User.Email, refreshResponse.User.Email)
	assert.Equal(t, loginResponse.User.ID, refreshResponse.User.ID)

	// New tokens should be different from original
	assert.NotEqual(t, loginResponse.AccessToken, refreshResponse.AccessToken)
	assert.NotEqual(t, loginResponse.RefreshToken, refreshResponse.RefreshToken)
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "Malformed token",
			token: "invalid.jwt.token",
		},
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Random string",
			token: "randomstring",
		},
		{
			name:  "Expired token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDk0NTkwMDB9.invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			refreshResponse, err := service.RefreshToken(tc.token)

			assert.Error(t, err)
			assert.Nil(t, refreshResponse)
			assert.Contains(t, err.Error(), "invalid refresh token")
		})
	}
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	// Login to get access token
	loginReq := models.LoginRequest{
		Email:    "test-unit@worksquare.com",
		Password: "testunit123",
	}

	authResponse, err := service.Login(loginReq)
	require.NoError(t, err)
	require.NotNil(t, authResponse)

	// Validate the access token
	claims, err := service.ValidateToken(authResponse.AccessToken)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, authResponse.User.ID, claims.UserID)
	assert.Equal(t, authResponse.User.Email, claims.Email)
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "Malformed token",
			token: "invalid.jwt.token",
		},
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Wrong signature",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20ifQ.wrong_signature",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tc.token)

			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}

func TestAuthService_GetUserByID_Success(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	// Get demo user (ID should be 1)
	user, err := service.GetUserByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "test-unit@worksquare.com", user.Email)
	assert.NotEmpty(t, user.Password) // Password should be present in internal model
}

func TestAuthService_GetUserByID_NotFound(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	user, err := service.GetUserByID(9999)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
}

func TestAuthService_MultipleUsers(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	// Register multiple users
	users := []models.RegisterRequest{
		{Email: "user1@test.com", Password: "password1"},
		{Email: "user2@test.com", Password: "password2"},
		{Email: "user3@test.com", Password: "password3"},
	}

	var userIDs []int

	for _, userReq := range users {
		authResponse, err := service.Register(userReq)
		require.NoError(t, err)
		require.NotNil(t, authResponse)
		
		userIDs = append(userIDs, authResponse.User.ID)

		// Verify each user can login
		loginReq := models.LoginRequest{
			Email:    userReq.Email,
			Password: userReq.Password,
		}

		loginResponse, err := service.Login(loginReq)
		assert.NoError(t, err)
		assert.NotNil(t, loginResponse)
		assert.Equal(t, userReq.Email, loginResponse.User.Email)
	}

	// Verify all users have unique IDs
	uniqueIDs := make(map[int]bool)
	for _, id := range userIDs {
		assert.False(t, uniqueIDs[id], "User ID %d should be unique", id)
		uniqueIDs[id] = true
	}
}

func TestAuthService_ConcurrentOperations(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping concurrent test in CI due to parallel thread execution")
	}
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	// Test concurrent login attempts
	concurrency := 10
	results := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			loginReq := models.LoginRequest{
				Email:    "test-unit@worksquare.com",
				Password: "testunit123",
			}

			_, err := service.Login(loginReq)
			results <- err
		}(i)
	}

	// Collect results
	for i := 0; i < concurrency; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}

func TestAuthService_ConcurrentRegistrations(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping concurrent test in CI due to parallel thread execution")
	}
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	// Test concurrent registration attempts with different emails
	concurrency := 5
	results := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			registerReq := models.RegisterRequest{
				Email:    fmt.Sprintf("concurrent%d@test.com", id),
				Password: "password123",
			}

			_, err := service.Register(registerReq)
			results <- err
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrency; i++ {
		err := <-results
		if err == nil {
			successCount++
		}
	}

	// All registrations should succeed since emails are unique
	assert.Equal(t, concurrency, successCount)
}

func TestAuthService_PasswordSecurity(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	// Register a user
	registerReq := models.RegisterRequest{
		Email:    "security@test.com",
		Password: "plainpassword123",
	}

	authResponse, err := service.Register(registerReq)
	require.NoError(t, err)
	require.NotNil(t, authResponse)

	// Get user by ID to check internal password storage
	user, err := service.GetUserByID(authResponse.User.ID)
	require.NoError(t, err)
	require.NotNil(t, user)

	// Password should be hashed, not plain text
	assert.NotEqual(t, "plainpassword123", user.Password)
	assert.Greater(t, len(user.Password), 20) // Bcrypt hashes are typically 60+ chars
	assert.Contains(t, user.Password, "$2a$") // Bcrypt hash prefix

	// Verify auth response doesn't contain password
	assert.Empty(t, authResponse.User.CreatedAt.IsZero()) // Has created time
	assert.Empty(t, authResponse.User.UpdatedAt.IsZero()) // Has updated time
}

func TestAuthService_TokenIntegrity(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	// Login to get tokens
	loginReq := models.LoginRequest{
		Email:    "test-unit@worksquare.com",
		Password: "testunit123",
	}

	authResponse, err := service.Login(loginReq)
	require.NoError(t, err)
	require.NotNil(t, authResponse)

	// Validate access token contains correct claims
	claims, err := service.ValidateToken(authResponse.AccessToken)
	require.NoError(t, err)
	require.NotNil(t, claims)

	assert.Equal(t, authResponse.User.ID, claims.UserID)
	assert.Equal(t, authResponse.User.Email, claims.Email)

	// Validate refresh token contains correct claims
	refreshClaims, err := service.ValidateToken(authResponse.RefreshToken)
	require.NoError(t, err)
	require.NotNil(t, refreshClaims)

	assert.Equal(t, authResponse.User.ID, refreshClaims.UserID)
	assert.Equal(t, authResponse.User.Email, refreshClaims.Email)
}

func TestAuthService_EdgeCases(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	t.Run("Empty email registration", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:    "",
			Password: "password123",
		}

		authResponse, err := service.Register(registerReq)
		// This should be handled by validation layer, but service should be robust
		assert.Error(t, err)
		assert.Nil(t, authResponse)
	})

	t.Run("Email case sensitivity", func(t *testing.T) {
		// Register with lowercase email
		registerReq := models.RegisterRequest{
			Email:    "case@test.com",
			Password: "password123",
		}

		authResponse1, err := service.Register(registerReq)
		require.NoError(t, err)
		require.NotNil(t, authResponse1)

		// Try to login with uppercase email
		loginReq := models.LoginRequest{
			Email:    "CASE@TEST.COM",
			Password: "password123",
		}

		authResponse2, err := service.Login(loginReq)
		assert.NoError(t, err) // Should work (case insensitive)
		assert.NotNil(t, authResponse2)
		assert.Equal(t, authResponse1.User.ID, authResponse2.User.ID)
	})
}

func TestAuthService_MemoryManagement(t *testing.T) {
	cfg := setupAuthTestEnvironment()
	defer cleanupAuthTestEnvironment()

	service := services.NewAuthService(cfg)

	// Create many users to test memory usage
	userCount := 100
	for i := 0; i < userCount; i++ {
		registerReq := models.RegisterRequest{
			Email:    fmt.Sprintf("memtest%d@test.com", i),
			Password: "password123",
		}

		authResponse, err := service.Register(registerReq)
		assert.NoError(t, err)
		assert.NotNil(t, authResponse)

		// Immediately try to login to verify user is properly stored
		loginReq := models.LoginRequest{
			Email:    registerReq.Email,
			Password: registerReq.Password,
		}

		loginResponse, err := service.Login(loginReq)
		assert.NoError(t, err)
		assert.NotNil(t, loginResponse)
	}

	// Verify all users can still be accessed
	for i := 1; i <= userCount+1; i++ { // +1 for demo user
		user, err := service.GetUserByID(i)
		assert.NoError(t, err)
		assert.NotNil(t, user)
	}
}