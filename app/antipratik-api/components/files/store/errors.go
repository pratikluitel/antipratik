package store

import "errors"

// ErrFileNotFound is returned by FileStore.Get when the requested file does not exist.
var ErrFileNotFound = errors.New("file not found")
