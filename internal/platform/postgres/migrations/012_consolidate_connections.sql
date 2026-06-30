-- +goose Up
-- +goose StatementBegin
CREATE TABLE connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    channel TEXT NOT NULL, -- 'whatsapp' (web), 'whatsapp_cloud' (WABA), 'telegram'
    sender_identity TEXT NOT NULL, -- Unique identifier (phone number or bot username, e.g. '+5511999990001', '@bot_username')
    status TEXT NOT NULL DEFAULT 'pending',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Encrypted Credentials (JSON object representing WABA config or Telegram Token)
    credentials BYTEA,
    key_id TEXT,
    key_version INT NOT NULL DEFAULT 1,
    
    -- WhatsApp Web (whatsmeow) specific fields
    jid TEXT,
    connected_since TIMESTAMPTZ,
    
    -- Traffic isolation proxy configuration
    proxy_url TEXT, -- SOCKS5/HTTP proxy string (e.g. 'socks5://user:pass@host:port')
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    UNIQUE (sender_identity)
);

CREATE INDEX idx_connections_workspace_id ON connections(workspace_id);
CREATE INDEX idx_connections_channel ON connections(channel);

-- Migrate existing devices (whatsmeow sessions)
INSERT INTO connections (
    id, workspace_id, name, channel, sender_identity, status, is_default, jid, connected_since, created_at, updated_at
)
SELECT 
    id,
    workspace_id,
    'WhatsApp Web - ' || COALESCE(phone, id::text),
    'whatsapp',
    COALESCE(phone, jid, id::text),
    status,
    FALSE,
    jid,
    connected_since,
    created_at,
    updated_at
FROM devices;

-- Migrate existing channel credentials (WABA and Telegram)
INSERT INTO connections (
    id, workspace_id, name, channel, sender_identity, status, is_default, credentials, key_id, key_version, created_at, updated_at
)
SELECT 
    id,
    workspace_id,
    CASE WHEN channel = 'telegram' THEN 'Telegram Bot' ELSE 'WhatsApp WABA' END,
    channel,
    'legacy_' || channel || '_' || id::text, -- Safe unique placeholder
    'active',
    TRUE,
    credentials,
    key_id,
    key_version,
    created_at,
    updated_at
FROM channel_credentials;

-- Ensure at least one default WhatsApp connection exists per workspace
WITH ranked_whatsapp AS (
    SELECT id, ROW_NUMBER() OVER (PARTITION BY workspace_id ORDER BY created_at) as rn
    FROM connections
    WHERE channel = 'whatsapp'
)
UPDATE connections
SET is_default = TRUE
WHERE id IN (SELECT id FROM ranked_whatsapp WHERE rn = 1);

-- Drop legacy tables
DROP TABLE IF EXISTS devices CASCADE;
DROP TABLE IF EXISTS channel_credentials CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE channel_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    channel TEXT NOT NULL,
    credentials BYTEA NOT NULL,
    key_id TEXT NOT NULL,
    key_version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (workspace_id, channel)
);

CREATE INDEX idx_channel_credentials_workspace ON channel_credentials(workspace_id);

CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id),
    channel TEXT NOT NULL,
    device_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    credentials BYTEA,
    key_id TEXT,
    key_version INT DEFAULT 1,
    jid TEXT,
    phone TEXT,
    connected_since TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_devices_workspace_id ON devices(workspace_id);
CREATE UNIQUE INDEX idx_devices_jid ON devices(jid) WHERE jid IS NOT NULL;

-- Populate devices from connections (only whatsapp)
INSERT INTO devices (
    id, workspace_id, channel, device_id, status, credentials, key_id, key_version, jid, phone, connected_since, created_at, updated_at
)
SELECT 
    id,
    workspace_id,
    channel,
    id::text,
    status,
    credentials,
    key_id,
    key_version,
    jid,
    sender_identity,
    connected_since,
    created_at,
    updated_at
FROM connections
WHERE channel = 'whatsapp';

-- Populate channel_credentials from connections (only telegram and whatsapp_cloud)
INSERT INTO channel_credentials (
    id, workspace_id, channel, credentials, key_id, key_version, created_at, updated_at
)
SELECT DISTINCT ON (workspace_id, channel)
    id,
    workspace_id,
    channel,
    credentials,
    key_id,
    key_version,
    created_at,
    updated_at
FROM connections
WHERE channel IN ('telegram', 'whatsapp_cloud');

DROP TABLE IF EXISTS connections CASCADE;
-- +goose StatementEnd
