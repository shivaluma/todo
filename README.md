# Todo API

A modern Todo List API built with Go, following CQRS architecture and 2025 best practices.

## Features

-   CRUD operations for Todo items
-   User authentication with JWT
-   CQRS pattern implementation
-   PostgreSQL database
-   Containerized with Docker
-   Comprehensive test coverage

## Architecture

This project follows a clean, layered architecture with CQRS pattern:

-   **Domain Layer**: Contains the core business logic and entities
-   **Application Layer**: Contains commands and queries that orchestrate the domain
-   **Infrastructure Layer**: Contains implementations of repositories and external services
-   **Interfaces Layer**: Contains API handlers and middleware

## Prerequisites

-   Go 1.23+
-   Docker and Docker Compose
-   PostgreSQL (if running locally)

## Getting Started

### Running with Docker

```bash
# Clone the repository
git clone https://github.com/sh1ro/todo-api.git
cd todo-api

# Start the application and database
docker-compose up -d
```

The API will be available at http://localhost:8080

### Running Locally

```bash
# Clone the repository
git clone https://github.com/sh1ro/todo-api.git
cd todo-api

# Install dependencies
go mod download

# Set up environment variables (see .env.example)
cp .env.example .env
# Edit .env with your configuration

# Run database migrations
go run cmd/migrate/main.go up

# Run the application
go run cmd/api/main.go
```

## API Endpoints

### Authentication

-   `POST /api/v1/auth/register` - Register a new user
-   `POST /api/v1/auth/login` - Login and get JWT token

### Todo Items

-   `GET /api/v1/todos` - Get all todos for the authenticated user
-   `GET /api/v1/todos/:id` - Get a specific todo
-   `POST /api/v1/todos` - Create a new todo
-   `PUT /api/v1/todos/:id` - Update a todo
-   `DELETE /api/v1/todos/:id` - Delete a todo

## Project Structure

```
.
├── cmd/                    # Application entry points
│   ├── api/                # Main API application
│   └── migrate/            # Database migration tool
├── internal/               # Private application code
│   └── app/
│       ├── domain/         # Domain layer
│       │   ├── model/      # Domain models
│       │   ├── repository/ # Repository interfaces
│       │   └── service/    # Domain services
│       ├── application/    # Application layer
│       │   ├── command/    # Command handlers
│       │   └── query/      # Query handlers
│       ├── infrastructure/ # Infrastructure layer
│       │   ├── persistence/# Repository implementations
│       │   └── auth/       # Authentication services
│       └── interfaces/     # Interface layer
│           ├── api/        # API handlers
│           └── middleware/ # HTTP middleware
├── pkg/                    # Public libraries
│   ├── validator/          # Validation utilities
│   ├── logger/             # Logging utilities
│   └── config/             # Configuration utilities
├── migrations/             # Database migrations
├── .env.example            # Example environment variables
├── docker-compose.yml      # Docker compose configuration
├── Dockerfile              # Docker build configuration
└── README.md               # Project documentation
```

## License

MIT

# Request Tracing

The API implements request tracing using unique request IDs. Each request is assigned a UUID that is:

-   Generated automatically if not provided
-   Accepted from clients via the `X-Request-ID` header
-   Included in all response headers as `X-Request-ID`
-   Added to all log entries related to the request
-   Available to all handlers via the context

This enables end-to-end tracing of requests across distributed systems and simplifies debugging.

## Standardized API Responses

This API uses standardized response formats for both successful and error responses to ensure consistency across all endpoints.

### Success Response Format

All successful responses follow this structure:

```json
{
  "status": "success",
  "message": "Human-readable success message",
  "data": { ... },
  "timestamp": "2023-03-01T12:34:56Z",
  "request_id": "optional-request-id-for-tracing"
}
```

### Paginated Response Format

For endpoints that return paginated data:

```json
{
  "status": "success",
  "message": "Human-readable success message",
  "data": [ ... ],
  "timestamp": "2023-03-01T12:34:56Z",
  "request_id": "optional-request-id-for-tracing",
  "meta": {
    "total_count": 100,
    "page": 1,
    "page_size": 10,
    "total_pages": 10
  }
}
```

### Error Response Format

All error responses follow this structure:

```json
{
	"status": "error",
	"message": "Human-readable error message",
	"timestamp": "2023-03-01T12:34:56Z",
	"request_id": "optional-request-id-for-tracing",
	"code": 400,
	"errors": {
		"field1": "Error message for field1",
		"field2": "Error message for field2"
	}
}
```

### Type-Safe Responses with Generics

The API also provides type-safe response templates using Go generics:

```go
// Example of a type-safe success response
type UserResponse struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

userData := UserResponse{
    ID:       user.ID.String(),
    Username: user.Username,
    Email:    user.Email,
}

response.RespondWithGenericCreated(c, "User registered successfully", userData)
```

This ensures that the response data is strongly typed and provides better compile-time type checking.

## Response Helper Functions

The API uses a set of helper functions to generate standardized responses:

### Standard Response Helpers

-   `RespondWithSuccess` - For general success responses
-   `RespondWithPaginated` - For paginated responses
-   `RespondWithError` - For error responses
-   `RespondWithValidationError` - For validation error responses
-   `RespondWithCreated` - For 201 Created responses
-   `RespondWithOK` - For 200 OK responses
-   `RespondWithNoContent` - For 204 No Content responses
-   `RespondWithBadRequest` - For 400 Bad Request responses
-   `RespondWithUnauthorized` - For 401 Unauthorized responses
-   `RespondWithForbidden` - For 403 Forbidden responses
-   `RespondWithNotFound` - For 404 Not Found responses
-   `RespondWithInternalError` - For 500 Internal Server Error responses

### Generic Response Helpers (Type-Safe)

-   `RespondWithGenericSuccess` - For type-safe success responses
-   `RespondWithGenericPaginated` - For type-safe paginated responses
-   `RespondWithGenericCreated` - For type-safe 201 Created responses
-   `RespondWithGenericOK` - For type-safe 200 OK responses
-   `RespondWithGenericError` - For type-safe error responses
-   `RespondWithGenericBadRequest` - For type-safe 400 Bad Request responses
-   `RespondWithGenericUnauthorized` - For type-safe 401 Unauthorized responses
-   `RespondWithGenericForbidden` - For type-safe 403 Forbidden responses
-   `RespondWithGenericNotFound` - For type-safe 404 Not Found responses
-   `RespondWithGenericInternalError` - For type-safe 500 Internal Server Error responses

## Running the API

```bash
# Build the API
go build -o bin/api ./cmd/api

# Run the API
./bin/api
```

## Environment Variables

The API uses the following environment variables:

-   `PORT` - The port to listen on (default: 8080)
-   `DB_HOST` - Database host
-   `DB_PORT` - Database port
-   `DB_USER` - Database user
-   `DB_PASSWORD` - Database password
-   `DB_NAME` - Database name
-   `JWT_SECRET` - Secret key for JWT tokens
-   `JWT_EXPIRATION` - JWT token expiration time in hours
-   `LOG_LEVEL` - Logging level (debug, info, warn, error)
-   `LOG_FORMAT` - Logging format (json, text)
-   `API_VERSION` - API version (default: v1)

# Monitoring with Prometheus and Grafana

The Todo API includes a comprehensive monitoring setup using Prometheus and Grafana.

## Monitoring Components

-   **Prometheus**: Collects and stores metrics from the API and other services
-   **Grafana**: Visualizes metrics with customizable dashboards
-   **Alertmanager**: Handles alerts from Prometheus and sends notifications
-   **Node Exporter**: Collects system metrics from the host
-   **cAdvisor**: Collects container metrics
-   **Postgres Exporter**: Collects metrics from PostgreSQL
-   **Nginx Exporter**: Collects metrics from Nginx

## Available Metrics

The API exposes the following metrics:

-   HTTP request count, duration, and error rates
-   Database operation count, duration, and error rates
-   Active request count
-   System metrics (CPU, memory, disk usage)
-   Container metrics
-   Database metrics
-   Nginx metrics

## Dashboards

Grafana comes pre-configured with dashboards for:

-   API overview (request rate, latency, error rate)
-   Database performance
-   System resources
-   Container resources

## Alerting

Prometheus is configured with alerting rules for:

-   High request latency
-   High error rates
-   API service availability
-   Database performance issues
-   System resource usage

## Production Deployment

For production deployment, use the production Docker Compose file:

```bash
# Start the production stack
docker-compose -f docker-compose.prod.yml up -d

# Start the full production stack with additional monitoring components
docker-compose -f docker-compose.prod.full.yml up -d
```

## Accessing Monitoring Tools

-   Grafana: http://localhost:3000 (default credentials: admin/admin)
-   Prometheus: http://localhost:9090
-   Alertmanager: http://localhost:9093

In production, these services are protected behind Nginx with proper authentication.

## Configuration

-   Prometheus configuration: `monitoring/prometheus/prometheus.prod.yml`
-   Alerting rules: `monitoring/prometheus/rules/alerts.yml`
-   Alertmanager configuration: `monitoring/alertmanager/alertmanager.yml`
-   Grafana dashboards: `monitoring/grafana/dashboards/`
-   Grafana datasources: `monitoring/grafana/provisioning/datasources/`
# todo
