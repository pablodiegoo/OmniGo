---
status: complete
date: 2026-06-28
description: Completed spike for interactive template selection and media upload in Developer Playground
---

# Spike: Playground Template Interactive - Summary

## Work Done
1. **Dynamic Template Dropdown**:
   - Added `GET /admin/playground/templates` to populate a template dropdown dynamically when the workspace is selected or changed.
2. **Template Details & Variable Form**:
   - Added `GET /admin/playground/templates/details` which fetches a WABA template and displays details.
   - Used a defensive metadata loop that dynamically shows Category, Status, Language, Meta ID, and Created Date only if they exist.
   - Parsed components JSON on the fly to count the number of variables (`{{1}}`, `{{2}}`, etc.) and render an input field for each.
   - Rendered media header inputs if the template demands an `IMAGE`, `VIDEO`, or `DOCUMENT` header.
3. **Local File Upload (MinIO/S3)**:
   - Added `POST /admin/playground/upload` endpoint.
   - Files are stored in S3/MinIO using the workspace prefix: `{workspace_id}/{uuid}-{filename}`.
   - Returns the local proxy path: `/media/{workspace_id}/{uuid}-{filename}`.
4. **Real-time Live Preview**:
   - Fixed the live preview by storing raw template variables as dataset attributes (`data-body-text` and `data-header-format`) on the DOM container, bypassing Go-to-JS string escaping issues.
   - The preview updates instantly when the user types in any variable fields or specifies a media URL.
5. **Cleaned Up Form & Redundant Inputs**:
   - Removed the old redundant "WABA Template Options (Optional)" form fields.
   - Tied the visibility and requirement of the "Message Body" field dynamically: it is shown and `required` when no template is selected or when a non-WABA channel is chosen; it is hidden and non-required when a template is active.
6. **WABA Media Parameters Support**:
   - Extended the `wabaParameter` struct and serialization mapping in `waba.go` to support nested `image`, `video`, and `document` structures expected by Meta's Graph API.
