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

	// Logging
	LogLevel string

	// API Configuration
	APIVersion string
	APIPrefix  string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Environment:          getEnv("NODE_ENV", "development"),
		Port:                 getEnv("PORT", "3000"),
		Host:                 getEnv("HOST", "localhost"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		APIVersion:          getEnv("API_VERSION", "v1"),
		APIPrefix:           getEnv("API_PREFIX", "/api"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

