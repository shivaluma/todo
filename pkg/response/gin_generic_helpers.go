package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sh1ro/todo-api/internal/app/application/query"
)

// RespondWithGenericSuccess sends a success response with type-safe data
func RespondWithGenericSuccess[T any](c *gin.Context, statusCode int, message string, data T) {
	requestID := getRequestID(c)
	var resp GenericResponse[T]

	if requestID != "" {
		resp = NewGenericSuccessWithRequestID(message, data, requestID)
	} else {
		resp = NewGenericSuccess(message, data)
	}

	c.JSON(statusCode, resp)
}

// RespondWithGenericPaginated sends a paginated response with type-safe data
func RespondWithGenericPaginated[T any](c *gin.Context, statusCode int, message string, data T, result *query.TodosResult) {
	requestID := getRequestID(c)
	var resp GenericPaginatedResponse[T]

	if requestID != "" {
		resp = NewGenericPaginatedWithRequestID(
			message,
			data,
			result.TotalCount,
			result.Page,
			result.PageSize,
			result.TotalPages,
			requestID,
		)
	} else {
		resp = NewGenericPaginated(
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

// RespondWithGenericCreated sends a created response with type-safe data
func RespondWithGenericCreated[T any](c *gin.Context, message string, data T) {
	RespondWithGenericSuccess(c, http.StatusCreated, message, data)
}

// RespondWithGenericOK sends an OK response with type-safe data
func RespondWithGenericOK[T any](c *gin.Context, message string, data T) {
	RespondWithGenericSuccess(c, http.StatusOK, message, data)
}
