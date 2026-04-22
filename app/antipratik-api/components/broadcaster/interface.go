package broadcaster

import (
	"context"
	"net/http"
	"time"
)

// BroadcasterLogic defines all broadcaster business operations.
type BroadcasterLogic interface {
	// Subscriber lifecycle
	Subscribe(ctx context.Context, subType, address string) error
	SendConfirmationEmails(ctx context.Context, subType string) (int, error)
	ConfirmSubscription(ctx context.Context, token string) error
	Unsubscribe(ctx context.Context, token string) error
	GetSubscribers(ctx context.Context, subType string) ([]SubscriberSummary, error)

	// Broadcast management
	CreateBroadcast(ctx context.Context, input BroadcastInput) (BroadcastPreview, error)
	UpdateBroadcast(ctx context.Context, id int64, input BroadcastUpdateInput) (BroadcastPreview, error)
	GetBroadcasts(ctx context.Context, bType string) ([]BroadcastSummary, error)

	// Dispatch — returns the number of BUFFERED send records created.
	// Emails are delivered asynchronously in a background goroutine.
	DispatchBroadcast(ctx context.Context, id int64) (int, error)

	// Contact form
	SendContactMessage(ctx context.Context, input ContactInput) error
}

// BroadcasterStore handles all broadcaster database operations.
type BroadcasterStore interface {
	// Subscriber operations
	RegisterSubscriber(ctx context.Context, subType, address, token string) error
	ConfirmSubscriber(ctx context.Context, token string) error
	UnsubscribeByToken(ctx context.Context, token string) error
	SetSubscriberToken(ctx context.Context, address, token string) error
	GetConfirmedSubscribers(ctx context.Context, subType string) ([]Subscriber, error)
	GetUnconfirmedSubscribers(ctx context.Context, subType string) ([]Subscriber, error)
	GetAllSubscribers(ctx context.Context, subType string) ([]Subscriber, error)

	// Broadcast operations
	CreateBroadcast(ctx context.Context, input StoreBroadcastInput, emailBody string) (int64, error)
	UpdateBroadcast(ctx context.Context, id int64, input StoreBroadcastUpdateInput, emailBody string) error
	GetBroadcasts(ctx context.Context, bType string) ([]BroadcastRow, error)
	GetBroadcast(ctx context.Context, id int64) (BroadcastRow, error)

	// Dispatch operations
	CreateBroadcastSends(ctx context.Context, broadcastID int64, sends []BroadcastSendInput) error
	GetDueSends(ctx context.Context, broadcastID int64) ([]BroadcastSend, error)
	GetRemainingBuffered(ctx context.Context, broadcastID int64) (int, error)
	UpdateSendStatus(ctx context.Context, sendID int64, status, message string, sentAt *time.Time) error
	GetBroadcastSendSummary(ctx context.Context, broadcastID int64) (BroadcastSendSummary, error)

	// Contact messages
	SaveContactMessage(ctx context.Context, name, email, message string) error
}

// BroadcasterAPI is the HTTP handler interface for broadcaster endpoints.
type BroadcasterAPI interface {
	Subscribe(w http.ResponseWriter, r *http.Request)
	ResendConfirmation(w http.ResponseWriter, r *http.Request)
	Confirm(w http.ResponseWriter, r *http.Request)
	Unsubscribe(w http.ResponseWriter, r *http.Request)
	GetSubscribers(w http.ResponseWriter, r *http.Request)
	CreateBroadcast(w http.ResponseWriter, r *http.Request)
	UpdateBroadcast(w http.ResponseWriter, r *http.Request)
	GetBroadcasts(w http.ResponseWriter, r *http.Request)
	DispatchBroadcast(w http.ResponseWriter, r *http.Request)
	Contact(w http.ResponseWriter, r *http.Request)
}

// SubscriberService exposes subscription capability to other components.
// Inject this interface rather than importing broadcaster/logic directly.
type SubscriberService interface {
	Subscribe(ctx context.Context, subType, address string) error
}
