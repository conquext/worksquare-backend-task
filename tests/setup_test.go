package tests

import (
	"os"
	"testing"

	"housing-api/internal/config"
)

var testConfig *config.Config

func TestMain(m *testing.M) {
	// Set test environment
	os.Setenv("NODE_ENV", "test")
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("DEMO_USER_EMAIL", "test@example.com")
	os.Setenv("DEMO_USER_PASSWORD", "testpassword")

	// Load test configuration
	var err error
	testConfig, err = config.Load()
	if err != nil {
		panic("Failed to load test config: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Exit(code)
}
