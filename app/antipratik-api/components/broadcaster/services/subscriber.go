// Package services exposes broadcaster capabilities as injectable interfaces
// for use by other components. Consumers depend on the interface, not the
// concrete broadcaster logic, keeping components decoupled.
package services

import (
	"context"

	"github.com/pratikluitel/antipratik/components/broadcaster"
)

type subscriberService struct {
	logic broadcaster.BroadcasterLogic
}

// NewSubscriberService returns a SubscriberService backed by the given BroadcasterLogic.
func NewSubscriberService(l broadcaster.BroadcasterLogic) broadcaster.SubscriberService {
	return &subscriberService{logic: l}
}

func (s *subscriberService) Subscribe(ctx context.Context, subType, address string) error {
	return s.logic.Subscribe(ctx, subType, address)
}
