package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/application/query"
	"github.com/sh1ro/todo-api/pkg/validator"
)

// RespondWithSuccess sends a success response
func RespondWithSuccess(c echo.Context, statusCode int, message string, data interface{}) error {
	requestID := getRequestID(c)
	var resp Response

	if requestID != "" {
		resp = NewSuccessWithRequestID(message, data, requestID)
	} else {
		resp = NewSuccess(message, data)
	}

	return c.JSON(statusCode, resp)
}

// RespondWithPaginated sends a paginated response
func RespondWithPaginated(c echo.Context, statusCode int, message string, data interface{}, result *query.TodosResult) error {
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

	return c.JSON(statusCode, resp)
}

// RespondWithError sends an error response
func RespondWithError(c echo.Context, statusCode int, message string) error {
	requestID := getRequestID(c)
	var resp ErrorResponse

	if requestID != "" {
		resp = NewErrorWithRequestID(message, statusCode, requestID)
	} else {
		resp = NewError(message, statusCode)
	}

	return c.JSON(statusCode, resp)
}

// RespondWithValidationError sends a validation error response
func RespondWithValidationError(c echo.Context, message string, errors []validator.ValidationError) error {
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

	return c.JSON(http.StatusBadRequest, resp)
}

// RespondWithCreated sends a created response
func RespondWithCreated(c echo.Context, message string, data interface{}) error {
	return RespondWithSuccess(c, http.StatusCreated, message, data)
}

// RespondWithOK sends an OK response
func RespondWithOK(c echo.Context, message string, data interface{}) error {
	return RespondWithSuccess(c, http.StatusOK, message, data)
}

// RespondWithNoContent sends a no content response
func RespondWithNoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// RespondWithBadRequest sends a bad request response
func RespondWithBadRequest(c echo.Context, message string) error {
	return RespondWithError(c, http.StatusBadRequest, message)
}

// RespondWithUnauthorized sends an unauthorized response
func RespondWithUnauthorized(c echo.Context, message string) error {
	return RespondWithError(c, http.StatusUnauthorized, message)
}

// RespondWithForbidden sends a forbidden response
func RespondWithForbidden(c echo.Context, message string) error {
	return RespondWithError(c, http.StatusForbidden, message)
}

// RespondWithNotFound sends a not found response
func RespondWithNotFound(c echo.Context, message string) error {
	return RespondWithError(c, http.StatusNotFound, message)
}

// RespondWithInternalError sends an internal server error response
func RespondWithInternalError(c echo.Context, message string) error {
	return RespondWithError(c, http.StatusInternalServerError, message)
}

// getRequestID retrieves the request ID from the context
func getRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
} 
