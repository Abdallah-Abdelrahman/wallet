# Simple Makefile for wallet API

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/main.go

# Run the application
run:
	@go run cmd/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

generate-docs: ## Generate API documentation using swag
	@rm -rf docs/swagger/*
	@swag init --dir ./cmd/,./internal/server --output ./docs
	@echo "Documentation generated successfully"

.PHONY: all build run test clean watch
