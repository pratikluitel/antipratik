// Package errors provides shared error types and validation helpers
// used across all components.
package errors

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError is returned when user input fails business-rule validation.
// API handlers check for this type to return 400 Bad Request instead of 500.
type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string { return e.msg }

// New returns a ValidationError with the given message.
func New(msg string) error { return &ValidationError{msg: msg} }

// Is reports whether err is (or wraps) a ValidationError.
func Is(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

// RequireNonEmpty returns a ValidationError if the trimmed value is empty.
func RequireNonEmpty(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return New(fmt.Sprintf("%s cannot be empty", field))
	}
	return nil
}

// RequirePositive returns a ValidationError if v is not greater than zero.
func RequirePositive(field string, v int) error {
	if v <= 0 {
		return New(fmt.Sprintf("%s must be greater than zero", field))
	}
	return nil
}
