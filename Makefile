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

# Tests
test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Docker
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 3000:3000 $(DOCKER_IMAGE)

docker-compose:
	@echo "Starting with Docker Compose..."
	@docker-compose up -d

docker-compose-down:
	@echo "Stopping Docker Compose..."
	@docker-compose down

# Swagger
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g main.go

# Clean
clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html

# Linting
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
