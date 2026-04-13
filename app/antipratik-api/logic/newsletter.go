package logic

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/pratikluitel/antipratik/store"
)

// emailRegex requires at least one non-whitespace/@ char on each side of @,
// and a dot somewhere after the @.
var emailRegex = regexp.MustCompile(`(?i)^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// NewsletterService implements NewsletterLogic.
type NewsletterService struct {
	store store.NewsletterStore
}

// NewNewsletterService creates a new NewsletterService backed by the given store.
func NewNewsletterService(s store.NewsletterStore) *NewsletterService {
	return &NewsletterService{store: s}
}

// Subscribe validates the email format, checks that the domain has MX records,
// and persists the subscriber. Returns a ValidationError for bad input or
// duplicates, and a wrapped store error for unexpected failures.
func (s *NewsletterService) Subscribe(ctx context.Context, email string) error {
	email = strings.TrimSpace(email)

	if !emailRegex.MatchString(email) {
		return validationErr("please enter a valid email address")
	}

	// Extract domain from email for MX lookup.
	parts := strings.SplitN(email, "@", 2)
	domain := parts[1]

	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return validationErr("email domain cannot receive mail")
	}

	if err := s.store.Subscribe(ctx, email); err != nil {
		if errors.Is(err, store.ErrDuplicate) {
			return validationErr("already subscribed")
		}
		return fmt.Errorf("NewsletterService.Subscribe: %w", err)
	}
	return nil
}
