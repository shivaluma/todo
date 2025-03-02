package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sh1ro/todo-api/internal/app/application/query"
	"github.com/sh1ro/todo-api/pkg/validator"
)

// RespondWithSuccess sends a success response
func RespondWithSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	requestID := getRequestID(c)
	var resp Response

	if requestID != "" {
		resp = NewSuccessWithRequestID(message, data, requestID)
	} else {
		resp = NewSuccess(message, data)
	}

	c.JSON(statusCode, resp)
}

// RespondWithPaginated sends a paginated response
func RespondWithPaginated(c *gin.Context, statusCode int, message string, data interface{}, result *query.TodosResult) {
	requestID := getRequestID(c)
	var resp PaginatedResponse

	if requestID != "" {
		resp = NewPaginatedWithRequestID(
			message,
			data,
			result.TotalCount,
			result.Page,
			result.PageSize,
			result.TotalPages,
			requestID,
		)
	} else {
		resp = NewPaginated(
			message,
			data,
			result.TotalCount,
			result.Page,
			result.PageSize,
			result.TotalPages,
		)
	}

	c.JSON(statusCode, resp)
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, statusCode int, message string) {
	requestID := getRequestID(c)
	var resp ErrorResponse

	if requestID != "" {
		resp = NewErrorWithRequestID(message, statusCode, requestID)
	} else {
		resp = NewError(message, statusCode)
	}

	c.JSON(statusCode, resp)
}

// RespondWithValidationError sends a validation error response
func RespondWithValidationError(c *gin.Context, message string, errors []validator.ValidationError) {
	requestID := getRequestID(c)

	// Convert validation errors to map
	errorMap := make(map[string]string)
	for _, err := range errors {
		errorMap[err.Field] = err.Message
	}

	var resp ErrorResponse
	if requestID != "" {
		resp = NewValidationErrorWithRequestID(message, errorMap, requestID)
	} else {
		resp = NewValidationError(message, errorMap)
	}

	c.JSON(http.StatusBadRequest, resp)
}

// RespondWithCreated sends a created response
func RespondWithCreated(c *gin.Context, message string, data interface{}) {
	RespondWithSuccess(c, http.StatusCreated, message, data)
}

// RespondWithOK sends an OK response
func RespondWithOK(c *gin.Context, message string, data interface{}) {
	RespondWithSuccess(c, http.StatusOK, message, data)
}

// RespondWithNoContent sends a no content response
func RespondWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// RespondWithBadRequest sends a bad request response
func RespondWithBadRequest(c *gin.Context, message string) {
	RespondWithError(c, http.StatusBadRequest, message)
}

// RespondWithUnauthorized sends an unauthorized response
func RespondWithUnauthorized(c *gin.Context, message string) {
	RespondWithError(c, http.StatusUnauthorized, message)
}

// RespondWithForbidden sends a forbidden response
func RespondWithForbidden(c *gin.Context, message string) {
	RespondWithError(c, http.StatusForbidden, message)
}

// RespondWithNotFound sends a not found response
func RespondWithNotFound(c *gin.Context, message string) {
	RespondWithError(c, http.StatusNotFound, message)
}

// RespondWithInternalError sends an internal server error response
func RespondWithInternalError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusInternalServerError, message)
}

// Helper function to get request ID from context
func getRequestID(c *gin.Context) string {
	if id, exists := c.Get("requestID"); exists {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	return ""
}
