package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"housing-api/api/routes"
	"housing-api/internal/config"
	"housing-api/internal/models"
)

func setupAuthTestApp() *fiber.App {
	// Set test environment variables
	os.Setenv("NODE_ENV", "test")
	os.Setenv("JWT_SECRET", "test-secret-key-for-auth-integration-tests")
	os.Setenv("DEMO_USER_EMAIL", "test-demo@worksquare.com")
	os.Setenv("DEMO_USER_PASSWORD", "testdemo123")

	app := fiber.New(fiber.Config{
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
			})
		},
	})

	cfg, _ := config.Load()
	routes.Setup(app, cfg)
	return app
}

func TestAuthFlow_Complete(t *testing.T) {
	app := setupAuthTestApp()

	t.Run("Get Demo Credentials", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/demo/credentials", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "Demo credentials")
	})

	t.Run("Login with Demo Credentials", func(t *testing.T) {
		loginData := models.LoginRequest{
			Email:    "test-demo@worksquare.com",
			Password: "testdemo123",
		}

		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Login successful", response.Message)

		// Verify response structure
		authData := response.Data.(map[string]interface{})
		assert.Contains(t, authData, "access_token")
		assert.Contains(t, authData, "refresh_token")
		assert.Contains(t, authData, "user")
		assert.Contains(t, authData, "expires_in")
	})
}

func TestLogin_Success(t *testing.T) {
	app := setupAuthTestApp()

	loginData := models.LoginRequest{
		Email:    "test-demo@worksquare.com",
		Password: "testdemo123",
	}

	jsonData, err := json.Marshal(loginData)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "Login successful", response.Message)
	assert.NotNil(t, response.Data)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	app := setupAuthTestApp()

	testCases := []struct {
		name     string
		email    string
		password string
		expected int
	}{
		{
			name:     "Wrong password",
			email:    "test-demo@worksquare.com",
			password: "wrongpassword",
			expected: http.StatusUnauthorized,
		},
		{
			name:     "Wrong email",
			email:    "wrong@email.com",
			password: "testdemo123",
			expected: http.StatusUnauthorized,
		},
		{
			name:     "Empty credentials",
			email:    "",
			password: "",
			expected: http.StatusUnprocessableEntity,
		},
		{
			name:     "Invalid email format",
			email:    "notanemail",
			password: "testdemo123",
			expected: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loginData := models.LoginRequest{
				Email:    tc.email,
				Password: tc.password,
			}

			jsonData, _ := json.Marshal(loginData)
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, resp.StatusCode)

			var response models.APIResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.False(t, response.Success)
		})
	}
}

func TestRegister_Success(t *testing.T) {
	app := setupAuthTestApp()

	registerData := models.RegisterRequest{
		Email:    "newuser@test.com",
		Password: "newpassword123",
	}

	jsonData, err := json.Marshal(registerData)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "Registration successful", response.Message)
	assert.NotNil(t, response.Data)

	// Verify auth response structure
	authData := response.Data.(map[string]interface{})
	assert.Contains(t, authData, "access_token")
	assert.Contains(t, authData, "refresh_token")
	assert.Contains(t, authData, "user")

	// Verify user data
	userData := authData["user"].(map[string]interface{})
	assert.Equal(t, "newuser@test.com", userData["email"])
	assert.NotContains(t, userData, "password") // Password should not be in response
}

func TestRegister_DuplicateEmail(t *testing.T) {
	app := setupAuthTestApp()

	// First registration
	registerData := models.RegisterRequest{
		Email:    "duplicate@test.com",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(registerData)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Second registration with same email
	req2 := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp2.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp2.Body).Decode(&response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "User already exists", response.Error.Message)
}

func TestRefreshToken_Success(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI environment due to timeout issues")
	}
	app := setupAuthTestApp()

	// First, login to get tokens
	loginData := models.LoginRequest{
		Email:    "test-demo@worksquare.com",
		Password: "testdemo123",
	}

	jsonData, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResponse models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)

	authData := loginResponse.Data.(map[string]interface{})
	refreshToken := authData["refresh_token"].(string)

	// Use refresh token to get new access token
	refreshData := map[string]string{
		"refresh_token": refreshToken,
	}

	jsonData, _ = json.Marshal(refreshData)
	req2 := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(jsonData))
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var refreshResponse models.APIResponse
	err = json.NewDecoder(resp2.Body).Decode(&refreshResponse)
	require.NoError(t, err)

	assert.True(t, refreshResponse.Success)
	assert.Equal(t, "Token refreshed successfully", refreshResponse.Message)

	// Verify new tokens are provided
	newAuthData := refreshResponse.Data.(map[string]interface{})
	assert.Contains(t, newAuthData, "access_token")
	assert.Contains(t, newAuthData, "refresh_token")
	assert.NotEqual(t, refreshToken, newAuthData["refresh_token"]) // Should be a new refresh token
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	app := setupAuthTestApp()

	refreshData := map[string]string{
		"refresh_token": "invalid.jwt.token",
	}

	jsonData, _ := json.Marshal(refreshData)
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Invalid refresh token", response.Error.Message)
}

func TestGetProfile_Success(t *testing.T) {
	app := setupAuthTestApp()

	// Login to get access token
	loginData := models.LoginRequest{
		Email:    "test-demo@worksquare.com",
		Password: "testdemo123",
	}

	jsonData, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResponse models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)

	authData := loginResponse.Data.(map[string]interface{})
	accessToken := authData["access_token"].(string)

	// Get profile with access token
	req2 := httptest.NewRequest("GET", "/api/v1/auth/profile", nil)
	req2.Header.Set("Authorization", "Bearer "+accessToken)

	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var profileResponse models.APIResponse
	err = json.NewDecoder(resp2.Body).Decode(&profileResponse)
	require.NoError(t, err)

	assert.True(t, profileResponse.Success)
	assert.Equal(t, "Profile retrieved successfully", profileResponse.Message)

	userData := profileResponse.Data.(map[string]interface{})
	assert.Equal(t, "test-demo@worksquare.com", userData["email"])
	assert.NotContains(t, userData, "password")
}

func TestGetProfile_Unauthorized(t *testing.T) {
	app := setupAuthTestApp()

	testCases := []struct {
		name   string
		header string
	}{
		{
			name:   "No authorization header",
			header: "",
		},
		{
			name:   "Invalid token format",
			header: "InvalidToken",
		},
		{
			name:   "Invalid Bearer token",
			header: "Bearer invalid.jwt.token",
		},
		{
			name:   "Missing Bearer prefix",
			header: "some.jwt.token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/auth/profile", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			var response models.APIResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.False(t, response.Success)
		})
	}
}

func TestLogout_Success(t *testing.T) {
	app := setupAuthTestApp()

	// Login to get access token
	loginData := models.LoginRequest{
		Email:    "test-demo@worksquare.com",
		Password: "testdemo123",
	}

	jsonData, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResponse models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)

	authData := loginResponse.Data.(map[string]interface{})
	accessToken := authData["access_token"].(string)

	// Logout with access token
	req2 := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	req2.Header.Set("Authorization", "Bearer "+accessToken)

	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var logoutResponse models.APIResponse
	err = json.NewDecoder(resp2.Body).Decode(&logoutResponse)
	require.NoError(t, err)

	assert.True(t, logoutResponse.Success)
	assert.Equal(t, "Logout successful", logoutResponse.Message)
}

func TestValidation_LoginRequest(t *testing.T) {
	app := setupAuthTestApp()

	testCases := []struct {
		name        string
		requestBody string
		expectedMsg string
	}{
		{
			name:        "Invalid JSON",
			requestBody: `{"email": "test@test.com", "password":}`,
			expectedMsg: "Invalid request body",
		},
		{
			name:        "Missing email",
			requestBody: `{"password": "password123"}`,
			expectedMsg: "Validation failed",
		},
		{
			name:        "Missing password",
			requestBody: `{"email": "test@test.com"}`,
			expectedMsg: "Validation failed",
		},
		{
			name:        "Short password",
			requestBody: `{"email": "test@test.com", "password": "123"}`,
			expectedMsg: "Validation failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.True(t, resp.StatusCode >= 400)

			var response models.APIResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.False(t, response.Success)
			assert.Contains(t, response.Error.Message, tc.expectedMsg)
		})
	}
}

func TestConcurrentAuth_Operations(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping concurrent test in CI due to parallel thread execution")
	}
	app := setupAuthTestApp()

	// Test concurrent login attempts
	concurrency := 10
	results := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			loginData := models.LoginRequest{
				Email:    "test-demo@worksquare.com",
				Password: "testdemo123",
			}

			jsonData, _ := json.Marshal(loginData)
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			results <- (err == nil && resp.StatusCode == http.StatusOK)
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrency; i++ {
		if <-results {
			successCount++
		}
	}

	// All concurrent requests should succeed
	assert.Equal(t, concurrency, successCount)
}