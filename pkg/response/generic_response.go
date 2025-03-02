package response

import (
	"time"
)

// GenericResponse is a generic API response structure with type-safe data
type GenericResponse[T any] struct {
	Status    Status    `json:"status"`              // Status of the response (success, error, warning)
	Message   string    `json:"message,omitempty"`   // Human-readable message
	Data      T         `json:"data,omitempty"`      // Response data (for success responses)
	Timestamp time.Time `json:"timestamp"`           // Timestamp of the response
	RequestID string    `json:"request_id,omitempty"` // Request ID for tracing
}

// GenericPaginatedResponse is a paginated response with type-safe data
type GenericPaginatedResponse[T any] struct {
	Status    Status    `json:"status"`              // Status of the response (success, error, warning)
	Message   string    `json:"message,omitempty"`   // Human-readable message
	Data      T         `json:"data,omitempty"`      // Response data (for success responses)
	Timestamp time.Time `json:"timestamp"`           // Timestamp of the response
	RequestID string    `json:"request_id,omitempty"` // Request ID for tracing
	Meta      MetaData  `json:"meta"`                // Pagination metadata
}

// NewGenericSuccess creates a new success response with type-safe data
func NewGenericSuccess[T any](message string, data T) GenericResponse[T] {
	return GenericResponse[T]{
		Status:    StatusSuccess,
		Message:   message,
		Data:      data,
		Timestamp: timeNow(),
	}
}

// NewGenericSuccessWithRequestID creates a new success response with type-safe data and request ID
func NewGenericSuccessWithRequestID[T any](message string, data T, requestID string) GenericResponse[T] {
	resp := NewGenericSuccess(message, data)
	resp.RequestID = requestID
	return resp
}

// NewGenericPaginated creates a new paginated response with type-safe data
func NewGenericPaginated[T any](message string, data T, totalCount, page, pageSize, totalPages int) GenericPaginatedResponse[T] {
	return GenericPaginatedResponse[T]{
		Status:    StatusSuccess,
		Message:   message,
		Data:      data,
		Timestamp: timeNow(),
		Meta: MetaData{
			TotalCount: totalCount,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	}
}

// NewGenericPaginatedWithRequestID creates a new paginated response with type-safe data and request ID
func NewGenericPaginatedWithRequestID[T any](message string, data T, totalCount, page, pageSize, totalPages int, requestID string) GenericPaginatedResponse[T] {
	resp := NewGenericPaginated(message, data, totalCount, page, pageSize, totalPages)
	resp.RequestID = requestID
	return resp
}
