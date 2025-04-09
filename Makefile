# Variables
BINARY_NAME=finance-tracker
DB_USER=finance
DB_PASSWORD=finance
DB_NAME=finance
DB_HOST=localhost
DB_PORT=5432
POSTGRESQL_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATION_DIR=./migrations
PORT=4000

# Default target
.PHONY: all
all: help

# Help message
.PHONY: help
help:
	@echo "Personal Finance Tracker Commands:"
	@echo "  make setup      - Set up the project (create DB, run migrations)"
	@echo "  make run        - Run the application"
	@echo "  make build      - Build the application"
	@echo "  make migrate-up - Run database migrations up"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean temporary files"

# Create the database
.PHONY: create-db
create-db:
	@echo "Creating database $(DB_NAME)..."
	@psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);" || echo "Database may already exist"

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/server/main.go

# Run migrations up
.PHONY: migrate-up
migrate-up:
	@echo "Running migrations up..."
	@migrate -database $(POSTGRESQL_URL) -path $(MIGRATION_DIR) up

# Run migrations down
.PHONY: migrate-down
migrate-down:
	@echo "Running migrations down..."
	@migrate -database $(POSTGRESQL_URL) -path $(MIGRATION_DIR) down

# Set up the project
.PHONY: setup
setup: create-db migrate-up
	@echo "Project setup complete"

# Run the application
.PHONY: run
run: build
	@echo "Starting Personal Finance Tracker on port $(PORT)..."
	@./$(BINARY_NAME)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./... -v

# Clean temporary files
.PHONY: clean
clean:
	@echo "Cleaning temporary files..."
	@rm -f $(BINARY_NAME)
	@go clean

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod tidy