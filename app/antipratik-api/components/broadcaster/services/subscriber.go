// Package services exposes broadcaster capabilities as injectable interfaces
// for use by other components. Consumers depend on the interface, not the
// concrete broadcaster logic, keeping components decoupled.
package services

import (
	"context"

	"github.com/pratikluitel/antipratik/components/broadcaster/logic"
)

// SubscriberService exposes the broadcaster's subscription capability to other
// components. Inject this interface rather than importing broadcaster/logic directly.
type SubscriberService interface {
	Subscribe(ctx context.Context, email string) error
}

type subscriberService struct {
	logic logic.NewsletterLogic
}

// NewSubscriberService returns a SubscriberService backed by the given NewsletterLogic.
func NewSubscriberService(l logic.NewsletterLogic) SubscriberService {
	return &subscriberService{logic: l}
}

func (s *subscriberService) Subscribe(ctx context.Context, email string) error {
	return s.logic.Subscribe(ctx, email)
}
