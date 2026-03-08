package middleware

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if cv == nil || cv.Validator == nil {
		return errors.New("validator is not initialized")
	}

	return cv.Validator.Struct(i)
}
