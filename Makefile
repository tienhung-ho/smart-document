all: build

.PHONY: build test clean docker migrate

# Build all services
build:
	@echo "Building all services..."
	@go build -o build/gateway ./apps/gateway/main.go
	@go build -o build/auth-api ./apps/auth/api/main.go
	@go build -o build/auth-rpc ./apps/auth/rpc/main.go
	@go build -o build/document-api ./apps/document/api/main.go
	@go build -o build/document-rpc ./apps/document/rpc/main.go
	@go build -o build/collaboration-api ./apps/collaboration/api/main.go
	@go build -o build/collaboration-rpc ./apps/collaboration/rpc/main.go
	@go build -o build/workspace-api ./apps/workspace/api/main.go
	@go build -o build/workspace-rpc ./apps/workspace/rpc/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf build/
	@rm -f coverage.out

# Docker commands
docker-build:
	@echo "Building Docker images..."
	@docker-compose build

docker-up:
	@echo "Starting services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping services..."
	@docker-compose down

# Database migrations
migrate-up:
	@echo "Running migrations up..."
	@./scripts/migrate.sh up

migrate-down:
	@echo "Running migrations down..."
	@./scripts/migrate.sh down

# Development
dev-setup:
	@echo "Setting up development environment..."
	@./scripts/setup.sh

# Generate code
generate:
	@echo "Generating code..."
	@./scripts/generate.sh

# Lint
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy