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
  Native Event producers        Third-party systems (proprietary payload)
  (Yoostart, BookingApp,        (Forgejo, GitHub, GitLab, Gitea,
   internal apps, custom APIs)   monitoring systems...)
             │                              │
             │                              ▼
             │                       an adapter
             │                     (payload → Event)
             │                              │
             └──────────────┬───────────────┘
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

A **native Event producer** already speaks Noticoel's Event model, so it publishes straight to `POST /api/v1/events/create` — no adapter involved. A **third-party system** has its own proprietary webhook format, so it goes through a dedicated adapter first, which converts that payload into an Event before handing it to the same pipeline — see [Adapters](#adapters).

---

# Request Flow

1. An application sends an HTTP request: either an Event directly, or a native payload to one of its adapter's routes.
2. If it came through an adapter, the adapter converts the native payload into an Event first.
3. Noticoel validates the Event.
4. The Event is persisted.
5. The Dispatcher forwards the Event to every enabled notifier.
6. Each notifier sends its own notification.

---

# Core Components

## API

The API exposes a minimal HTTP interface.

Endpoints:

```
POST /api/v1/events/create
GET  /api/v1/events/list
```

Receive or list events, in Noticoel's internal Event shape.

```
POST /api/v1/adapters/{name}
```

Receive a native payload from a specific third-party adapter (e.g. `forgejo`, `github`), converted to an Event before entering the same pipeline as the endpoint above.

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

The Event is the internal representation of a notification request — the common shape every producer publishes to Noticoel, whether directly (a native Event producer) or through an adapter (a third-party system).

```go
type Event struct {
    Source   string            `json:"source"`
    Category string            `json:"category,omitempty"`
    Type     string            `json:"type"`
    Severity Severity          `json:"severity"`
    Title    string            `json:"title"`
    Message  string            `json:"message"`
    Metadata map[string]string `json:"metadata,omitempty"`
}
```

- **Source** — the application or service that published the event.
- **Category** — an optional grouping for related event types (`billing`, `ci`, `auth`...).
- **Type** — the specific event within Source/Category, e.g. `subscription.created`.
- **Severity** — a closed set of values (`info`, `warning`, `error`, `critical`) driving how notifiers present the event and, later, routing rules.
- **Metadata** — arbitrary producer-specific context that doesn't belong in Title/Message.

Every notifier receives the same Event object.

---

## Adapters

Noticoel has two kinds of clients.

**Native Event producers** already speak Noticoel's Event model — web applications, SaaS platforms, internal business apps, custom services, anything you (or someone) built to call Noticoel directly. They publish straight to `POST /api/v1/events/create`. No adapter, because there is nothing to adapt: the payload already is an Event.

**Third-party systems** are external systems with a proprietary webhook format you don't control — Forgejo, GitHub, GitLab, Gitea, a monitoring tool. Each gets a dedicated adapter that converts its native payload into an Event.

A web application is a native Event producer, not an adapter — it never needs, and should never get, an adapter package. `internal/adapters/` exists only for third-party systems whose payload Noticoel doesn't control.

An adapter owns:

- its native payload shape (its own `Payload` type, decoded straight from the request body)
- the conversion from that payload into an Event
- its own route, `POST /api/v1/adapters/{name}`

and hands the resulting Event to the same `event.Service` the generic API uses, so persistence, validation and dispatch happen in exactly one place regardless of whether the Event arrived as-is or was adapted.

Adapters are independent, concrete packages — there is no shared `Adapter` interface, because nothing in the codebase needs to treat them polymorphically (unlike notifiers, which the Dispatcher fans a single Event out to, adapters each own a distinct route and are never iterated over as a group). Adding one means adding a new package under `internal/adapters/`, not modifying an existing one, an existing adapter, or the Event API.

The Dispatcher and every notifier only ever see an Event — they have no notion of Forgejo, GitHub, GitLab, Gitea, or any other producer, native or third-party.

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

        config/
        database/
        adapters/
            forgejo/
            gitea/
            github/
            gitlab/
        dispatcher/
        logger/
        modules/
            auth/
            event/
            health/
        notifier/
            telegram/
        router/
        version/

    config/
        config.yaml.example

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

Builds the application's top-level `http.Handler` on a standard library `http.ServeMux`, by registering each module's and each adapter's routes onto it. No third-party router.

---

## modules

Each module owns one feature end-to-end: its handler, its routes (via a `RegisterRoutes` method), and, where relevant, its own model and business logic.

- **health** — `GET /health` and `GET /version`.
- **event** — `POST /api/v1/events/create` and `GET /api/v1/events/list`; owns the `Event` model, its validation, and persists received events to SQLite via its `Service`. Its `Service` is also reused directly by every adapter, so an event's persistence and dispatch path is identical no matter where it came from.
- **auth** — the bearer token middleware guarding authenticated routes.

---

## adapters

Each adapter converts one third-party system's native webhook payload into an `event.Event` — see [Adapters](#adapters) above. Follows the same `Module` / `RegisterRoutes` shape as `modules`, kept in its own top-level package because it is a different concern: modules serve Noticoel's own API, adapters translate someone else's. Native Event producers don't get a package here at all — they use `event`'s routes directly.

- **forgejo** — `POST /api/v1/adapters/forgejo`; converts a Forgejo release webhook.
- **github** — `POST /api/v1/adapters/github`; converts a GitHub Actions `workflow_run` webhook.
- **gitlab** — `POST /api/v1/adapters/gitlab`; converts a GitLab pipeline webhook.
- **gitea** — `POST /api/v1/adapters/gitea`; converts a Gitea push webhook.

A web application, SaaS platform, or any other native Event producer does **not** get a package here — see [Adapters](#adapters) above.

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

Schema lives in `internal/database/migrations` (Goose); queries live in `internal/database/queries` and are compiled into typed Go code under `internal/database/sqlc` by running `make sqlc` (sqlc is a codegen tool only — the generated code has no sqlc runtime dependency, just `database/sql`). Each module that needs persistence (e.g. `event`) consumes the generated `*sqlc.Queries` directly from its own `Service`, with no repository layer in between.

---

## config

Loads infrastructure configuration (server, database, notifiers, secrets) entirely from `NOTICOEL_*` environment variables — see [Configuration](#configuration). The YAML file it also accepts is optional: reserved for future business configuration, and read with gopkg.in/yaml.v3 only as a fallback for deployments still on a pre-2.0 `config.yaml`.

---

## logger

Centralized structured logging, built entirely on the Go standard library's `log/slog` — no third-party logging dependency.

In development (`NOTICOEL_DEBUG=true`), it writes human-readable text to stdout at DEBUG level. In production, it writes JSON to stdout at INFO level.

---

## version

Application version information.

---

# Configuration

Infrastructure configuration — server, database, which notifiers are enabled, their credentials — is read entirely from `NOTICOEL_*` environment variables, following the [Twelve-Factor App](https://12factor.net/config) convention. There is no config file to template or mount for a standard deployment.

```go
cfg := config.Load("config/config.yaml")
```

`Load` reads every field from its environment variable, e.g. `NOTICOEL_SERVER_PORT`, `NOTICOEL_DATABASE_PATH`, `NOTICOEL_TELEGRAM_ENABLED`, falling back to a hardcoded default (`8080`, `./data/noticoel.db`, `true`...) when unset. `NOTICOEL_AUTH_TOKEN` has no default and is always required; the Telegram credentials are required only while `NOTICOEL_TELEGRAM_ENABLED` is `true`.

The `path` argument is optional in every sense: if the file doesn't exist, `Load` just skips it. It exists for two things:

1. **A pre-2.0 `config.yaml`** (with `server:`, `database:`, `notifiers:` sections) — if present, it's decoded and used as a fallback for any field whose environment variable isn't set, so an in-place upgrade doesn't break an existing deployment.
2. **Future business configuration** — routing rules, notification templates — that doesn't map cleanly onto environment variables. Nothing uses this yet.

Set secrets via a `.env` file at the repository root for local development (copy `.env.example` to `.env`), or inject them directly as real environment variables in production. See the README's [Configuration](../README.md#configuration) section for the full variable list.

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