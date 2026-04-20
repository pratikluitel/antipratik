// Package logic contains the broadcaster business logic layer.
package logic

import "context"

// NewsletterLogic defines newsletter subscription operations.
type NewsletterLogic interface {
	// Subscribe validates the email and persists it as a subscriber.
	Subscribe(ctx context.Context, email string) error
}
