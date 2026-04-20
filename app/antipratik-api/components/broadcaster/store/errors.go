package store

import "errors"

// ErrDuplicate is returned when an INSERT would violate a UNIQUE constraint.
var ErrDuplicate = errors.New("duplicate entry")
