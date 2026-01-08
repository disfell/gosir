# Makefile for Gosir project

.PHONY: help build run test clean swagger docs docker-build docker-run docker-down

# Default target
help:
	@echo "Available targets:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make swagger     - Generate Swagger documentation"
	@echo "  make docs        - Same as swagger (alias)"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker containers"
	@echo "  make docker-down  - Stop Docker containers"

# Build the application
build:
	@echo "Building Gosir..."
	@go build -o bin/gosir cmd/server/main.go

# Run the application
run:
	@echo "Running Gosir..."
	@go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf data.db
	@rm -rf logs/

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs

# Alias for swagger
docs: swagger

# Docker operations
docker-build:
	@echo "Building Docker image..."
	@cd docker && docker compose build

docker-run:
	@echo "Starting Docker containers..."
	@cd docker && docker compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@cd docker && docker compose down
