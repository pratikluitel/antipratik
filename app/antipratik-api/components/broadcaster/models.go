package broadcaster

import "time"

// ── Logic-layer models ────────────────────────────────────────────────────────

// BroadcastInput is the input for creating a broadcast.
type BroadcastInput struct {
	Type    string   // "email"
	Title   string   // email subject
	Caption string   // optional caption shown above post content
	PostIDs []string // post IDs to include, in display order
}

// BroadcastUpdateInput is the input for updating a broadcast's content.
// Type and ID cannot be changed.
type BroadcastUpdateInput struct {
	Title   string
	Caption string
	PostIDs []string
}

// BroadcastPreview is returned after create/update with the rendered HTML.
type BroadcastPreview struct {
	HTML string
	ID   int64
}

// SubscriberSummary is the API representation of a subscriber.
type SubscriberSummary struct {
	ConfirmedAt    *string `json:"confirmedAt,omitempty"`
	UnsubscribedAt *string `json:"unsubscribedAt,omitempty"`
	Type           string  `json:"type"`
	Address        string  `json:"address"`
	CreatedAt      string  `json:"createdAt"`
	Confirmed      bool    `json:"confirmed"`
}

// BroadcastSummary is a broadcast row with send counts and parsed data.
type BroadcastSummary struct {
	Type      string   `json:"type"`
	Title     string   `json:"title"`
	Caption   string   `json:"caption"`
	EmailBody string   `json:"emailBody"`
	PostIDs   []string `json:"postIDs"`
	ID        int64    `json:"id"`
	Buffered  int      `json:"buffered"`
	Success   int      `json:"success"`
	Failed    int      `json:"failed"`
}

// ContactInput is the input for the contact form.
type ContactInput struct {
	Name    string
	Email   string
	Message string
}

// ── Store-layer models ────────────────────────────────────────────────────────

// Subscriber represents a single subscriber row.
type Subscriber struct {
	CreatedAt      time.Time
	ConfirmedAt    *time.Time
	UnsubscribedAt *time.Time
	Type           string
	Address        string
	Token          string
	ID             int64
	Confirmed      bool
}

// BroadcastRow represents a broadcast with its optional rendered email body.
type BroadcastRow struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      string
	Title     string
	Data      string
	EmailBody string
	ID        int64
}

// StoreBroadcastInput is used when creating a broadcast in the store.
type StoreBroadcastInput struct {
	Type  string
	Title string
	Data  string // JSON
}

// StoreBroadcastUpdateInput is used when updating a broadcast in the store.
type StoreBroadcastUpdateInput struct {
	Title string
	Data  string // JSON
}

// BroadcastSendInput is one row to insert into broadcast_sends.
type BroadcastSendInput struct {
	ScheduledAt  time.Time
	SubscriberID int64
}

// BroadcastSend is a single row from broadcast_sends joined with subscriber address/token.
type BroadcastSend struct {
	ScheduledAt       time.Time
	SentAt            *time.Time
	SubscriberAddress string
	SubscriberToken   string
	Status            string
	Message           string
	ID                int64
	BroadcastID       int64
	SubscriberID      int64
}

// BroadcastSendSummary holds aggregate counts for a broadcast's sends by status.
type BroadcastSendSummary struct {
	Buffered int
	Success  int
	Failed   int
}

// BroadcastSendDetail is a single send row returned by the admin sends endpoint.
type BroadcastSendDetail struct {
	Address     string  `json:"address"`
	Status      string  `json:"status"`
	Message     string  `json:"message,omitempty"`
	ScheduledAt string  `json:"scheduledAt"`
	SentAt      *string `json:"sentAt,omitempty"`
}
