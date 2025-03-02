package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GenericErrorResponse is a generic error response with type-safe details
type GenericErrorResponse[T any] struct {
	Response
	Code    int            `json:"code"`              // HTTP status code or application-specific error code
	Details T              `json:"details,omitempty"` // Additional error details with type safety
	Errors  map[string]string `json:"errors,omitempty"`  // Validation errors by field
}

// NewGenericError creates a new error response with type-safe details
func NewGenericError[T any](message string, code int, details T) GenericErrorResponse[T] {
	return GenericErrorResponse[T]{
		Response: Response{
			Status:    StatusError,
			Message:   message,
			Timestamp: timeNow(),
		},
		Code:    code,
		Details: details,
	}
}

// NewGenericErrorWithRequestID creates a new error response with type-safe details and request ID
func NewGenericErrorWithRequestID[T any](message string, code int, details T, requestID string) GenericErrorResponse[T] {
	resp := NewGenericError(message, code, details)
	resp.RequestID = requestID
	return resp
}

// RespondWithGenericError sends an error response with type-safe details
func RespondWithGenericError[T any](c *gin.Context, statusCode int, message string, details T) {
	requestID := getRequestID(c)
	var resp GenericErrorResponse[T]

	if requestID != "" {
		resp = NewGenericErrorWithRequestID(message, statusCode, details, requestID)
	} else {
		resp = NewGenericError(message, statusCode, details)
	}

	c.JSON(statusCode, resp)
}

// RespondWithGenericBadRequest sends a bad request response with type-safe details
func RespondWithGenericBadRequest[T any](c *gin.Context, message string, details T) {
	RespondWithGenericError(c, http.StatusBadRequest, message, details)
}

// RespondWithGenericUnauthorized sends an unauthorized response with type-safe details
func RespondWithGenericUnauthorized[T any](c *gin.Context, message string, details T) {
	RespondWithGenericError(c, http.StatusUnauthorized, message, details)
}

// RespondWithGenericForbidden sends a forbidden response with type-safe details
func RespondWithGenericForbidden[T any](c *gin.Context, message string, details T) {
	RespondWithGenericError(c, http.StatusForbidden, message, details)
}

// RespondWithGenericNotFound sends a not found response with type-safe details
func RespondWithGenericNotFound[T any](c *gin.Context, message string, details T) {
	RespondWithGenericError(c, http.StatusNotFound, message, details)
}

// RespondWithGenericInternalError sends an internal server error response with type-safe details
func RespondWithGenericInternalError[T any](c *gin.Context, message string, details T) {
	RespondWithGenericError(c, http.StatusInternalServerError, message, details)
}

// Helper function to get current time (extracted for testing)
func timeNow() time.Time {
	return time.Now().UTC()
}
