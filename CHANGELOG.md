# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Adapters for third-party webhooks — Forgejo, GitHub, GitLab, Gitea — each converting its native payload into the internal Event model behind its own `POST /api/v1/adapters/{name}` route, then flowing through the same validation, persistence and dispatch as the generic Event API
- Infrastructure configuration (server, database, which notifiers are enabled, their credentials) can now be set entirely through `NOTICOEL_*` environment variables — a Docker deployment no longer needs to mount a config file, just an `environment:` block

### Changed

- `config/config.yaml` is now optional and reserved for future business configuration; an existing pre-2.0 file (with `server:`, `database:`, `notifiers:` sections) still works as a fallback for any environment variable left unset

## [0.1.3] - 2026-07-19

### Added

- `category` field on events, for grouping related event types (e.g. `billing`, `ci`)
- `Severity` type (`info`, `warning`, `error`, `critical`), validated on every event instead of an arbitrary string

### Changed

- **Breaking:** renamed the `status` event field to `severity`; it must now be one of `info`, `warning`, `error`, `critical`
- **Breaking:** renamed the `data` event field to `metadata`

## [0.1.0] - 2026-07-19

First release. Noticoel receives events over HTTP and dispatches notifications to Telegram.

### Added

- `POST /api/v1/events` to receive an event, `GET /api/v1/events` to list stored events with pagination (`limit`/`offset`)
- Bearer token authentication on the events API
- SQLite persistence via Goose migrations and sqlc-generated queries
- Telegram notifier: events are dispatched to a Telegram chat through a pluggable notifier/dispatcher
- `GET /health` and `GET /version` endpoints
- Structured logging (`log/slog`), request logging and panic recovery middleware
- Graceful shutdown
- YAML configuration, with secrets (`NOTICOEL_AUTH_TOKEN`, `NOTICOEL_TELEGRAM_BOT_TOKEN`, `NOTICOEL_TELEGRAM_CHAT_ID`) read from environment variables / `.env`
- Example scripts (`send.sh`, `list.sh`) and sample event payloads
- Single-binary distribution for Linux, macOS and Windows via GoReleaser, plus an OCI image
