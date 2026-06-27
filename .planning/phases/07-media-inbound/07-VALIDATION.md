---
phase: 7
slug: media-inbound
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-06-27
---

# Phase 7 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none |
| **Quick run command** | `go test ./internal/domain/... ./internal/channel/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/domain/... ./internal/channel/...` (or local package tests)
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 15 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 07-01-01 | 01 | 1 | MEDIA-01 | — | N/A | unit | `go test -run TestCreateMessageRequest_Validate ./internal/domain` | ❌ W0 | ⬜ pending |
| 07-01-02 | 01 | 1 | MEDIA-03 | T-07-01 | Rejects size > 25MB, download timeouts and failures early | unit | `go test -run TestDownloadAndValidate ./internal/platform/media` | ❌ W0 | ⬜ pending |
| 07-01-03 | 01 | 1 | MEDIA-03 | — | N/A | integration | `go test -run TestS3Client ./internal/platform/media` | ❌ W0 | ⬜ pending |
| 07-02-01 | 01 | 2 | MEDIA-02 | — | N/A | integration | `go test -run TestWhatsAppAdapter_Media ./internal/channel/whatsapp` | ❌ W0 | ⬜ pending |
| 07-02-02 | 01 | 2 | MEDIA-02 | — | N/A | integration | `go test -run TestTelegramAdapter_Media ./internal/channel/telegram` | ❌ W0 | ⬜ pending |
| 07-02-03 | 01 | 2 | MEDIA-02 | — | N/A | integration | `go test -run TestWABAAdapter_Media ./internal/channel/whatsapp` | ❌ W0 | ⬜ pending |
| 07-03-01 | 01 | 3 | INBD-01 | — | N/A | unit | `go test -run TestInboundDeduplicate ./internal/repository` | ❌ W0 | ⬜ pending |
| 07-03-02 | 01 | 3 | INBD-01 | — | N/A | integration | `go test -run TestTelegramWebhook_Inbound ./internal/api/handler` | ❌ W0 | ⬜ pending |
| 07-03-03 | 01 | 3 | INBD-01 | — | N/A | integration | `go test -run TestWABAWebhook_Inbound ./internal/api/handler` | ❌ W0 | ⬜ pending |
| 07-03-04 | 01 | 3 | INBD-01 | — | N/A | integration | `go test -run TestWhatsAppInbound ./internal/session` | ❌ W0 | ⬜ pending |
| 07-04-01 | 01 | 4 | INBD-02 | — | N/A | integration | `go test -run TestWebhookWorker_Inbound ./internal/platform/queue` | ❌ W0 | ⬜ pending |
| 07-04-02 | 01 | 4 | INBD-03 | — | N/A | integration | `go test -run TestInboundAuditLogging ./internal/api/handler` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/platform/media/media_test.go` — stubs for media download & validation
- [ ] `internal/api/handler/waba_webhook_test.go` — stubs for WABA inbound verification/payload test
- [ ] `internal/repository/inbound_dedup_test.go` — stub for DB-level inbound deduplication test

*If none: "Existing infrastructure covers all phase requirements."*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Webhook PII Opt-In payload check | MEDIA-01, INBD-02 | Verifies compliance flow visually & logically | Inspect webhook logs under workspaces with pii_opt_in=false and pii_opt_in=true |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 15s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
