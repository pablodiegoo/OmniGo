-- +goose Up
CREATE TABLE recipient_sessions (
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    recipient_phone TEXT NOT NULL,
    channel TEXT NOT NULL,
    last_inbound_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (workspace_id, recipient_phone, channel)
);

CREATE INDEX idx_recipient_sessions_workspace ON recipient_sessions(workspace_id);

-- +goose Down
DROP TABLE IF EXISTS recipient_sessions;
