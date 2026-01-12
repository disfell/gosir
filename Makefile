# Makefile for Gosir project

.PHONY: help build run test clean swagger docs \
	docker-build docker-deploy docker-deploy-full docker-start docker-stop docker-restart docker-logs docker-clean docker-status docker-help

# 默认端口配置
PORT ?= 1323

# Default target
help:
	@echo "========================================="
	@echo "  Gosir Makefile"
	@echo "========================================="
	@echo ""
	@echo "应用开发："
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make swagger      - Generate Swagger documentation"
	@echo "  make docs         - Same as swagger (alias)"
	@echo ""
	@echo "Docker 操作（使用统一的部署脚本）："
	@echo "  make docker-build            - Build Docker image"
	@echo "  make docker-deploy           - Deploy Docker containers (existing image)"
	@echo "  make docker-deploy-full      - Build and deploy (rebuild image)"
	@echo "  make docker-start            - Start Docker containers"
	@echo "  make docker-stop             - Stop Docker containers"
	@echo "  make docker-restart          - Restart Docker containers"
	@echo "  make docker-logs             - Show Docker logs"
	@echo "  make docker-clean            - Clean Docker unused resources (reclaim disk space)"
	@echo "  make docker-status           - Show Docker status"
	@echo "  make docker-help             - Show Docker deployment help"
	@echo ""
	@echo "端口参数："
	@echo "  PORT          应用端口 (默认: 1323)"
	@echo ""
	@echo "示例："
	@echo "  make docker-deploy              # 使用默认端口 1323 部署"
	@echo "  make docker-deploy PORT=8080    # 指定端口 8080 部署"
	@echo ""

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
	@swag init -g cmd/server/main.go -o docs

docker-deploy:
	@cd docker && ./deploy.sh deploy-with-build $(PORT)

docker-start:
	@cd docker && ./deploy.sh start $(PORT)

docker-stop:
	@cd docker && docker compose down

docker-restart:
	@cd docker && ./deploy.sh restart $(PORT)

docker-logs:
	@cd docker && ./deploy.sh logs

docker-clean:
	@cd docker && ./deploy.sh clean

docker-status:
	@cd docker && ./deploy.sh status

docker-help:
	@cd docker && ./deploy.sh help
