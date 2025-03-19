# Simple Makefile for wallet API
RED := \e[31m
GREEN := \e[32m
CYAN := \033[36m
BOLD := \e[1m
RESET := \033[0m

build: ## Build the application
	@echo "Building..."
	
	
	@CGO_ENABLED=0 go build -o bin/wallet -ldflags "-s -w" cmd/main.go


run: ## Run the application
	@go run cmd/main.go


test: ## Test the application
	@echo "Testing..."
	@go test ./... -v


clean: ## Clean the binary
	@echo "Cleaning..."
	@rm -f main


watch: ## Live Reload
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

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'

.DEFAULT_GOAL:= help
.PHONY: all build run test clean watch
