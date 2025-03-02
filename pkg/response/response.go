package response

import (
	"net/http"
	"time"
)

// Status represents the status of a response
type Status string

const (
	// StatusSuccess represents a successful response
	StatusSuccess Status = "success"
	// StatusError represents an error response
	StatusError Status = "error"
	// StatusWarning represents a warning response
	StatusWarning Status = "warning"
)

// Response is a generic API response structure
type Response struct {
	Status    Status      `json:"status"`              // Status of the response (success, error, warning)
	Message   string      `json:"message,omitempty"`   // Human-readable message
	Data      interface{} `json:"data,omitempty"`      // Response data (for success responses)
	Timestamp time.Time   `json:"timestamp"`           // Timestamp of the response
	RequestID string      `json:"request_id,omitempty"` // Request ID for tracing
}

// MetaData represents metadata for paginated responses
type MetaData struct {
	TotalCount int `json:"total_count"` // Total number of items
	Page       int `json:"page"`        // Current page
	PageSize   int `json:"page_size"`   // Number of items per page
	TotalPages int `json:"total_pages"` // Total number of pages
}

// PaginatedResponse is a response with pagination metadata
type PaginatedResponse struct {
	Response
	Meta MetaData `json:"meta"` // Pagination metadata
}

// ErrorResponse is a detailed error response
type ErrorResponse struct {
	Response
	Code    int                    `json:"code"`              // HTTP status code or application-specific error code
	Details map[string]interface{} `json:"details,omitempty"` // Additional error details
	Errors  map[string]string      `json:"errors,omitempty"`  // Validation errors by field
}

// NewSuccess creates a new success response
func NewSuccess(message string, data interface{}) Response {
	return Response{
		Status:    StatusSuccess,
		Message:   message,
		Data:      data,
		Timestamp: timeNow(),
	}
}

// NewSuccessWithRequestID creates a new success response with request ID
func NewSuccessWithRequestID(message string, data interface{}, requestID string) Response {
	resp := NewSuccess(message, data)
	resp.RequestID = requestID
	return resp
}

// NewPaginated creates a new paginated response
func NewPaginated(message string, data interface{}, totalCount, page, pageSize, totalPages int) PaginatedResponse {
	return PaginatedResponse{
		Response: Response{
			Status:    StatusSuccess,
			Message:   message,
			Data:      data,
			Timestamp: timeNow(),
		},
		Meta: MetaData{
			TotalCount: totalCount,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	}
}

// NewPaginatedWithRequestID creates a new paginated response with request ID
func NewPaginatedWithRequestID(message string, data interface{}, totalCount, page, pageSize, totalPages int, requestID string) PaginatedResponse {
	resp := NewPaginated(message, data, totalCount, page, pageSize, totalPages)
	resp.RequestID = requestID
	return resp
}

// NewError creates a new error response
func NewError(message string, code int) ErrorResponse {
	return ErrorResponse{
		Response: Response{
			Status:    StatusError,
			Message:   message,
			Timestamp: timeNow(),
		},
		Code: code,
	}
}

// NewErrorWithRequestID creates a new error response with request ID
func NewErrorWithRequestID(message string, code int, requestID string) ErrorResponse {
	resp := NewError(message, code)
	resp.RequestID = requestID
	return resp
}

// NewValidationError creates a new validation error response
func NewValidationError(message string, errors map[string]string) ErrorResponse {
	return ErrorResponse{
		Response: Response{
			Status:    StatusError,
			Message:   message,
			Timestamp: timeNow(),
		},
		Code:   http.StatusBadRequest,
		Errors: errors,
	}
}

// NewValidationErrorWithRequestID creates a new validation error response with request ID
func NewValidationErrorWithRequestID(message string, errors map[string]string, requestID string) ErrorResponse {
	resp := NewValidationError(message, errors)
	resp.RequestID = requestID
	return resp
}

// WithDetails adds details to an error response
func (e ErrorResponse) WithDetails(details map[string]interface{}) ErrorResponse {
	e.Details = details
	return e
}

// WithErrors adds validation errors to an error response
func (e ErrorResponse) WithErrors(errors map[string]string) ErrorResponse {
	e.Errors = errors
	return e
}
