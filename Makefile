.PHONY: build run clean test lint migrate-up migrate-down migrate-create sqlc docker-up docker-down

# Binary output
BIN_DIR = bin
BIN_NAME = api

# Build settings
GO = go
GOFLAGS = -v
BUILD_DIR = ./cmd/api

# Database settings
DB_USER = postgres
DB_PASSWORD = postgres
DB_NAME = yadwy
DB_HOST = localhost
DB_PORT = 5432
MIGRATION_DIR = migrations

# Default target
all: build

# Build the application
build:
	@echo "Building application..."
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(BIN_NAME) $(BUILD_DIR)

# Run the application
run: build
	@echo "Running application..."
	./$(BIN_DIR)/$(BIN_NAME)

# Run the application in development mode with hot reload
dev:
	@echo "Running in development mode..."
	air -c .air.toml

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	$(GO) clean

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Create a new migration file
migrate-create:
	@echo "Creating migration file..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $${name}

# Run migrations up
migrate-up:
	@echo "Running migrations up..."
	migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

# Run migrations down
migrate-down:
	@echo "Running migrations down..."
	migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# Start PostgreSQL in Docker
docker-up:
	@echo "Starting PostgreSQL in Docker..."
	docker run --name yadwy-postgres -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_DB=$(DB_NAME) -p $(DB_PORT):5432 -d postgres:15

# Stop and remove PostgreSQL Docker container
docker-down:
	@echo "Stopping PostgreSQL Docker container..."
	docker stop yadwy-postgres || true
	docker rm yadwy-postgres || true

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Run with hot reload (requires air)"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  lint           - Run linter"
	@echo "  migrate-create - Create a new migration file"
	@echo "  migrate-up     - Run migrations up"
	@echo "  migrate-down   - Run migrations down"
	@echo "  docker-up      - Start PostgreSQL in Docker"
	@echo "  docker-down    - Stop and remove PostgreSQL Docker container"
	@echo "  help           - Show this help message"
