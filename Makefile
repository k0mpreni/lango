# Simple Makefile for a Go project

# Build the application
all: build

install:
	@echo "Installing..."
	cd ./cmd/web && npm i && cd ../.. 

build:
	@echo "Building..."
	cd ./cmd/web && npm run css && cd ../.. 
	@templ generate
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

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
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

# up: ## Database migration up
# 		@echo 'Running up migrations...'
# 		migrate -path ./internal/database/migrate/migrations -database ${DB_DSN} --verbose up

up: ## Database migration up
		@go run ./internal/database/migrate/migrate.go up

migration:
		@migrate create -ext sql -dir internal/database/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))


.PHONY: all build run test clean
