package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"sociomile-be/internal/http/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ValidationErrorItem struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "internal server error"
	var validationItems []ValidationErrorItem

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		code = httpErr.Code
		if msg, ok := httpErr.Message.(string); ok && msg != "" {
			message = msg
		} else if statusText := http.StatusText(code); statusText != "" {
			message = strings.ToLower(statusText)
		}

		if verrs, ok := httpErr.Internal.(validator.ValidationErrors); ok {
			message = "validation failed"
			validationItems = formatValidationErrors(verrs)
		}
	} else if verrs, ok := err.(validator.ValidationErrors); ok {
		code = http.StatusBadRequest
		message = "validation failed"
		validationItems = formatValidationErrors(verrs)
	}

	if len(validationItems) > 0 {
		_ = c.JSON(code, response.ErrorEnvelope{
			Status:  "ERROR",
			Message: message,
			Errors:  validationItems,
		})
		return
	}

	_ = c.JSON(code, response.ErrorEnvelope{
		Status:  "ERROR",
		Message: message,
	})
}

func formatValidationErrors(verrs validator.ValidationErrors) []ValidationErrorItem {
	items := make([]ValidationErrorItem, 0, len(verrs))
	for _, ve := range verrs {
		field := strings.ToLower(ve.Field())
		items = append(items, ValidationErrorItem{
			Field:   field,
			Rule:    ve.Tag(),
			Message: validationMessage(field, ve.Tag(), ve.Param()),
		})
	}

	return items
}

func validationMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	default:
		return fmt.Sprintf("%s failed %s validation", field, tag)
	}
}
