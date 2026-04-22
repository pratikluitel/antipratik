package store

import "errors"

// ErrDuplicate is returned when an INSERT would violate a UNIQUE constraint.
var ErrDuplicate = errors.New("duplicate entry")

// ErrNotFound is returned when a requested record does not exist.
var ErrNotFound = errors.New("not found")
