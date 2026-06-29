---
status: complete
date: 2026-06-29
description: Permitir a configuração e exibição do Verify Token do Webhook do WhatsApp Cloud (WABA) no painel de controle
---

# Quick Task: configuracao-verify-token - Summary

## Work Done
1. **Root Cause Resolved**:
   - The handler verification `HandleGet` checked `creds.VerifyToken` or a default fallback token format `pergo_verify_token_[workspace_id]`.
   - However, the `VerifyToken` field was missing from the `WABAConfig` struct, was not present in the HTML configuration form, and was not bound in the admin workspace handler.
   - Consequently, the `VerifyToken` parameter was never saved in the database, rendering any custom Meta verification tokens invalid in production.
2. **Implementation**:
   - Added `VerifyToken` with JSON tag `verify_token` to `WABAConfig` in both `templates/pages/workspaces.templ` and `internal/channel/whatsapp/waba.go`.
   - Updated the form in `templates/pages/workspaces.templ` to expose the "Webhook Verify Token" input field.
   - Enhanced the workspace details list in `templates/pages/workspaces.templ` to show the active or default verify token.
   - Bound the form field `verify_token` inside the workspace handler `internal/api/handler/admin/workspace.go`.
3. **Verification**:
   - Successfully compiled the regenerated templates via `make generate` and verified that all backend and API tests pass (`make test`).
   - Verified that both custom and default webhook verification tokens correctly respond with 200 OK and return the challenge query parameter.
