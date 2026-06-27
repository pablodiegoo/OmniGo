---
phase: 6
slug: webhook-delivery-dlq
status: approved
nyquist_compliant: true
wave_0_complete: false
created: 2026-06-26
---

# Phase 6 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test (stdlib) |
| **Config file** | none |
| **Quick run command** | `go test ./...` |
| **Full suite command** | `go test ./... -race -count=1` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./...`
- **After every plan wave:** Run `go test ./... -race -count=1`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 06-01-01 | 01 | 1 | WHOOK-04 | — | N/A | integration | `go test ./internal/repository/... -run TestWebhookDLQRepository -v` | ❌ W0 | ⬜ pending |
| 06-01-02 | 01 | 2 | WHOOK-01 | — | N/A | integration | `go test ./internal/platform/queue/... -run TestPublishWebhookEvent -v` | ❌ W0 | ⬜ pending |
| 06-01-03 | 01 | 2 | WHOOK-02, WHOOK-03, WHOOK-05 | T-06-01 | HMAC signature verification | integration | `go test ./internal/platform/queue/... -run TestWebhookWorker -v` | ❌ W0 | ⬜ pending |
| 06-01-04 | 01 | 3 | WHOOK-04 | — | Workspace tenant isolation | integration | `go test ./cmd/omnigo/... -run TestAdminWebhookDLQHandlers -v` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/repository/webhook_dlq_test.go` — stubs for WHOOK-04 verification.
- [ ] `internal/platform/queue/webhook_worker_test.go` — stubs for WHOOK-01, WHOOK-02, WHOOK-03, WHOOK-05.
- [ ] `cmd/omnigo/admin_webhook_dlq_test.go` — stubs for WHOOK-04 UI API verification.

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Sidebar badge & HTMX list render | WHOOK-04 | Visual layout and badge rendering | Load dashboard, trigger dead-letter event, check sidebar badge updates via HTMX poll/refresh. |

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify or Wave 0 dependencies
- [x] Sampling continuity: no 3 consecutive tasks without automated verify
- [x] Wave 0 covers all MISSING references
- [x] No watch-mode flags
- [x] Feedback latency < 15s
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** approved 2026-06-26
