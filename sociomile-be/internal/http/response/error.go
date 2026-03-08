package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewHTTPError(code int, message string) *echo.HTTPError {
	return echo.NewHTTPError(code, message)
}

func BadRequest(message string) *echo.HTTPError {
	return NewHTTPError(http.StatusBadRequest, message)
}

func Unauthorized(message string) *echo.HTTPError {
	return NewHTTPError(http.StatusUnauthorized, message)
}

func Forbidden(message string) *echo.HTTPError {
	return NewHTTPError(http.StatusForbidden, message)
}

func NotFound(message string) *echo.HTTPError {
	return NewHTTPError(http.StatusNotFound, message)
}

func Conflict(message string) *echo.HTTPError {
	return NewHTTPError(http.StatusConflict, message)
}

func TooManyRequests(message string) *echo.HTTPError {
	return NewHTTPError(http.StatusTooManyRequests, message)
}

func InternalServerError(message string) *echo.HTTPError {
	return NewHTTPError(http.StatusInternalServerError, message)
}

func ValidationFailed(err error) *echo.HTTPError {
	return BadRequest("validation failed").SetInternal(err)
}
