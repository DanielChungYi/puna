.PHONY: help build run test clean docker-build docker-run docker-compose-up docker-compose-down

# Default target
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-compose-up    - Start all services with docker-compose"
	@echo "  docker-compose-down  - Stop all services"
	@echo "  install-deps   - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/puna cmd/web/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/web/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t puna .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 puna

# Start all services with docker-compose
docker-compose-up:
	@echo "Starting services with docker-compose..."
	docker-compose up -d

# Stop all services
docker-compose-down:
	@echo "Stopping services..."
	docker-compose down

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Build for different platforms
build-all: build-linux build-windows build-macos

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/puna-linux cmd/web/main.go

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o bin/puna-windows.exe cmd/web/main.go

build-macos:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o bin/puna-macos cmd/web/main.go

# Development helpers
dev:
	@echo "Starting development server..."
	air

# Database helpers
db-create:
	@echo "Creating database..."
	createdb puna_db

db-drop:
	@echo "Dropping database..."
	dropdb puna_db

db-reset: db-drop db-create
	@echo "Database reset complete"

# Security helpers
generate-secret:
	@echo "Generating secret key..."
	openssl rand -base64 32

# Release helpers
release:
	@echo "Creating release..."
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
	@echo "Release v$(VERSION) created"

# Help for release
release-help:
	@echo "Usage: make release VERSION=x.x.x"
	@echo "Example: make release VERSION=1.0.0"
