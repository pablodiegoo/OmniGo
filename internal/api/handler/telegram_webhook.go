package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/pablojhp.omnigo/internal/repository"
)

// TelegramWebhookHandler handles inbound webhooks from Telegram.
type TelegramWebhookHandler struct {
	credsRepo            *repository.CredentialsRepository
	recipientSessionRepo *repository.RecipientSessionRepository
}

// NewTelegramWebhookHandler creates a new TelegramWebhookHandler.
func NewTelegramWebhookHandler(credsRepo *repository.CredentialsRepository, recipientSessionRepo *repository.RecipientSessionRepository) *TelegramWebhookHandler {
	return &TelegramWebhookHandler{
		credsRepo:            credsRepo,
		recipientSessionRepo: recipientSessionRepo,
	}
}

type telegramUpdate struct {
	Message *telegramMessage `json:"message"`
}

type telegramMessage struct {
	Chat *telegramChat `json:"chat"`
}

type telegramChat struct {
	ID int64 `json:"id"`
}

type telegramConfig struct {
	Token       string `json:"token"`
	SecretToken string `json:"secret_token"`
}

// Handle processes the incoming Telegram webhook POST request.
func (h *TelegramWebhookHandler) Handle(c *echo.Context) error {
	workspaceIDStr, err := echo.PathParam[string](c, "workspace_id")
	if err != nil || workspaceIDStr == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// Retrieve secret token from headers
	receivedToken := c.Request().Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if receivedToken == "" {
		return c.NoContent(http.StatusForbidden)
	}

	// Load registered credentials for the workspace
	credsBytes, err := h.credsRepo.Get(c.Request().Context(), workspaceID, "telegram")
	if err != nil {
		return c.NoContent(http.StatusForbidden)
	}

	var config telegramConfig
	if err := json.Unmarshal(credsBytes, &config); err != nil {
		return c.NoContent(http.StatusForbidden)
	}

	// Validate secret token
	if config.SecretToken == "" || receivedToken != config.SecretToken {
		return c.NoContent(http.StatusForbidden)
	}

	// Bind request body
	var update telegramUpdate
	if err := c.Bind(&update); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if update.Message != nil && update.Message.Chat != nil {
		chatIDStr := strconv.FormatInt(update.Message.Chat.ID, 10)
		err := h.recipientSessionRepo.Upsert(c.Request().Context(), workspaceID, chatIDStr, "telegram", time.Now().UTC())
		if err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
