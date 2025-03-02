# Metrics Package

This package provides Prometheus metrics collection for the Todo API.

## Overview

The metrics package offers:

1. HTTP request metrics (count, duration, active requests)
2. Database operation metrics (count, duration)
3. Error metrics
4. Middleware for automatic collection of HTTP metrics
5. Helper functions for measuring database operations

## Usage

### Middleware Registration

Register the metrics middleware in your application:

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sh1ro/todo-api/pkg/metrics"
)

func main() {
    router := gin.Default()

    // Register metrics middleware
    router.Use(metrics.MetricsMiddleware())

    // Register metrics endpoint
    metrics.RegisterMetricsEndpoint(router)

    // ... rest of your application setup
}
```

### Database Operation Measurement

Wrap database operations with the measurement helper:

```go
import (
    "context"
    "github.com/sh1ro/todo-api/pkg/metrics"
)

func (r *Repository) FindByID(ctx context.Context, id string) (*Entity, error) {
    var entity Entity

    err := metrics.MeasureDatabaseOperation("find", "entity", func() error {
        // Your database operation here
        return db.QueryRowContext(ctx, "SELECT * FROM entities WHERE id = $1", id).Scan(&entity.ID, &entity.Name)
    })

    if err != nil {
        return nil, err
    }

    return &entity, nil
}
```

## Available Metrics

### HTTP Metrics

-   `http_requests_total` - Total number of HTTP requests (labels: code, method, path)
-   `http_request_duration_seconds` - Duration of HTTP requests (labels: code, method, path)
-   `http_requests_active` - Number of active HTTP requests

### Database Metrics

-   `database_operations_total` - Total number of database operations (labels: operation, entity)
-   `database_operation_duration_seconds` - Duration of database operations (labels: operation, entity)

### Error Metrics

-   `http_errors_total` - Total number of HTTP errors (labels: type)

## Prometheus Integration

These metrics are automatically exposed at the `/api/v1/metrics` endpoint in Prometheus format. You can configure Prometheus to scrape this endpoint to collect metrics from your application.

## Grafana Dashboards

Pre-configured Grafana dashboards are available in the `monitoring/grafana/dashboards` directory.
