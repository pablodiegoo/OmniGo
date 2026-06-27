-- +goose Up
CREATE TABLE webhook_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    secret BYTEA NOT NULL,
    key_id TEXT NOT NULL,
    key_version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (workspace_id)
);

CREATE TABLE webhook_dlqs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    trace_id TEXT NOT NULL,
    message_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    webhook_url TEXT NOT NULL,
    last_attempt_at TIMESTAMPTZ NOT NULL,
    failure_reason TEXT,
    attempts INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_webhook_configs_workspace_id ON webhook_configs(workspace_id);
CREATE INDEX idx_webhook_dlqs_workspace_id ON webhook_dlqs(workspace_id);
CREATE INDEX idx_webhook_dlqs_trace_id ON webhook_dlqs(trace_id);
CREATE INDEX idx_webhook_dlqs_message_id ON webhook_dlqs(message_id);

-- +goose Down
DROP TABLE IF EXISTS webhook_dlqs;
DROP TABLE IF EXISTS webhook_configs;
