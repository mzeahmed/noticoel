# Noticoel Architecture

> **Noticoel is a lightweight notification service for self-hosted infrastructures.**

---

# Introduction

Noticoel receives events over HTTP and dispatches notifications to one or more channels.

The first version has a single goal:

> Receive Forgejo workflow events and notify users.

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

A containerized deployment (Docker) may be offered later as an alternative, but it is not part of the current architecture.

---

## Lightweight

Noticoel should be easy to install, configure and maintain.

It uses an embedded SQLite database for its own operational data, so no external database service is required — only the configured notification channels depend on the outside world.

---

## Extensible

Although the first release targets Forgejo notifications, the architecture should allow new notification channels and event sources to be added without major changes.

---

# High-Level Architecture

```
            Forgejo
                │
                ▼
         HTTP REST API
                │
                ▼
          Event Handler
                │
                ▼
           Dispatcher
                │
      ┌─────────┴─────────┐
      ▼         ▼         ▼
   Discord     ntfy     Email
```

---

# Request Flow

1. Forgejo sends an HTTP request.
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
  driver: sqlite
  path: ./data/noticoel.db

notifications:

  discord:
    enabled: true
    webhook: https://discord...

  ntfy:
    enabled: false

  email:
    enabled: false
```

Secrets are not part of the YAML file: the bearer auth token and the Telegram credentials are read from `NOTICOEL_AUTH_TOKEN`, `NOTICOEL_TELEGRAM_BOT_TOKEN` and `NOTICOEL_TELEGRAM_CHAT_ID`. Set them via a `.env` file at the repository root for local development (copy `.env.example` to `.env`), or inject them directly as real environment variables in production. See the README's [Configuration](../README.md#configuration) section for details.

---

# Deployment

Noticoel ships as a single, self-contained binary. No container runtime is required.

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

GoReleaser produces prebuilt binaries and archives for Linux, macOS and Windows on every release — see the [Releases](https://github.com/mzeahmed/noticoel/releases) page.

The first version is designed to run on the same server as Forgejo:

```
Forgejo
      │
localhost
      │
      ▼
 Noticoel
      │
      ▼
 Discord
```

No public endpoint is required.

A Docker image is not provided today but may be introduced later as an optional, additional way to deploy Noticoel.

---

# Future Evolution

The current architecture intentionally focuses on notification delivery.

As the project grows, Noticoel may evolve into a more generic event routing platform by introducing features such as:

- routing rules
- multiple event sources
- dashboards
- plugins
- optional containerized deployment (Docker)

These features will build on the existing architecture without changing its core philosophy.

---

# Development Principles

- Keep packages small.
- Prefer composition over complexity.
- Avoid unnecessary abstractions.
- One responsibility per package.
- Keep dependencies minimal.
- Build only what is needed.