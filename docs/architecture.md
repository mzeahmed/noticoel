# Noticoel Architecture

> **Noticoel is a lightweight event hub for self-hosted infrastructures.**

---

# Introduction

Noticoel receives events over HTTP from any application or infrastructure service and dispatches notifications to one or more channels.

Applications publish events; they know nothing about Telegram, Discord or Email — Noticoel owns that decision.

> Forgejo workflow events are one example among many — see [Request Flow](#request-flow).

The architecture intentionally remains small and focused to provide a reliable foundation before introducing more advanced features.

---

# Design Principles

## Simplicity

Noticoel solves one problem well.

Avoid unnecessary abstractions and keep every package focused on a single responsibility.

---

## Self-hosted First

Noticoel is distributed as a single, dependency-free binary for Linux, macOS and Windows.

Typical deployments include:

- VPS
- Home server
- Any machine that can run a native binary

No external service is required.

Every release also publishes a matching OCI image as an alternative to the standalone binary — see [Deployment](#deployment).

---

## Lightweight

Noticoel should be easy to install, configure and maintain.

It uses an embedded SQLite database for its own operational data, so no external database service is required — only the configured notification channels depend on the outside world.

---

## Extensible

Noticoel is not tied to any single event source. The architecture allows new event producers and new notification channels to be added without major changes.

---

# High-Level Architecture

```
    Forgejo   Yoostart   BookingApp   Monitoring   Cron Jobs
        │        │           │            │            │
        └────────┴─────┬─────┴────────────┴────────────┘
                        ▼
                 HTTP REST API
                        │
                        ▼
                  Event Handler
                        │
                        ▼
                   Dispatcher
                        │
        ┌───────────────┼───────────────┬───────────────┐
        ▼                ▼               ▼               ▼
    Telegram          Discord           ntfy            Email
```

Any application capable of sending an HTTP request can publish events to Noticoel — Forgejo is just one example.

---

# Request Flow

1. An application sends an HTTP request describing an event (Forgejo, Yoostart, a monitoring system, a cron job...).
2. Noticoel validates the request.
3. The request is converted into an Event.
4. The Dispatcher forwards the Event to every enabled notifier.
5. Each notifier sends its own notification.

---

# Core Components

## API

The API exposes a minimal HTTP interface.

Endpoints:

```
POST /api/v1/events
```

Receive an event.

```
GET /health
```

Health check.

```
GET /version
```

Application version.

---

## Event

The Event is the internal representation of a notification request.

```go
type Event struct {
    Source  string            `json:"source"`
    Type    string            `json:"type"`
    Status  string            `json:"status"`
    Title   string            `json:"title"`
    Message string            `json:"message"`
    Data    map[string]string `json:"data,omitempty"`
}
```

Every notifier receives the same Event object.

---

## Dispatcher

The Dispatcher coordinates notification delivery.

Responsibilities:

- receive an Event
- call every enabled notifier
- collect delivery errors
- report failures

The Dispatcher contains no notification-specific logic.

---

## Notifiers

A notifier is responsible for delivering an Event to an external service.

Examples:

- Telegram
- Discord
- ntfy
- Email
- Webhook

Every notifier implements the same interface.

```go
type Notifier interface {
    Notify(ctx context.Context, event Event) error
}
```

---

# Project Structure

```
noticoel/

app/
    cmd/
        main.go

    internal/

        app/
        config/
        database/
        logger/
        modules/
            auth/
            events/
            health/
        response/
        router/
        version/

    config/
        config.yaml

    data/

assets/

docs/

examples/
    events/
    scripts/

Makefile
```

---

# Package Responsibilities

## router

Builds the application's top-level `http.Handler` on a standard library `http.ServeMux`, by registering each module's routes onto it. No third-party router.

---

## modules

Each module owns one feature end-to-end: its handler, its routes (via a `RegisterRoutes` method), and, where relevant, its own model and business logic.

- **health** — `GET /health` and `GET /version`.
- **events** — `POST /api/v1/events`; owns the `Event` model, its validation, and persists received events to SQLite via its `Service`.
- **auth** — the bearer token middleware guarding authenticated routes.

---

## response

Shared helper for writing JSON HTTP responses, used by every module's handlers.

---

## dispatcher

Coordinates notification delivery.

---

## notifier

Contains every notification implementation.

Each notifier is independent.

---

## database

Opens the SQLite database and runs pending Goose migrations at startup.

Schema lives in `internal/database/migrations` (Goose); queries live in `internal/database/queries` and are compiled into typed Go code under `internal/database/sqlc` by running `make sqlc` (sqlc is a codegen tool only — the generated code has no sqlc runtime dependency, just `database/sql`). Each module that needs persistence (e.g. `events`) consumes the generated `*sqlc.Queries` directly from its own `Service`, with no repository layer in between.

---

## config

Loads the application configuration from a YAML file using the Go standard library together with gopkg.in/yaml.v3. Secrets are the one exception: the bearer auth token and the Telegram credentials are read from environment variables instead, so they never need to be committed to the YAML file (see [Configuration](#configuration)).

---

## logger

Centralized structured logging, built entirely on the Go standard library's `log/slog` — no third-party logging dependency.

In development (`debug: true`), it writes human-readable text to stdout at DEBUG level. In production, it writes JSON to stdout at INFO level.

---

## version

Application version information.

---

# Configuration

Example:

```yaml
server:
  port: 8080

database:
  path: ./data/noticoel.db

notifiers:

  telegram:
    enabled: true

  ntfy:
    enabled: false

  webhook:
    enabled: false

  discord:
    enabled: false

  email:
    enabled: false
```

Secrets are not part of the YAML file: the bearer auth token and the Telegram credentials are read from `NOTICOEL_AUTH_TOKEN`, `NOTICOEL_TELEGRAM_BOT_TOKEN` and `NOTICOEL_TELEGRAM_CHAT_ID`. Set them via a `.env` file at the repository root for local development (copy `.env.example` to `.env`), or inject them directly as real environment variables in production. See the README's [Configuration](../README.md#configuration) section for details.

---

# Deployment

Noticoel ships as a single, self-contained binary, and every release also publishes a matching OCI image (built with [Ko](https://ko.build), no Dockerfile maintained) for container-based deployments.

Run it from source:

```
go -C app run ./cmd
```

Or build and execute the binary directly:

```
go -C app build -o noticoel ./cmd
./app/noticoel
```

(`make run` and `make build` wrap these same commands.)

For local development, [Air](https://github.com/air-verse/air) (configured in `.air.toml`) rebuilds and restarts Noticoel automatically whenever a `.go`, `.yaml` or `.sql` file changes. It is a developer-only convenience and plays no part in the release build.

GoReleaser produces prebuilt binaries, archives and the OCI image for Linux, macOS and Windows on every release — see the [Releases](https://github.com/mzeahmed/noticoel/releases) page.

A common deployment runs Noticoel alongside the applications that publish events to it, for example on the same server as Forgejo:

```
Forgejo
      │
localhost
      │
      ▼
 Noticoel
      │
      ▼
 Telegram
```

No public endpoint is required in that scenario — but Noticoel works just as well as a shared, reachable endpoint for several applications at once.

---

# Future Evolution

Noticoel already accepts events from any HTTP-capable application and routes them to multiple channels. The current architecture intentionally focuses on getting that core loop right.

As the project grows, Noticoel will build on this foundation by introducing:

- routing rules and filters
- event and delivery history
- a lightweight dashboard
- a plugin system for new connectors

These features build on the existing architecture without changing its core philosophy.

---

# Development Principles

- Keep packages small.
- Prefer composition over complexity.
- Avoid unnecessary abstractions.
- One responsibility per package.
- Keep dependencies minimal.
- Build only what is needed.