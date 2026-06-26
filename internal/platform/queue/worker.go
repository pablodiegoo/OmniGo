package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"

	"github.com/pablojhp.omnigo/internal/channel"
)

const (
	defaultMaxRetries  = 5
	defaultMaxBackoff  = 60 * time.Second
	defaultBaseBackoff = 1 * time.Second
)

// Worker reads messages from a JetStream consumer and dispatches them.
// It supports retry with exponential backoff, TTL enforcement, and
// delivery deduplication.
type Worker struct {
	consumer   jetstream.Consumer
	cancel     context.CancelFunc
	done       chan struct{}

	maxRetries int
	maxBackoff time.Duration

	// dispatchers provides channel-specific dispatchers.
	// When nil or a channel has no registered dispatcher, messages
	// are logged and acked (no-op).
	dispatchers *channel.Registry

	// dispatched tracks message IDs already processed (in-memory dedup).
	dispatched sync.Map // message_id string → struct{}
}

// NewWorker starts a goroutine that reads messages from consumer and processes them.
// The goroutine stops when ctx is cancelled. Call Stop() to initiate shutdown.
// dispatchers may be nil for log-only mode (useful in tests).
func NewWorker(ctx context.Context, consumer jetstream.Consumer, maxRetries int, maxBackoff time.Duration, dispatchers *channel.Registry) *Worker {
	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}
	if maxBackoff <= 0 {
		maxBackoff = defaultMaxBackoff
	}

	ctx, cancel := context.WithCancel(ctx)
	w := &Worker{
		consumer:   consumer,
		cancel:     cancel,
		done:       make(chan struct{}),
		maxRetries: maxRetries,
		maxBackoff: maxBackoff,
	}

	go w.run(ctx)
	return w
}

// run is the main consumer loop. It reads messages from the consumer,
// checks TTL and dedup, then dispatches. On failure, NAKs with
// exponential backoff up to maxRetries.
func (w *Worker) run(ctx context.Context) {
	defer close(w.done)

	msgCtx, err := w.consumer.Messages()
	if err != nil {
		slog.Error("worker: failed to create messages context", "error", err)
		return
	}
	defer msgCtx.Stop()

	slog.Info("message worker started",
		"consumer", w.consumer.CachedInfo().Config.Name,
		"max_retries", w.maxRetries,
		"max_backoff", w.maxBackoff,
	)

	for {
		msg, err := msgCtx.Next()
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("message worker stopped")
				return
			}
			slog.Error("worker: failed to get next message", "error", err)
			continue
		}

		w.processMessage(ctx, msg)
	}
}

// processMessage handles a single message: dedup check, TTL check, dispatch.
func (w *Worker) processMessage(ctx context.Context, msg jetstream.Msg) {
	// Extract trace_id from Nats-Msg-Id header
	traceID := ""
	if headers := msg.Headers(); headers != nil {
		traceID = headers.Get("Nats-Msg-Id")
	}

	// Extract retry attempt from header
	attempt := w.retryAttempt(msg)

	data := msg.Data()
	preview := string(data)
	if len(preview) > 100 {
		preview = preview[:100] + "..."
	}

	// --- Delivery dedup: skip if already processed ---
	if traceID != "" {
		if _, loaded := w.dispatched.LoadOrStore(traceID, struct{}{}); loaded {
			slog.Info("worker: duplicate delivery prevented",
				"trace_id", traceID,
				"subject", msg.Subject(),
			)
			_ = msg.Ack()
			return
		}
	}

	// --- TTL enforcement: drop expired messages ---
	if w.isExpired(msg) {
		slog.Warn("worker: message expired (TTL), dropping",
			"trace_id", traceID,
			"subject", msg.Subject(),
			"payload_preview", preview,
		)
		_ = msg.Ack()
		return
	}

	// --- Dispatch ---
	slog.Info("worker: dispatching message",
		"trace_id", traceID,
		"subject", msg.Subject(),
		"attempt", attempt,
	)

	dispatchErr := w.dispatch(ctx, msg)

	if dispatchErr != nil {
		w.handleFailure(msg, traceID, attempt)
		return
	}

	if err := msg.Ack(); err != nil {
		slog.Error("worker: failed to ack message", "error", err, "trace_id", traceID)
	}
}

// dispatch looks up the channel dispatcher from the payload and sends.
// Returns nil if no dispatchers are registered or channel is unknown
// (the message is effectively acked as a no-op — useful for testing).
func (w *Worker) dispatch(ctx context.Context, msg jetstream.Msg) error {
	if w.dispatchers == nil {
		return nil // no-op mode
	}

	// Extract channel from payload
	var payload struct {
		Channel string `json:"channel"`
		To      string `json:"to"`
		Body    string `json:"body"`
		TraceID string `json:"trace_id"`
	}
	if err := json.Unmarshal(msg.Data(), &payload); err != nil {
		return fmt.Errorf("worker: unmarshal payload: %w", err)
	}

	dispatcher, ok := w.dispatchers.Get(payload.Channel)
	if !ok {
		slog.Warn("worker: no dispatcher for channel, skipping",
			"channel", payload.Channel,
			"trace_id", payload.TraceID,
		)
		return nil // not an error — channel may not be configured yet
	}

	return dispatcher.Dispatch(ctx, &channel.MessagePayload{
		MessageID: msg.Headers().Get("Nats-Msg-Id"),
		TraceID:  payload.TraceID,
		To:       payload.To,
		Channel:  payload.Channel,
		Body:     payload.Body,
	})
}

// handleFailure NAKs the message with exponential backoff or marks as terminal failure.
func (w *Worker) handleFailure(msg jetstream.Msg, traceID string, attempt int) {
	if attempt >= w.maxRetries {
		// Terminal failure — ack and log
		slog.Error("worker: terminal failure after max retries",
			"trace_id", traceID,
			"subject", msg.Subject(),
			"attempts", attempt,
			"max_retries", w.maxRetries,
		)
		_ = msg.Ack()
		return
	}

	// Calculate exponential backoff: base * 2^attempt, capped at maxBackoff
	delay := time.Duration(math.Pow(2, float64(attempt))) * defaultBaseBackoff
	if delay > w.maxBackoff {
		delay = w.maxBackoff
	}

	slog.Warn("worker: retrying message with backoff",
		"trace_id", traceID,
		"attempt", attempt+1,
		"backoff", delay,
	)

	// NAK with delay for redelivery
	nakErr := msg.NakWithDelay(delay)
	if nakErr != nil {
		slog.Error("worker: failed to NAK message",
			"error", nakErr,
			"trace_id", traceID,
		)
		// If NAK fails, ack to prevent infinite loop
		_ = msg.Ack()
	}
}

// retryAttempt extracts the retry attempt count from message headers.
func (w *Worker) retryAttempt(msg jetstream.Msg) int {
	if headers := msg.Headers(); headers != nil {
		val := headers.Get("X-Retry-Attempt")
		if val != "" {
			var n int
			if _, err := fmt.Sscanf(val, "%d", &n); err == nil {
				return n
			}
		}
	}
	return 0
}

// isExpired checks if the message has exceeded its TTL.
// It looks for queued_at and ttl_seconds in the message metadata.
func (w *Worker) isExpired(msg jetstream.Msg) bool {
	data := msg.Data()

	// Parse minimal fields from the JSON payload to check TTL.
	// The payload is a CreateMessageRequest serialized to JSON.
	type ttlPayload struct {
		TTLSeconds *int `json:"ttl_seconds"`
		QueuedAt   string `json:"queued_at"`
	}

	var p ttlPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return false // can't parse TTL info, don't expire
	}

	if p.TTLSeconds == nil || *p.TTLSeconds <= 0 {
		return false // no TTL set
	}

	queuedAt, err := time.Parse(time.RFC3339Nano, p.QueuedAt)
	if err != nil {
		return false // can't parse queued_at, don't expire
	}

	expiry := queuedAt.Add(time.Duration(*p.TTLSeconds) * time.Second)
	return time.Now().UTC().After(expiry)
}



// Stop cancels the consumer context and waits for the goroutine to finish.
func (w *Worker) Stop() {
	w.cancel()
	<-w.done
}
