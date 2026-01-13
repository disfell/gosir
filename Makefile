# Makefile for Gosir project

.PHONY: build run test clean swagger docs \
	docker-deploy docker-stop docker-restart docker-clean docker-help

# 默认端口配置
PORT ?= 1323

# Build the application
build: swagger
	@echo "Building Gosir..."
	@go build -o bin/gosir cmd/server/main.go
	@echo "Done!"

# Run the application
run: swagger
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
	@echo "Done!"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go --parseDependency --parseInternal -o docs

docker-deploy:
	@cd docker && ./deploy.sh build_and_deploy $(PORT)

docker-stop:
	@cd docker && docker compose down

docker-restart:
	@cd docker && ./deploy.sh restart $(PORT)

docker-clean:
	@cd docker && ./deploy.sh clean

docker-help:
	@cd docker && ./deploy.sh help
