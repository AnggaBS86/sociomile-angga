package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SuccessEnvelope struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorEnvelope struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func OK(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, SuccessEnvelope{Status: "OK", Message: message, Data: data})
}

func Created(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusCreated, SuccessEnvelope{Status: "OK", Message: message, Data: data})
}

func Accepted(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusAccepted, SuccessEnvelope{Status: "OK", Message: message, Data: data})
}
