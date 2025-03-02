package metrics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestsTotal counts the number of HTTP requests processed
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, partitioned by status code, method, and path",
		},
		[]string{"code", "method", "path"},
	)

	// RequestDuration tracks the duration of HTTP requests
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"code", "method", "path"},
	)

	// ActiveRequests tracks the number of in-flight requests
	ActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_active",
			Help: "Number of active HTTP requests",
		},
	)

	// ErrorsTotal counts the number of errors
	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of errors",
		},
		[]string{"type"},
	)

	// DatabaseOperationsTotal counts the number of database operations
	DatabaseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "entity"},
	)

	// DatabaseOperationDuration tracks the duration of database operations
	DatabaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_operation_duration_seconds",
			Help:    "Duration of database operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "entity"},
	)
)

// MetricsMiddleware returns a middleware that collects metrics for HTTP requests
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Increment active requests
		ActiveRequests.Inc()

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Decrement active requests
		ActiveRequests.Dec()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get status code as string
		statusCode := c.Writer.Status()
		statusCodeStr := http.StatusText(statusCode)

		// Get request method and path
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		// Record metrics
		RequestsTotal.WithLabelValues(statusCodeStr, method, path).Inc()
		RequestDuration.WithLabelValues(statusCodeStr, method, path).Observe(duration)

		// Record errors
		if statusCode >= 400 {
			ErrorsTotal.WithLabelValues("http").Inc()
		}
	}
}

// RegisterMetricsEndpoint registers the metrics endpoint
func RegisterMetricsEndpoint(router *gin.RouterGroup) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// MeasureDatabaseOperation measures the duration of a database operation
func MeasureDatabaseOperation(operation, entity string, fn func() error) error {
	// Start timer
	start := time.Now()

	// Execute operation
	err := fn()

	// Calculate duration
	duration := time.Since(start).Seconds()

	// Record metrics
	DatabaseOperationsTotal.WithLabelValues(operation, entity).Inc()
	DatabaseOperationDuration.WithLabelValues(operation, entity).Observe(duration)

	// Record errors
	if err != nil {
		ErrorsTotal.WithLabelValues("database").Inc()
	}

	return err
}
