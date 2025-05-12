.PHONY: build run dev test clean migrate migrate-create migrate-down mock

# Application name
APP_NAME = personal-website-backend

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOTEST = $(GOCMD) test
GOTIDY = $(GOCMD) mod tidy
GOMOD = $(GOCMD) mod
GOGET = $(GOCMD) get
GOINSTALL = $(GOCMD) install

# Main package
MAIN_PKG = ./cmd/api

# Build variables
BINARY_PATH = ./bin

# Build the application
build:
	@echo "Building application..."
	$(GOBUILD) -o $(BINARY_PATH)/$(APP_NAME) $(MAIN_PKG)

# Run the application
run: build
	@echo "Running application..."
	$(BINARY_PATH)/$(APP_NAME)

dev:
	@echo "Running application in development mode..."
	$(GORUN) ./cmd/api

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BINARY_PATH)

# Run database migrations
migrate:
	@echo "Running database migrations..."
	$(GORUN) $(MAIN_PKG) db:migrate

# Create a new migration
migrate-create:
	@echo "Creating migration..."
	@read -p "Enter migration name: " name; \
	$(GORUN) $(MAIN_PKG)/main.go db:create $$name

# Rollback database migrations
migrate-down:
	@echo "Rolling back database migrations..."
	$(GORUN) $(MAIN_PKG)/main.go db:rollback

# Generate mocks for testing
mock:
	@echo "Generating mocks..."
	mockery --name=AuthService --dir=internal/service --output=internal/service/mocks
	mockery --name=ArticleService --dir=internal/service --output=internal/service/mocks
	mockery --name=PortfolioService --dir=internal/service --output=internal/service/mocks
	mockery --name=UserRepository --dir=internal/repository --output=internal/repository/mocks
	mockery --name=ArticleRepository --dir=internal/repository --output=internal/repository/mocks
	mockery --name=PortfolioRepository --dir=internal/repository --output=internal/repository/mocks

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOTIDY)
	$(GOGET) -u all

# Generate a new module
generate-module:
	@echo "Generating module..."
	@read -p "Enter module name: " name; \
	mkdir -p internal/$(name); \
	touch internal/$(name)/$(name).go

# Help command
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Build and run the application"
	@echo "  make dev            - Run the application in development mode"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make migrate        - Run database migrations"
	@echo "  make migrate-create - Create a new migration"
	@echo "  make migrate-down   - Rollback database migrations"
	@echo "  make mock           - Generate mocks for testing"
	@echo "  make deps           - Install dependencies"
	@echo "  make generate-module - Generate a new module"
	@echo "  make help           - Display this help message"

# Default target
default: help 