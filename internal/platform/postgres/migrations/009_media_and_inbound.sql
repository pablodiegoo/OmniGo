-- Migration 009_media_and_inbound.sql
-- +goose Up
ALTER TABLE workspaces ADD COLUMN pii_opt_in BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE inbound_dedups (
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    channel VARCHAR(50) NOT NULL,
    provider_message_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workspace_id, channel, provider_message_id)
);

CREATE INDEX idx_inbound_dedups_created_at ON inbound_dedups(created_at);

-- +goose Down
DROP TABLE IF EXISTS inbound_dedups;
ALTER TABLE workspaces DROP COLUMN IF EXISTS pii_opt_in;
