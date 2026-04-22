package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pratikluitel/antipratik/components/broadcaster"
)

// sqliteBroadcasterStore is the SQLite implementation of BroadcasterStore.
type sqliteBroadcasterStore struct {
	db *sql.DB
}

// NewBroadcasterStore creates a new sqliteBroadcasterStore.
func NewBroadcasterStore(db *sql.DB) broadcaster.BroadcasterStore {
	return &sqliteBroadcasterStore{db: db}
}

// ── Subscriber operations ─────────────────────────────────────────────────────

// RegisterSubscriber inserts a new subscriber.
// Returns ErrDuplicate if the address is already registered.
func (s *sqliteBroadcasterStore) RegisterSubscriber(ctx context.Context, subType, address, token string) error {
	address = strings.ToLower(strings.TrimSpace(address))
	res, err := s.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO subscribers (type, address, token) VALUES (?, ?, ?)`,
		subType, address, token)
	if err != nil {
		return fmt.Errorf("RegisterSubscriber: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RegisterSubscriber rows: %w", err)
	}
	if n == 0 {
		return ErrDuplicate
	}
	return nil
}

// ConfirmSubscriber marks the subscriber with the given token as confirmed.
// Returns ErrNotFound if no such token exists.
func (s *sqliteBroadcasterStore) ConfirmSubscriber(ctx context.Context, token string) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE subscribers SET confirmed = TRUE, confirmed_at = CURRENT_TIMESTAMP WHERE token = ?`, token)
	if err != nil {
		return fmt.Errorf("ConfirmSubscriber: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ConfirmSubscriber rows: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// UnsubscribeByToken sets unsubscribed_at for the subscriber with the given token.
// Returns ErrNotFound if no such token exists.
func (s *sqliteBroadcasterStore) UnsubscribeByToken(ctx context.Context, token string) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE subscribers SET unsubscribed_at = CURRENT_TIMESTAMP WHERE token = ? AND unsubscribed_at IS NULL`, token)
	if err != nil {
		return fmt.Errorf("UnsubscribeByToken: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UnsubscribeByToken rows: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// SetSubscriberToken updates the token for an existing subscriber by address.
func (s *sqliteBroadcasterStore) SetSubscriberToken(ctx context.Context, address, token string) error {
	address = strings.ToLower(strings.TrimSpace(address))
	_, err := s.db.ExecContext(ctx,
		`UPDATE subscribers SET token = ? WHERE address = ?`, token, address)
	if err != nil {
		return fmt.Errorf("SetSubscriberToken: %w", err)
	}
	return nil
}

// GetConfirmedSubscribers returns all confirmed, active (not unsubscribed) subscribers of the given type.
func (s *sqliteBroadcasterStore) GetConfirmedSubscribers(ctx context.Context, subType string) ([]broadcaster.Subscriber, error) {
	return s.querySubscribers(ctx,
		`SELECT id, type, address, token, confirmed, created_at, confirmed_at, unsubscribed_at
		 FROM subscribers
		 WHERE type = ? AND confirmed = TRUE AND unsubscribed_at IS NULL`, subType)
}

// GetUnconfirmedSubscribers returns all unconfirmed, active subscribers of the given type.
func (s *sqliteBroadcasterStore) GetUnconfirmedSubscribers(ctx context.Context, subType string) ([]broadcaster.Subscriber, error) {
	return s.querySubscribers(ctx,
		`SELECT id, type, address, token, confirmed, created_at, confirmed_at, unsubscribed_at
		 FROM subscribers
		 WHERE type = ? AND confirmed = FALSE AND unsubscribed_at IS NULL`, subType)
}

// DeleteSubscriber hard-deletes the subscriber row with the given address.
// Returns ErrNotFound if no row matches.
func (s *sqliteBroadcasterStore) DeleteSubscriber(ctx context.Context, address string) error {
	address = strings.ToLower(strings.TrimSpace(address))
	res, err := s.db.ExecContext(ctx, `DELETE FROM subscribers WHERE address = ?`, address)
	if err != nil {
		return fmt.Errorf("DeleteSubscriber: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteSubscriber rows: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// GetAllSubscribers returns all subscribers of the given type, including inactive ones.
func (s *sqliteBroadcasterStore) GetAllSubscribers(ctx context.Context, subType string) ([]broadcaster.Subscriber, error) {
	return s.querySubscribers(ctx,
		`SELECT id, type, address, token, confirmed, created_at, confirmed_at, unsubscribed_at
		 FROM subscribers
		 WHERE type = ?
		 ORDER BY created_at DESC`, subType)
}

func (s *sqliteBroadcasterStore) querySubscribers(ctx context.Context, query, arg string) ([]broadcaster.Subscriber, error) {
	rows, err := s.db.QueryContext(ctx, query, arg)
	if err != nil {
		return nil, fmt.Errorf("querySubscribers: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []broadcaster.Subscriber
	for rows.Next() {
		var sub broadcaster.Subscriber
		var confirmedAt, unsubscribedAt sql.NullTime
		if err := rows.Scan(
			&sub.ID, &sub.Type, &sub.Address, &sub.Token, &sub.Confirmed,
			&sub.CreatedAt, &confirmedAt, &unsubscribedAt,
		); err != nil {
			return nil, fmt.Errorf("querySubscribers scan: %w", err)
		}
		if confirmedAt.Valid {
			sub.ConfirmedAt = &confirmedAt.Time
		}
		if unsubscribedAt.Valid {
			sub.UnsubscribedAt = &unsubscribedAt.Time
		}
		out = append(out, sub)
	}
	return out, rows.Err()
}

// ── Broadcast operations ──────────────────────────────────────────────────────

// CreateBroadcast inserts a broadcast and its rendered email body atomically.
func (s *sqliteBroadcasterStore) CreateBroadcast(ctx context.Context, input broadcaster.StoreBroadcastInput, emailBody string) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("CreateBroadcast begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx,
		`INSERT INTO broadcasts (type, title, data) VALUES (?, ?, ?)`,
		input.Type, input.Title, input.Data)
	if err != nil {
		return 0, fmt.Errorf("CreateBroadcast insert: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateBroadcast last id: %w", err)
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO email_broadcasts (broadcast_id, email_body) VALUES (?, ?)`, id, emailBody); err != nil {
		return 0, fmt.Errorf("CreateBroadcast email body: %w", err)
	}

	return id, tx.Commit()
}

// UpdateBroadcast updates a broadcast's title and data, and re-saves the email body.
func (s *sqliteBroadcasterStore) UpdateBroadcast(ctx context.Context, id int64, input broadcaster.StoreBroadcastUpdateInput, emailBody string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("UpdateBroadcast begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx,
		`UPDATE broadcasts SET title = ?, data = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		input.Title, input.Data, id); err != nil {
		return fmt.Errorf("UpdateBroadcast update: %w", err)
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO email_broadcasts (broadcast_id, email_body) VALUES (?, ?)
		 ON CONFLICT(broadcast_id) DO UPDATE SET email_body = excluded.email_body`,
		id, emailBody); err != nil {
		return fmt.Errorf("UpdateBroadcast email body: %w", err)
	}

	return tx.Commit()
}

// GetBroadcasts returns all broadcasts of the given type, with their email bodies.
func (s *sqliteBroadcasterStore) GetBroadcasts(ctx context.Context, bType string) ([]broadcaster.BroadcastRow, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT b.id, b.type, b.title, b.data, COALESCE(eb.email_body, ''), b.created_at, b.updated_at
		FROM broadcasts b
		LEFT JOIN email_broadcasts eb ON eb.broadcast_id = b.id
		WHERE b.type = ?
		ORDER BY b.created_at DESC`, bType)
	if err != nil {
		return nil, fmt.Errorf("GetBroadcasts: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []broadcaster.BroadcastRow
	for rows.Next() {
		var r broadcaster.BroadcastRow
		if err := rows.Scan(&r.ID, &r.Type, &r.Title, &r.Data, &r.EmailBody, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("GetBroadcasts scan: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// GetBroadcast returns a single broadcast by ID.
// Returns ErrNotFound if the broadcast does not exist.
func (s *sqliteBroadcasterStore) GetBroadcast(ctx context.Context, id int64) (broadcaster.BroadcastRow, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT b.id, b.type, b.title, b.data, COALESCE(eb.email_body, ''), b.created_at, b.updated_at
		FROM broadcasts b
		LEFT JOIN email_broadcasts eb ON eb.broadcast_id = b.id
		WHERE b.id = ?`, id)

	var r broadcaster.BroadcastRow
	if err := row.Scan(&r.ID, &r.Type, &r.Title, &r.Data, &r.EmailBody, &r.CreatedAt, &r.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return broadcaster.BroadcastRow{}, ErrNotFound
		}
		return broadcaster.BroadcastRow{}, fmt.Errorf("GetBroadcast: %w", err)
	}
	return r, nil
}

// ── Dispatch operations ───────────────────────────────────────────────────────

// CreateBroadcastSends bulk-inserts broadcast_sends rows in a single transaction.
func (s *sqliteBroadcasterStore) CreateBroadcastSends(ctx context.Context, broadcastID int64, sends []broadcaster.BroadcastSendInput) error {
	if len(sends) == 0 {
		return nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("CreateBroadcastSends begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT OR IGNORE INTO broadcast_sends (broadcast_id, subscriber_id, scheduled_at) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("CreateBroadcastSends prepare: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	for _, send := range sends {
		// SQLite CURRENT_TIMESTAMP evaluates to UTC 'YYYY-MM-DD HH:MM:SS'
		schedStr := send.ScheduledAt.UTC().Format("2006-01-02 15:04:05")
		if _, err := stmt.ExecContext(ctx, broadcastID, send.SubscriberID, schedStr); err != nil {
			return fmt.Errorf("CreateBroadcastSends exec: %w", err)
		}
	}
	return tx.Commit()
}

// GetDueSends returns all BUFFERED sends for the broadcast whose scheduled_at is in the past.
func (s *sqliteBroadcasterStore) GetDueSends(ctx context.Context, broadcastID int64) ([]broadcaster.BroadcastSend, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT bs.id, bs.broadcast_id, bs.subscriber_id,
		       sub.address, sub.token,
		       bs.status, COALESCE(bs.message, ''), bs.scheduled_at, bs.sent_at
		FROM broadcast_sends bs
		JOIN subscribers sub ON sub.id = bs.subscriber_id
		WHERE bs.broadcast_id = ?
		  AND bs.status = 'BUFFERED'
		  AND bs.scheduled_at <= CURRENT_TIMESTAMP`, broadcastID)
	if err != nil {
		return nil, fmt.Errorf("GetDueSends: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []broadcaster.BroadcastSend
	for rows.Next() {
		var send broadcaster.BroadcastSend
		var sentAt sql.NullTime
		if err := rows.Scan(
			&send.ID, &send.BroadcastID, &send.SubscriberID,
			&send.SubscriberAddress, &send.SubscriberToken,
			&send.Status, &send.Message, &send.ScheduledAt, &sentAt,
		); err != nil {
			return nil, fmt.Errorf("GetDueSends scan: %w", err)
		}
		if sentAt.Valid {
			send.SentAt = &sentAt.Time
		}
		out = append(out, send)
	}
	return out, rows.Err()
}

// GetRemainingBuffered returns the count of BUFFERED sends remaining for the broadcast.
func (s *sqliteBroadcasterStore) GetRemainingBuffered(ctx context.Context, broadcastID int64) (int, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM broadcast_sends WHERE broadcast_id = ? AND status = 'BUFFERED'`, broadcastID)
	var n int
	if err := row.Scan(&n); err != nil {
		return 0, fmt.Errorf("GetRemainingBuffered: %w", err)
	}
	return n, nil
}

// UpdateSendStatus updates the status (and optionally sent_at) for a single broadcast_send row.
func (s *sqliteBroadcasterStore) UpdateSendStatus(ctx context.Context, sendID int64, status, message string, sentAt *time.Time) error {
	if sentAt != nil {
		_, err := s.db.ExecContext(ctx,
			`UPDATE broadcast_sends SET status = ?, message = ?, sent_at = ? WHERE id = ?`,
			status, message, sentAt, sendID)
		return err
	}
	_, err := s.db.ExecContext(ctx,
		`UPDATE broadcast_sends SET status = ?, message = ? WHERE id = ?`,
		status, message, sendID)
	return err
}

// GetAllBroadcastSends returns every send row for a broadcast with subscriber address, ordered by status then address.
func (s *sqliteBroadcasterStore) GetAllBroadcastSends(ctx context.Context, broadcastID int64) ([]broadcaster.BroadcastSend, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT bs.id, bs.broadcast_id, bs.subscriber_id,
		       sub.address, sub.token,
		       bs.status, COALESCE(bs.message, ''), bs.scheduled_at, bs.sent_at
		FROM broadcast_sends bs
		JOIN subscribers sub ON sub.id = bs.subscriber_id
		WHERE bs.broadcast_id = ?
		ORDER BY bs.status, sub.address`, broadcastID)
	if err != nil {
		return nil, fmt.Errorf("GetAllBroadcastSends: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []broadcaster.BroadcastSend
	for rows.Next() {
		var send broadcaster.BroadcastSend
		var sentAt sql.NullTime
		if err := rows.Scan(
			&send.ID, &send.BroadcastID, &send.SubscriberID,
			&send.SubscriberAddress, &send.SubscriberToken,
			&send.Status, &send.Message, &send.ScheduledAt, &sentAt,
		); err != nil {
			return nil, fmt.Errorf("GetAllBroadcastSends scan: %w", err)
		}
		if sentAt.Valid {
			send.SentAt = &sentAt.Time
		}
		out = append(out, send)
	}
	return out, rows.Err()
}

// GetBroadcastSendSummary returns aggregate counts grouped by status for a broadcast.
func (s *sqliteBroadcasterStore) GetBroadcastSendSummary(ctx context.Context, broadcastID int64) (broadcaster.BroadcastSendSummary, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT status, COUNT(*) FROM broadcast_sends WHERE broadcast_id = ? GROUP BY status`, broadcastID)
	if err != nil {
		return broadcaster.BroadcastSendSummary{}, fmt.Errorf("GetBroadcastSendSummary: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var summary broadcaster.BroadcastSendSummary
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return broadcaster.BroadcastSendSummary{}, fmt.Errorf("GetBroadcastSendSummary scan: %w", err)
		}
		switch status {
		case "BUFFERED":
			summary.Buffered = count
		case "SUCCESS":
			summary.Success = count
		case "FAILED":
			summary.Failed = count
		}
	}
	return summary, rows.Err()
}

// SaveContactMessage inserts an inbound contact form message.
func (s *sqliteBroadcasterStore) SaveContactMessage(ctx context.Context, name, email, message string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO contact_messages (name, email, message) VALUES (?, ?, ?)`, name, email, message)
	if err != nil {
		return fmt.Errorf("SaveContactMessage: %w", err)
	}
	return nil
}
