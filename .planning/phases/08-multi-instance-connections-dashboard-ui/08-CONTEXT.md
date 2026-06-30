# Phase 8: Multi-Instance Connections & Dashboard UI - Context

**Gathered:** 2026-06-29
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase delivers the database consolidation for multi-instance messaging (merging credentials and devices into `connections`), dynamic message routing on dispatch (loading credentials by ConnectionID or loading whatsmeow sessions by JID), and a Notion-style monochromatic admin dashboard built on Tailwind/daisyUI with dynamic onboarding logic.

</domain>

<decisions>
## Implementation Decisions

### Database & Migration
- **D-01 (Database Migration):** Write a Goose SQL migration that consolidates the `devices` and `channel_credentials` tables into a single unfiled `connections` table. The migration must automatically copy existing records (re-encrypting or copying encrypted byte payloads) to the new structure, then drop the legacy tables to ensure a transparent update in development/production environments.

### Outbound Routing & Concurrency
- **D-02 (Re-pairing Flow):** Disconnected or logged-out WhatsApp Web (whatsmeow) sessions change status to "Disconnected" (visual indicator: red badge) and expose a "Re-pair" button. Clicar neste botão abre o modal para gerar um novo QR Code associado ao mesmo ID de conexão, preservando logs e histórico de auditoria associados àquela linha.
- **D-03 (Workspace Connection Limits):** Restrict active whatsmeow WebSocket connections per workspace via an environment variable `PERGO_MAX_WHATSAPP_CONNECTIONS` (defaulting to 5) to protect the server from memory exhaustion (OOM). If exceeded during a pairing attempt, return HTTP 422. Stateless channels (Telegram, WABA) remain unrestricted.

### Interface & Navigation (Notion Vibe)
- **D-04 (UI Design Theme & Sidebar):** Apply an essentially monochromatic, clean, and light-themed layout (Notion style) using Tailwind CSS and daisyUI. Grayscale borders and backgrounds define structure. Accent and functional colors (green/blue) are reserved strictly for system states (active/connected/received events). The default navigation is a wide, collapsible left sidebar (Variant A).
- **D-05 (Dynamic Onboarding Logic):** Do not store dashboard UI state in the database. Compute the onboarding vs metrics layout dynamically on page request by checking active connections and API keys for the workspace. If `Count(APIKeys) == 0` or `Count(Connections) == 0`, show the 4-step progressive onboarding checklist; otherwise, render the operational metrics dashboard.

### Folded Todos
- **implement-dynamic-onboarding-logic.md:** Implement dynamic onboarding logic check in workspace repository / dashboard handlers.
- **implement-socks5-proxy-support-whatsmeow.md:** Add SOCKS5/HTTP proxy configuration support to the whatsmeow client connection setup to allow traffic isolation.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Spike & Sketch Findings (Blueprints)
- `.agents/skills/spike-findings-pergo/references/multi-instance-routing.md` — Multi-instance connection routing blueprint.
- `.agents/skills/sketch-findings-pergo/references/admin-dashboard-ui.md` — Admin dashboard layouts and daisyUI structures blueprint.

### Architecture & Strategy Notes
- `.planning/notes/cpaas-monetization-architecture.md` — SaaS strategy and proxy traffic isolation design.
- `.planning/notes/cpaas-dashboard-ux-architecture.md` — Notion styling guidelines and dynamic onboarding rules.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/session/manager.go`: ReconnectAll/reconnectDevice manages whatsmeow sessions (needs to load specific devices via JID).
- `internal/channel/whatsapp/client.go`: whatsmeow client initialization (needs custom SOCKS5/HTTP proxy dialer support).
- `internal/channel/registry.go`: Registry maps static dispatchers (`whatsapp`, `whatsapp_cloud`, `telegram`).

### Established Patterns
- **Static Dispatcher Registry:** Keep the registry keys static (`"telegram"`, `"whatsapp"`). Roteamento de credenciais e de sessões ativas é feito sob demanda na hora do `.Dispatch(ctx, payload)` usando `ConnectionID` (para stateless) ou `JID` (para stateful whatsmeow) em vez de instanciar múltiplos dispatchers dinamicamente.

### Integration Points
- `cmd/pergo/main.go`: Composition root wiring for adapters, queue workers, and session manager.
- `internal/platform/queue/worker.go`: Outbound queue worker (needs to pass `ConnectionID` / JID inside the message payload).

</code_context>

<specifics>
## Specific Ideas
- Monochromatic gray dashboard design with Lucide icons.
- Webhook simulation trigger in metrics dashboard for developer testing.

</specifics>

<deferred>
## Deferred Ideas
- stripe/billing managed cloud backend (recorded in seed `cpaas-billing-saas-infra.md`).

</deferred>

---

*Phase: 08-multi-instance-connections-dashboard-ui*
*Context gathered: 2026-06-29*
