# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/mzeahmed/noticoel/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/mzeahmed/noticoel/releases/tag/v0.1.0
