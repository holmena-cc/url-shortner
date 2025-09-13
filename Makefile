# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

tailwind:
	npx tailwindcss -i ./web/static/input.css -o ./web/static/output.css

# Create DB container
docker-run:
	@docker compose up --build

# Shutdown DB container
docker-down:
	@docker compose down
# run the migrations
migrate:
	docker exec -i urlshortner-psql_bp-1 bash -c "for f in /docker-entrypoint-initdb.d/migrations/*.sql; do psql -U melkey -d blueprint -f $$f; done"
# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Format code and tidy modules
format:
	@go fmt ./...
	@go mod tidy 

.PHONY: build run tailwind docker-run docker-down test itest clean format
