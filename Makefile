.PHONY: build run test clean migrate-up migrate-down docker-up docker-down docker-up-logs docker-rebuild deps deps-update lint docs help dev prod-build prod-up prod-down prod-logs prod-full-up prod-full-down gen-certs prod-migrate-up prod-migrate prod-build-migrate prod-build-all

# Build the application
build:
	go build -o bin/api ./cmd/api

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run database migrations up
migrate-up:
	go run cmd/migrate/main.go up

# Run database migrations down
migrate-down:
	go run cmd/migrate/main.go down

# Create a new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	go run cmd/migrate/main.go create $$name

# Start Docker containers
docker-up:
	docker-compose up -d

# Stop Docker containers
docker-down:
	docker-compose down

# Start Docker containers and follow logs
docker-up-logs:
	docker-compose up

# Rebuild Docker containers
docker-rebuild:
	docker-compose up -d --build

# Download Go dependencies
deps:
	go mod download

# Update Go dependencies
deps-update:
	go get -u ./...
	go mod tidy

# Run linter
lint:
	go vet ./...

# Generate Go documentation
docs:
	godoc -http=:6060

# Run the application with hot reload using Air
dev:
	@command -v air > /dev/null 2>&1 || { \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	}
	air

# Help command
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run the application with hot reload using Air"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make migrate-up     - Run database migrations up"
	@echo "  make migrate-down   - Run database migrations down"
	@echo "  make migrate-create - Create a new migration"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-up-logs - Start Docker containers and follow logs"
	@echo "  make docker-rebuild - Rebuild Docker containers"
	@echo "  make deps           - Download Go dependencies"
	@echo "  make deps-update    - Update Go dependencies"
	@echo "  make lint           - Run linter"
	@echo "  make docs           - Generate Go documentation"
	@echo "  make prod-build     - Build production images"
	@echo "  make prod-up        - Start production stack"
	@echo "  make prod-down      - Stop production stack"
	@echo "  make prod-logs      - View production logs"
	@echo "  make prod-full-up   - Start full production stack with monitoring"
	@echo "  make prod-full-down - Stop full production stack"
	@echo "  make prod-migrate-up - Run database migrations in production (full stack)"
	@echo "  make prod-migrate   - Run database migrations in production (basic stack)"
	@echo "  make prod-build-migrate - Build migration image"
	@echo "  make prod-build-all - Build all production images"
	@echo "  make help           - Show this help message"

# Production deployment targets

# Build production images
prod-build:
	docker-compose -f docker-compose.prod.yml build

# Build migration image
prod-build-migrate:
	docker build -t todo-api-migrate:latest -f Dockerfile.migrate .

# Build all production images
prod-build-all: prod-build prod-build-migrate

# Start production stack
prod-up:
	docker-compose -f docker-compose.prod.yml up -d

# Stop production stack
prod-down:
	docker-compose -f docker-compose.prod.yml down

# View production logs
prod-logs:
	docker-compose -f docker-compose.prod.yml logs -f

# Start full production stack with all monitoring components
prod-full-up:
	docker-compose -f docker-compose.prod.full.yml up -d

# Stop full production stack
prod-full-down:
	docker-compose -f docker-compose.prod.full.yml down

# Generate self-signed SSL certificates for development
gen-certs:
	mkdir -p nginx/ssl
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout nginx/ssl/todo-api.key -out nginx/ssl/todo-api.crt \
		-subj "/C=US/ST=State/L=City/O=Organization/CN=todo-api.example.com"

# Run database migrations in production (full stack)
prod-migrate-up:
	docker-compose -f docker-compose.prod.full.yml run --rm migrations

# Run database migrations in production (basic stack)
prod-migrate:
	docker-compose -f docker-compose.prod.yml run --rm migrations
