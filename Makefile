.PHONY: build run clean test lint migrate-up migrate-down migrate-create sqlc docker-up docker-down

# Binary output
BIN_DIR = bin
BIN_NAME = api

# Build settings
GO = go
GOFLAGS = -v
BUILD_DIR = ./cmd/api

# Database settings
DB_USER = user
DB_PASSWORD = pass
DB_NAME = yadwy_db
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


tidy:
	@echo "Running go mod tidy..."
	$(GO) mod tidy

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	$(GO) clean

test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

migrate-create:
	@echo "Creating migration file..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $${name}

migrate-up:
	@echo "Running migrations up..."
	migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	@echo "Running migrations down..."
	migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

docker-up:
	@echo "Starting PostgreSQL in Docker..."
	docker-compose -f docker-compose.yml up

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
