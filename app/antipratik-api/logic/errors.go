package logic

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError is returned when user input fails validation.
// API handlers check for this type to return 400 Bad Request instead of 500.
// All future error types specific to the logic layer also live here.
type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string { return e.msg }

func validationErr(msg string) error { return &ValidationError{msg: msg} }

// IsValidationError reports whether err is (or wraps) a ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

// requireNonEmpty returns a ValidationError if the trimmed value is empty.
func requireNonEmpty(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return validationErr(fmt.Sprintf("%s cannot be empty", field))
	}
	return nil
}

// requirePositive returns a ValidationError if v is not greater than zero.
func requirePositive(field string, v int) error {
	if v <= 0 {
		return validationErr(fmt.Sprintf("%s must be greater than zero", field))
	}
	return nil
}
