.PHONY: help build run test clean docker-up docker-down docker-logs lint security

# Default target
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application locally"
	@echo "  test        - Run tests with coverage"
	@echo "  test-watch  - Run tests in watch mode"
	@echo "  lint        - Run linter"
	@echo "  security    - Run security scan"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-up   - Start services with Docker Compose"
	@echo "  docker-down - Stop Docker Compose services"
	@echo "  docker-logs - View Docker Compose logs"

# Build the application with version info
build:
	@echo "Building Go application..."
	@CGO_ENABLED=0 GOOS=linux go build \
		-ldflags="-X main.Version=$$(git describe --tags --always --dirty) -X main.BuildTime=$$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.GitCommit=$$(git rev-parse HEAD)" \
		-o bin/server cmd/server/main.go
	@echo "✅ Build completed"

# Run the application locally
run:
	@echo "Starting Go application..."
	@go run cmd/server/main.go

# Run tests with coverage
test:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Tests completed. Coverage report: coverage.html"

# Run tests in watch mode (requires entr)
test-watch:
	@echo "Running tests in watch mode..."
	@find . -name "*.go" | entr -c go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run --timeout=5m
	@echo "✅ Linting completed"

# Run security scan
security:
	@echo "Running security scan..."
	@gosec -fmt json -out gosec-report.json -stdout ./...
	@echo "✅ Security scan completed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ coverage.out coverage.html gosec-report.json
	@echo "✅ Clean completed"

# Docker Compose commands
docker-up:
	@echo "Starting services with Docker Compose..."
	@docker compose up -d
	@echo "✅ Services started"

docker-down:
	@echo "Stopping Docker Compose services..."
	@docker compose down
	@echo "✅ Services stopped"

docker-logs:
	@docker compose logs -f

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@go mod tidy
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "✅ Development environment setup completed"

# Format code
fmt:
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Build and push Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t go-auth:latest \
		--build-arg VERSION=$$(git describe --tags --always --dirty) \
		--build-arg BUILD_TIME=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
		--build-arg GIT_COMMIT=$$(git rev-parse HEAD) .
	@echo "✅ Docker image built"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@docker compose -f docker-compose.test.yml up --build --abort-on-container-exit
	@docker compose -f docker-compose.test.yml down
	@echo "✅ Integration tests completed"

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	@swag init -g cmd/server/main.go -o docs/
	@echo "✅ API documentation generated" 