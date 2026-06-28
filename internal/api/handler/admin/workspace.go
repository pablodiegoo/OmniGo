package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"

	mw "github.com/pablojhp.omnigo/internal/api/middleware"
	"github.com/pablojhp.omnigo/internal/repository"
	"github.com/pablojhp.omnigo/templates/pages"
)

// WorkspaceHandler holds dependencies for workspace admin operations.
type WorkspaceHandler struct {
	Repo        *repository.WorkspaceRepository
	APIKeys     *repository.APIKeyRepository
	Credentials *repository.CredentialsRepository
	Templates   *repository.WABATemplateRepository
}

// List renders the workspace list page or HTMX fragment.
func (h *WorkspaceHandler) List(c *echo.Context) error {
	workspaces, err := h.Repo.List(c.Request().Context(), 50)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to load workspaces")
	}

	if mw.IsHTMX(c) {
		return mw.Render(c, http.StatusOK, pages.WorkspaceListContent(workspaces))
	}
	return mw.Render(c, http.StatusOK, pages.WorkspaceListPage(workspaces))
}

// Create handles workspace creation via POST form.
func (h *WorkspaceHandler) Create(c *echo.Context) error {
	name := c.FormValue("name")
	if name == "" {
		return c.String(http.StatusBadRequest, "name is required")
	}

	ws, err := h.Repo.Create(c.Request().Context(), name)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create workspace")
	}

	return mw.Render(c, http.StatusOK, pages.WorkspaceRow(*ws))
}

// Detail renders the workspace detail page with API keys and channel configuration.
func (h *WorkspaceHandler) Detail(c *echo.Context) error {
	idStr, err := echo.PathParam[string](c, "id")
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid workspace ID")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid workspace ID")
	}

	ws, err := h.Repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusNotFound, "workspace not found")
	}

	var keys []repository.APIKey
	if h.APIKeys != nil {
		keys, err = h.APIKeys.ListByWorkspace(c.Request().Context(), id)
		if err != nil {
			keys = nil // degrade gracefully
		}
	}

	var waba pages.WABAConfig
	var tg pages.TelegramConfig

	if h.Credentials != nil {
		wabaBytes, err := h.Credentials.Get(c.Request().Context(), id, "whatsapp_cloud")
		if err == nil {
			_ = json.Unmarshal(wabaBytes, &waba)
		}
		tgBytes, err := h.Credentials.Get(c.Request().Context(), id, "telegram")
		if err == nil {
			_ = json.Unmarshal(tgBytes, &tg)
		}
	}

	return mw.Render(c, http.StatusOK, pages.WorkspaceDetailPage(*ws, keys, waba, tg))
}

// ConfirmDelete returns an HTMX modal fragment for delete confirmation.
func (h *WorkspaceHandler) ConfirmDelete(c *echo.Context) error {
	idStr, err := echo.PathParam[string](c, "id")
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid workspace ID")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid workspace ID")
	}

	ws, err := h.Repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusNotFound, "workspace not found")
	}

	return mw.Render(c, http.StatusOK, pages.WorkspaceDeleteConfirm(*ws))
}

// Delete removes a workspace and returns empty 200 for HTMX to remove the row.
func (h *WorkspaceHandler) Delete(c *echo.Context) error {
	idStr, err := echo.PathParam[string](c, "id")
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid workspace ID")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid workspace ID")
	}

	if err := h.Repo.Delete(c.Request().Context(), id); err != nil {
		return c.String(http.StatusInternalServerError, "failed to delete workspace")
	}

	return c.NoContent(http.StatusOK)
}

// SaveCredentials handles form submission via POST for channel credentials.
func (h *WorkspaceHandler) SaveCredentials(c *echo.Context) error {
	idStr, err := echo.PathParam[string](c, "id")
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid workspace ID")
	}
	workspaceID, err := uuid.Parse(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid workspace ID")
	}

	channel, err := echo.PathParam[string](c, "channel")
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid channel param")
	}
	if channel != "whatsapp_cloud" && channel != "telegram" {
		return c.String(http.StatusBadRequest, "invalid channel")
	}

	var payload []byte
	if channel == "whatsapp_cloud" {
		payload, err = json.Marshal(pages.WABAConfig{
			PhoneNumberID: c.FormValue("phone_number_id"),
			Token:         c.FormValue("token"),
			WABAAccountID: c.FormValue("waba_account_id"),
		})
	} else {
		payload, err = json.Marshal(pages.TelegramConfig{
			Token: c.FormValue("token"),
		})
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to marshal credentials")
	}

	if err := h.Credentials.Save(c.Request().Context(), workspaceID, channel, payload); err != nil {
		return c.String(http.StatusInternalServerError, "failed to save credentials")
	}

	if channel == "whatsapp_cloud" {
		var waba pages.WABAConfig
		_ = json.Unmarshal(payload, &waba)
		// Run sync in background so HTTP response is returned immediately
		go h.syncTemplatesFromMeta(context.Background(), workspaceID, waba)
		return mw.Render(c, http.StatusOK, pages.WABACredentialsCard(idStr, waba))
	} else {
		var tg pages.TelegramConfig
		_ = json.Unmarshal(payload, &tg)
		return mw.Render(c, http.StatusOK, pages.TelegramCredentialsCard(idStr, tg))
	}
}

func (h *WorkspaceHandler) syncTemplatesFromMeta(ctx context.Context, workspaceID uuid.UUID, config pages.WABAConfig) {
	baseURL := "https://graph.facebook.com/v18.0"
	metaURL := fmt.Sprintf("%s/%s/message_templates?limit=100", baseURL, config.WABAAccountID)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, metaURL, nil)
	if err != nil {
		slog.Error("failed to create meta sync request", "error", err, "workspace_id", workspaceID)
		return
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failed to fetch meta templates", "error", err, "workspace_id", workspaceID)
		return
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read meta templates response", "error", err, "workspace_id", workspaceID)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("meta templates api error status", "status", resp.StatusCode, "body", string(respBytes), "workspace_id", workspaceID)
		return
	}

	type metaTemplate struct {
		ID         string            `json:"id"`
		Name       string            `json:"name"`
		Language   string            `json:"language"`
		Status     string            `json:"status"`
		Category   string            `json:"category"`
		Components []json.RawMessage `json:"components"`
	}

	type metaTemplatesResponse struct {
		Data []metaTemplate `json:"data"`
	}

	var metaResp metaTemplatesResponse
	if err := json.Unmarshal(respBytes, &metaResp); err != nil {
		slog.Error("failed to unmarshal meta templates", "error", err, "workspace_id", workspaceID)
		return
	}

	slog.Info("syncing templates from Meta", "count", len(metaResp.Data), "workspace_id", workspaceID)

	for _, t := range metaResp.Data {
		componentsJSON, err := json.Marshal(t.Components)
		if err != nil {
			slog.Error("failed to marshal components", "error", err, "template", t.Name)
			continue
		}

		dbTmpl := &repository.WABATemplate{
			WorkspaceID:    workspaceID,
			MetaTemplateID: t.ID,
			Name:           t.Name,
			Language:       t.Language,
			Status:         t.Status,
			Category:       t.Category,
			Components:     componentsJSON,
		}

		if h.Templates != nil {
			_, err = h.Templates.Upsert(ctx, dbTmpl)
			if err != nil {
				slog.Error("failed to upsert template in local DB", "error", err, "template", t.Name)
			}
		}
	}
}

// DeleteCredentials revokes channel credentials.
func (h *WorkspaceHandler) DeleteCredentials(c *echo.Context) error {
	idStr, err := echo.PathParam[string](c, "id")
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid workspace ID")
	}
	workspaceID, err := uuid.Parse(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid workspace ID")
	}

	channel, err := echo.PathParam[string](c, "channel")
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid channel param")
	}
	if channel != "whatsapp_cloud" && channel != "telegram" {
		return c.String(http.StatusBadRequest, "invalid channel")
	}

	if err := h.Credentials.Delete(c.Request().Context(), workspaceID, channel); err != nil {
		return c.String(http.StatusInternalServerError, "failed to delete credentials")
	}

	if channel == "whatsapp_cloud" {
		return mw.Render(c, http.StatusOK, pages.WABACredentialsCard(idStr, pages.WABAConfig{}))
	} else {
		return mw.Render(c, http.StatusOK, pages.TelegramCredentialsCard(idStr, pages.TelegramConfig{}))
	}
}

