# Milestones

## v1.0 v1.0 (Shipped: 2026-06-27)

**Phases completed:** 8 phases, 23 plans, 32 tasks

**Key accomplishments:**

- Echo v5 server scaffold with pgxpool dual-access PostgreSQL, goose embedded migrations, health/readiness endpoints, and Docker Compose topology
- SHA-256 API key hashing with prefix lookup, AES-256-GCM envelope encryption, tenant-context convention, workspace/API key repositories with in-memory cache, and Echo v5 auth middleware
- Trace-ID middleware with UUID generation and header extraction, structured slog with trace context, buffered batch audit writer with pgx.CopyFrom, and monthly partitioned audit_logs with BRIN index
- pprof debug server on localhost:6060, expvar metrics with audit_drops counter, and LIFO shutdown orchestrator wiring Echo → debug → audit → NATS → pgxpool → sqlDB
- Server-rendered admin panel with templ + HTMX, HMAC-signed session auth, sidebar navigation, login/logout flow, and dashboard landing page with workspace count and audit activity
- Workspace CRUD and API key lifecycle management via admin panel with HTMX fragment updates, modal confirmations, and one-time plaintext key display
- Audit log review with parameterized filtering (workspace, trace_id, event_type, time range), 50-row pagination, CSV export, and HTMX fragment updates via admin panel
- POST /messages with JSON validation, structured error responses, and X-Trace-Id correlation header
- WorkQueue stream with publish-side dedup, worker stub, and composition root lifecycle wiring
- Production-hardens the message ingestion path with per-workspace rate limiting, queue depth backpressure, worker retry with exponential backoff, TTL enforcement, and delivery deduplication.
- Channel abstraction layer and WhatsApp Web adapter stub.
- Session management layer — device persistence, in-memory registry, and startup reconnection with storm protection.
- Admin UI for QR pairing and session telemetry — completes the operator-facing WhatsApp Web experience.

---
