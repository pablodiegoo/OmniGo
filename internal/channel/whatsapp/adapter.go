package whatsapp

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"

	"github.com/pablojhp.omnigo/internal/channel"
)

// minStagger and maxStagger define the random delay range for
// ban-risk mitigation on WhatsApp Web sends.
const (
	minStagger = 1 * time.Second
	maxStagger = 3 * time.Second
)

// WhatsAppAdapter implements channel.Dispatcher for WhatsApp Web via whatsmeow.
// It adds staggered dispatch (1-3s random delay) before each send to
// minimize account suspension risk.
type WhatsAppAdapter struct {
	client *WhatsAppClient
	log    *slog.Logger
}

// NewWhatsAppAdapter creates a dispatcher backed by the given WhatsApp client.
func NewWhatsAppAdapter(client *WhatsAppClient) *WhatsAppAdapter {
	return &WhatsAppAdapter{
		client: client,
		log:    slog.With("component", "whatsapp-adapter"),
	}
}

// Dispatch sends a text message via WhatsApp Web with staggered delay.
// The recipient in payload.To should be a phone number (digits only or
// with country code). It's converted to a JID with @s.whatsapp.net suffix.
//
// Returns channel.TerminalError for 403/logged-out errors (non-retryable).
// Returns regular error for transient failures (retryable).
func (a *WhatsAppAdapter) Dispatch(ctx context.Context, m *channel.MessagePayload) error {
	if a.client == nil || a.client.Client() == nil {
		return fmt.Errorf("whatsapp: client not connected")
	}

	// Staggered dispatch: random delay before send
	stagger := time.Duration(rand.Int64N(int64(maxStagger-minStagger))) + minStagger
	select {
	case <-time.After(stagger):
	case <-ctx.Done():
		return ctx.Err()
	}

	// Convert phone number to JID
	recipientJID, parseErr := types.ParseJID(phoneToJID(m.To))
	if parseErr != nil {
		return fmt.Errorf("whatsapp: invalid recipient %q: %w", m.To, parseErr)
	}

	a.log.Info("whatsapp: dispatching message",
		"trace_id", m.TraceID,
		"to", recipientJID.String(),
		"stagger", stagger,
	)

	// Send text message via whatsmeow
	body := m.Body
	_, err := a.client.Client().SendMessage(ctx, recipientJID, &waE2E.Message{
		Conversation: &body,
	})
	if err != nil {
		a.log.Error("whatsapp: send failed",
			"error", err,
			"trace_id", m.TraceID,
			"to", recipientJID.String(),
		)
		if isTerminalWhatsAppError(err) {
			return channel.NewTerminalError(fmt.Errorf("whatsapp terminal: %w", err))
		}
		return fmt.Errorf("whatsapp send: %w", err)
	}

	a.log.Info("whatsapp: message sent",
		"trace_id", m.TraceID,
		"to", recipientJID.String(),
	)
	return nil
}

// phoneToJID converts a phone number string to a WhatsApp JID.
// Strips non-digit characters and appends @s.whatsapp.net.
func phoneToJID(phone string) string {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	return digits + "@s.whatsapp.net"
}

// isTerminalWhatsAppError classifies whatsmeow errors as terminal
// (non-retryable) vs transient.
func isTerminalWhatsAppError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "403") ||
		strings.Contains(msg, "logged out") ||
		strings.Contains(msg, "unpaired")
}
