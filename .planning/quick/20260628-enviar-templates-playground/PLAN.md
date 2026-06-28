---
status: complete
date: 2026-06-28
description: Implement WABA message template enqueuing support inside Developer Playground
---

# Plan - Enviar Templates no Playground

Support testing WABA message templates directly from the Developer Playground UI.

## Tasks
1. **Playground Form Inputs**: Add fields for `template_name`, `template_language`, and `template_components` to the Playground send form.
2. **Vanilla JS UI Toggles**: Toggle visibility of WABA-specific options when `whatsapp_cloud` is selected, and dynamically toggle `required` state of message body.
3. **Queue Message Processing**: Update `Send` handler in `playground.go` to parse WABA template parameters JSON and populate the enqueued `QueueMessage` struct.
4. **Verification**: Run `templ generate`, `go build`, and verify tests pass.
