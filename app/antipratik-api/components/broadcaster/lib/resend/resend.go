// Package resend provides an SMTP client for Resend email delivery.
// This package has no dependencies on other antipratik packages.
package resend

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

// Logger is a minimal logging interface injected by callers.
type Logger interface {
	Error(msg string, args ...any)
}

// Config holds Resend SMTP connection settings.
// Fields ordered by decreasing alignment for minimal padding.
type Config struct {
	APIKey   string // Resend API key — used as SMTP password
	Host     string // Default: smtp.resend.com
	From     string // Sender email address
	FromName string // Sender display name
	Port     int    // Default: 587 (STARTTLS)
}

// ErrTransient wraps failures that are safe to retry.
type ErrTransient struct{ Cause error }

func (e ErrTransient) Error() string { return "transient: " + e.Cause.Error() }
func (e ErrTransient) Unwrap() error { return e.Cause }

// ErrPermanent wraps failures that must not be retried.
type ErrPermanent struct{ Cause error }

func (e ErrPermanent) Error() string { return "permanent: " + e.Cause.Error() }
func (e ErrPermanent) Unwrap() error { return e.Cause }

// IsTransient returns true if err is retry-able.
func IsTransient(err error) bool {
	var t ErrTransient
	return errors.As(err, &t)
}

// SendRequest is the input for sending an HTML email.
// Fields ordered by decreasing alignment.
type SendRequest struct {
	Subject string
	HTML    string
	To      []string
}

// Client sends emails via Resend SMTP.
// Fields ordered by decreasing alignment.
type Client struct {
	logger Logger
	cfg    Config
}

// NewClient creates a new Client. Applies defaults for Host (smtp.resend.com) and Port (587).
func NewClient(cfg Config, logger Logger) *Client {
	if cfg.Host == "" {
		cfg.Host = "smtp.resend.com"
	}
	if cfg.Port == 0 {
		cfg.Port = 587
	}
	return &Client{cfg: cfg, logger: logger}
}

// Send delivers an HTML email via Resend's SMTP interface (STARTTLS on port 587).
func (c *Client) Send(_ context.Context, req SendRequest) error {
	if len(req.To) == 0 {
		return ErrPermanent{Cause: fmt.Errorf("no recipients")}
	}

	from := fmt.Sprintf("%s <%s>", c.cfg.FromName, c.cfg.From)
	var buf strings.Builder
	buf.WriteString("From: " + from + "\r\n")
	buf.WriteString("To: " + strings.Join(req.To, ", ") + "\r\n")
	buf.WriteString("Subject: " + req.Subject + "\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(req.HTML)

	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
	tlsCfg := &tls.Config{
		ServerName: c.cfg.Host,
		MinVersion: tls.VersionTLS12,
	}

	client, dialErr := smtp.Dial(addr)
	if dialErr != nil {
		return ErrTransient{Cause: fmt.Errorf("dial %s: %w", addr, dialErr)}
	}
	defer func() { _ = client.Quit() }()

	if startTLSErr := client.StartTLS(tlsCfg); startTLSErr != nil {
		return ErrTransient{Cause: fmt.Errorf("starttls: %w", startTLSErr)}
	}

	auth := smtp.PlainAuth("", "resend", c.cfg.APIKey, c.cfg.Host)
	if authErr := client.Auth(auth); authErr != nil {
		return ErrPermanent{Cause: fmt.Errorf("smtp auth: %w", authErr)}
	}

	if mailErr := client.Mail(c.cfg.From); mailErr != nil {
		return ErrPermanent{Cause: fmt.Errorf("smtp MAIL FROM: %w", mailErr)}
	}

	for _, to := range req.To {
		if rcptErr := client.Rcpt(to); rcptErr != nil {
			return ErrPermanent{Cause: fmt.Errorf("smtp RCPT TO %s: %w", to, rcptErr)}
		}
	}

	w, dataErr := client.Data()
	if dataErr != nil {
		return ErrTransient{Cause: fmt.Errorf("smtp DATA: %w", dataErr)}
	}

	if _, writeErr := fmt.Fprint(w, buf.String()); writeErr != nil {
		return ErrTransient{Cause: fmt.Errorf("smtp write: %w", writeErr)}
	}

	if closeErr := w.Close(); closeErr != nil {
		return ErrTransient{Cause: fmt.Errorf("smtp close: %w", closeErr)}
	}

	return nil
}
