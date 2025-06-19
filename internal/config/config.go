package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server Configuration
	Environment string
	Port        string
	Host        string

	// JWT Configuration
	JWTSecret           string
	JWTExpiresIn        time.Duration
	JWTRefreshExpiresIn time.Duration

	// Rate Limiting
	RateLimitWindowMS   time.Duration
	RateLimitMaxRequests int

	// Logging
	LogLevel string

	// API Configuration
	APIVersion string
	APIPrefix  string

	// Demo User Credentials
	DemoUserEmail    string
	DemoUserPassword string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Environment:          getEnv("NODE_ENV", "development"),
		Port:                 getEnv("PORT", "3000"),
		Host:                 getEnv("HOST", "localhost"),
		JWTSecret:           getEnv("JWT_SECRET", "super-secret-jwt-key"),
		JWTExpiresIn:        parseDuration(getEnv("JWT_EXPIRES_IN", "24h")),
		JWTRefreshExpiresIn: parseDuration(getEnv("JWT_REFRESH_EXPIRES_IN", "168h")), // 7 days
		RateLimitWindowMS:   parseDuration(getEnv("RATE_LIMIT_WINDOW_MS", "3600000ms")), // 1 hour
		RateLimitMaxRequests: parseInt(getEnv("RATE_LIMIT_MAX_REQUESTS", "100")),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		APIVersion:          getEnv("API_VERSION", "v1"),
		APIPrefix:           getEnv("API_PREFIX", "/api"),
		DemoUserEmail:       getEnv("DEMO_USER_EMAIL", "demo@worksquare.com"),
		DemoUserPassword:    getEnv("DEMO_USER_PASSWORD", "demo123456"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}