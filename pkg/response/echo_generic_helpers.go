package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RespondWithGenericOK sends a generic OK response
func RespondWithGenericOK[T any](c echo.Context, message string, data T) error {
	return RespondWithGenericSuccess(c, http.StatusOK, message, data)
}

// RespondWithGenericCreated sends a generic created response
func RespondWithGenericCreated[T any](c echo.Context, message string, data T) error {
	return RespondWithGenericSuccess(c, http.StatusCreated, message, data)
}

// RespondWithGenericSuccess sends a generic success response
func RespondWithGenericSuccess[T any](c echo.Context, statusCode int, message string, data T) error {
	requestID := getRequestID(c)
	var resp GenericResponse[T]

	if requestID != "" {
		resp = NewGenericSuccessWithRequestID(message, data, requestID)
	} else {
		resp = NewGenericSuccess(message, data)
	}

	return c.JSON(statusCode, resp)
}

// RespondWithGenericError sends a generic error response
func RespondWithGenericError[T any](c echo.Context, statusCode int, message string, data T) error {
	requestID := getRequestID(c)
	var resp GenericErrorResponse[T]

	if requestID != "" {
		resp = NewGenericErrorWithRequestID(message, statusCode, data, requestID)
	} else {
		resp = NewGenericError(message, statusCode, data)
	}

	return c.JSON(statusCode, resp)
}

// RespondWithGenericBadRequest sends a generic bad request response
func RespondWithGenericBadRequest[T any](c echo.Context, message string, data T) error {
	return RespondWithGenericError(c, http.StatusBadRequest, message, data)
}

// RespondWithGenericUnauthorized sends a generic unauthorized response
func RespondWithGenericUnauthorized[T any](c echo.Context, message string, data T) error {
	return RespondWithGenericError(c, http.StatusUnauthorized, message, data)
}

// RespondWithGenericForbidden sends a generic forbidden response
func RespondWithGenericForbidden[T any](c echo.Context, message string, data T) error {
	return RespondWithGenericError(c, http.StatusForbidden, message, data)
}

// RespondWithGenericNotFound sends a generic not found response
func RespondWithGenericNotFound[T any](c echo.Context, message string, data T) error {
	return RespondWithGenericError(c, http.StatusNotFound, message, data)
}

// RespondWithGenericInternalError sends a generic internal server error response
func RespondWithGenericInternalError[T any](c echo.Context, message string, data T) error {
	return RespondWithGenericError(c, http.StatusInternalServerError, message, data)
} 
