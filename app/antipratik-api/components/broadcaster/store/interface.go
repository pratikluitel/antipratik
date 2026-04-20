// Package store contains the broadcaster data persistence layer.
package store

import "context"

// NewsletterStore handles newsletter subscriber database operations.
type NewsletterStore interface {
	// Subscribe inserts an email into newsletter_subscribers.
	// Returns ErrDuplicate if the email already exists.
	Subscribe(ctx context.Context, email string) error
}
