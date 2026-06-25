// Package domain defines the core message types, validation rules, and
// error contracts for the OmniGo message ingestion API.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// MessageStatus represents the lifecycle state of a message.
type MessageStatus string

const (
	StatusQueued    MessageStatus = "queued"
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
	StatusFailed    MessageStatus = "failed"
)

// ValidChannels defines the set of accepted channel values.
var ValidChannels = map[string]bool{
	"whatsapp":      true,
	"whatsapp_cloud": true,
	"telegram":      true,
}

// CreateMessageRequest is the JSON payload for POST /messages.
type CreateMessageRequest struct {
	To         string            `json:"to"`
	Channel    string            `json:"channel"`
	Body       string            `json:"body"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	TTLSeconds *int              `json:"ttl_seconds,omitempty"`
}

// CreateMessageResponse is returned on successful enqueue (HTTP 202).
type CreateMessageResponse struct {
	MessageID uuid.UUID     `json:"message_id"`
	Status    MessageStatus `json:"status"`
	QueuedAt  time.Time     `json:"queued_at"`
}

// ErrorResponse is the structured error format per API-04.
type ErrorResponse struct {
	Code     string       `json:"code"`
	Message  string       `json:"message"`
	MoreInfo string       `json:"more_info,omitempty"`
	Details  []FieldError `json:"details,omitempty"`
}

// FieldError provides field-level validation error details.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateMessage checks a CreateMessageRequest for correctness and returns
// an ErrorResponse if any validation rule fails. Returns nil when valid.
func ValidateMessage(req *CreateMessageRequest) *ErrorResponse {
	var details []FieldError

	if req.To == "" {
		details = append(details, FieldError{
			Field:   "to",
			Message: "is required",
		})
	}

	if req.Channel == "" {
		details = append(details, FieldError{
			Field:   "channel",
			Message: "is required",
		})
	} else if !ValidChannels[req.Channel] {
		details = append(details, FieldError{
			Field:   "channel",
			Message: "must be one of: whatsapp, whatsapp_cloud, telegram",
		})
	}

	if req.TTLSeconds != nil && *req.TTLSeconds <= 0 {
		details = append(details, FieldError{
			Field:   "ttl_seconds",
			Message: "must be a positive integer",
		})
	}

	if len(details) > 0 {
		return &ErrorResponse{
			Code:    "invalid_payload",
			Message: "request body validation failed",
			Details: details,
		}
	}

	return nil
}
