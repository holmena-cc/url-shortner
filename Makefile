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
	npx tailwindcss -i ./web/static/input.css -o ./web/static/output.css --watch

# Create DB container
docker-run:
	@docker compose up --build

# Shutdown DB container
docker-down:
	@docker compose down

# Build and run the Docker container with live reload (Air)
# ! Include PORT in the command (e.g. make docker-dev PORT=5000)
docker-dev:
	@echo "Building Docker image..."
	@docker build --target dev -t url-shortener-dev .
	@echo "Running container with live reload..."
	@docker run -p ${PORT}:${PORT} --rm \
		-v "${CURDIR}:/app" \
		-v "${CURDIR}/tmp:/app/tmp" \
		--name url-shortener-air \
		url-shortener-dev

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
	go fmt ./...
	go mod tidy 

# Live Reload
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"

.PHONY: all build run test clean watch docker-run docker-down itest
