# Phase 2: Admin Shell — Phase Context

**Gathered:** 2026-06-25
**Status:** Ready for planning
**Mode:** Auto-generated (smart discuss — autonomous mode, recommended defaults applied)

<domain>
## Phase Boundary

Operators can manage workspaces, API keys, and review audit logs through a server-rendered admin panel built on Echo + Templ + HTMX. The admin panel provides a visual interface for the identity and audit infrastructure built in Phase 1.

</domain>

<decisions>
## Implementation Decisions

### UI Architecture
- **Templ** for compile-time type-safe HTML templates (not hand-written HTML)
- **HTMX 2.x** for fragment-based interactions (no full-page reloads)
- **Echo v5** serves both API and admin routes on the same port
- Admin routes prefixed with `/admin/` to separate from public API
- HTMX fragment detection via `HX-Request` header check in Echo middleware

### Layout & Navigation
- **Sidebar navigation** with sections: Workspaces, API Keys, Audit Logs
- Dashboard landing page showing workspace count, recent audit activity, system health
- Responsive design using simple CSS (no framework dependency)

### Workspace Management
- Workspace list view with search/filter
- Create workspace form: name (required), description (optional)
- Workspace detail view showing API keys and recent audit entries
- Delete workspace with confirmation dialog (HTMX modal)

### API Key Management
- API key list per workspace with status indicator (active/revoked)
- Generate new key: button triggers key generation, displays once, then hashes
- Revoke key: confirmation dialog, immediate revocation
- Key display: show prefix + masked suffix for identification

### Audit Log Review
- Table view with columns: timestamp, workspace, trace_id, event_type, status
- Filter by: workspace dropdown, time range picker, event type dropdown
- Search by trace_id (exact match)
- Export as CSV (server-side generation, download link)
- Pagination with 50 rows per page

### HTMX Fragment Strategy
- List views return HTML fragments (no layout wrapper)
- Form submissions return fragments for target update area
- Delete confirmations use HTMX modal pattern
- Navigation clicks return full pages (layout + content)

### Styling
- Minimal CSS with CSS custom properties for theming
- No external CSS framework (keep dependency footprint small)
- Consistent spacing and typography using CSS variables

</decisions>

<code_context>
## Existing Code Insights

### Reusable Assets
- Phase 1 provides: workspace/repository/apikey repository, auth middleware, audit batch writer, health handlers
- Echo v5 already configured with middleware stack
- pgxpool and stdlib bridge already wired

### Established Patterns
- Repository pattern: `internal/repository/workspace.go`, `internal/repository/apikey.go`
- Handler pattern: `internal/api/handler/health.go`
- Middleware pattern: `internal/api/middleware/auth.go`, `internal/api/middleware/trace.go`

### Integration Points
- Admin routes mount on existing Echo instance from Phase 1
- Auth middleware extended to support session-based auth for admin (not just API keys)
- Audit log queries reuse existing `audit_logs` table and batch writer

</code_context>

<specifics>
## Specific Ideas

- Admin panel should feel like a lightweight operator console, not a full CRM
- Focus on functional clarity over visual polish
- HTMX interactions should feel snappy (fragment responses, not full page loads)

</specifics>

<deferred>
## Deferred Ideas

- Real-time WebSocket updates for session status (Phase 4 — WhatsApp Web integration)
- Dashboard charts/graphs (can be added later without schema changes)
- Multi-user admin authentication with roles (MVP uses single-operator model)
- API-only admin management (CLI/SDK) — admin panel is the primary interface for MVP

</deferred>
