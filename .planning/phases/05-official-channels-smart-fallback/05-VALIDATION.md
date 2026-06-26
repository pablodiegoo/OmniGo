---
phase: 5
slug: official-channels-smart-fallback
status: draft
nyquist_compliant: true
wave_0_complete: false
created: 2026-06-26
---

# Phase 5 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — stdlib testing |
| **Quick run command** | `go test -v ./internal/channel/... ./internal/session/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test -v ./internal/channel/... ./internal/session/...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 05-01-01 | 01 | 1 | WABA-01 | — | REST adapters for WABA and Telegram implement Dispatcher interface | unit | `go test -v ./internal/channel/...` | ✅ | ⬜ pending |
| 05-01-02 | 01 | 1 | TGRAM-01 | — | WABA and Telegram token encryption at rest | unit | `go test -v ./internal/channel/...` | ✅ | ⬜ pending |
| 05-02-01 | 02 | 2 | WABA-02 | — | Template CRUD database repository and endpoints | unit | `go test -v ./internal/session/...` | ✅ | ⬜ pending |
| 05-02-02 | 02 | 2 | WABA-03 | — | WABA template message sending support | unit | `go test -v ./internal/channel/whatsapp/...` | ✅ | ⬜ pending |
| 05-03-01 | 03 | 3 | TGRAM-02 | — | Telegram Bot webhook secret token validation | unit | `go test -v ./internal/api/handler/...` | ✅ | ⬜ pending |
| 05-03-02 | 03 | 3 | WABA-04 | — | 24-hour customer window checking logic | unit | `go test -v ./internal/session/...` | ✅ | ⬜ pending |
| 05-04-01 | 04 | 4 | FALL-01 | — | Fallback channels array iterative dispatch in worker | unit | `go test -v ./internal/platform/queue/...` | ✅ | ⬜ pending |
| 05-04-02 | 04 | 4 | FALL-02 | — | Terminal error bypasses queue retry to fallback | unit | `go test -v ./internal/platform/queue/...` | ✅ | ⬜ pending |
| 05-04-03 | 04 | 4 | FALL-03 | — | Fallback-aware deduplication via message_dispatches table | unit | `go test -v ./internal/platform/queue/...` | ✅ | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/channel/whatsapp/waba_test.go` — stubs for WABA-01 and WABA-03
- [ ] `internal/platform/queue/fallback_test.go` — stubs for FALL-01, FALL-02, and FALL-03

*If none: "Existing infrastructure covers all phase requirements."*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| WABA webhook handshake verification | WABA-02 | Meta webhook subscription verification | Mock Meta subscription GET request to webhooks and verify challenge response |
| Telegram webhook end-to-end webhook validation | TGRAM-02 | External webhook integration | Send simulated webhook POST request with Telegram signature and verify response |

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify or Wave 0 dependencies
- [x] Sampling continuity: no 3 consecutive tasks without automated verify
- [x] Wave 0 covers all MISSING references
- [x] No watch-mode flags
- [x] Feedback latency < 5s
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
