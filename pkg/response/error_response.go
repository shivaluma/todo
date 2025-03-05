package response

// GenericErrorResponse is a generic error response with type-safe details
type GenericErrorResponse[T any] struct {
	Response
	Code    int               `json:"code"`              // HTTP status code or application-specific error code
	Details T                 `json:"details,omitempty"` // Additional error details with type safety
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
