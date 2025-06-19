.PHONY: build run test clean docker-build docker-run dev deps swagger

# Variables
APP_NAME=housing-api
BINARY_NAME=main
DOCKER_IMAGE=housing-api:latest

# Development
dev:
	@echo "Starting development server..."
	@go run main.go

# Build
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BINARY_NAME) .

# Run
run: build
	@echo "Running $(APP_NAME)..."
	@./$(BINARY_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
