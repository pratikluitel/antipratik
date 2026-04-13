package logic

import (
	"errors"
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
