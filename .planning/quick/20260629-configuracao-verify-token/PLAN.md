---
status: complete
date: 2026-06-29
description: Permitir a configuração e exibição do Verify Token do Webhook do WhatsApp Cloud (WABA) no painel de controle
---

# Plan - Configuração de Verify Token do WABA

Enable operators to view and customize the webhook verification token for WhatsApp Cloud (WABA) inside the workspace configuration dashboard.

## Tasks
1. **Update structs**:
   - Add `VerifyToken string json:"verify_token"` to `WABAConfig` in `templates/pages/workspaces.templ` and `internal/channel/whatsapp/waba.go`.
2. **Update Admin Form & Details View**:
   - Add Webhook Verify Token input field to the form in `templates/pages/workspaces.templ`.
   - Render the active verify token in the details list when configured (showing default fallback `pergo_verify_token_[workspace_id]` if empty).
3. **Bind form field**:
   - In `internal/api/handler/admin/workspace.go`, read and assign the `verify_token` field from the HTTP form request.
4. **Verification**:
   - Run `make generate` to compile the templates.
   - Run `make test` to verify handler and channel tests pass successfully.
