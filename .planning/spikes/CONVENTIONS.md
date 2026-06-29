# Spike Conventions

Patterns and design choices established across multi-instance redesign spikes.

## Stack
- **Database:** PostgreSQL (with `pgcrypto` for transparent credentials encryption).
- **Language:** Go 1.25+ with Echo v5 (router) and pgx/v5 (database connectivity).

## Structure
- Spike artifacts are stored under `.planning/spikes/<NNN>-<name>/` containing a `README.md` and test verification files.
- Local compose-based containers run PostgreSQL on port `5433` and NATS on port `4222`.

## Patterns
- **Unified Connections Table:** Consolidate `devices` and `channel_credentials` into a single `connections` table, using a globally unique `sender_identity` column (representing the bot username, phone number, or JID) as the business routing key.
- **Dynamic Instance Routing:** Keep dispatchers statically registered (e.g. one `WABAAdapter` instance, one `TelegramAdapter` instance). The worker/API passes the connection ID/identities to the dispatcher payload, and the dispatcher resolves the specific credentials from the DB or active socket sessions in memory at dispatch time. This avoids dynamic pool creation overhead and memory leaks.
- **API `from` Field Routing:** `POST /api/v1/messages` resolves the target connection using the `from` parameter. If empty, it falls back to the connection where `is_default = true` for the requested channel.
