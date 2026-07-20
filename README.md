<div align="center">

<img src="assets/logo.png" alt="Noticoel Logo" width="220">

<br>

**A lightweight event hub for self-hosted infrastructures.**

Applications publish events. Noticoel decides what happens next.

<br>

[![Go Reference](https://pkg.go.dev/badge/github.com/mzeahmed/noticoel.svg)](https://pkg.go.dev/github.com/mzeahmed/noticoel)
[![Release](https://img.shields.io/github/v/release/mzeahmed/noticoel)](https://github.com/mzeahmed/noticoel/releases)
[![License](https://img.shields.io/github/license/mzeahmed/noticoel)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mzeahmed/noticoel?filename=app%2Fgo.mod)](app/go.mod)

</div>

---

# Overview

Noticoel is a lightweight event hub designed for self-hosted environments.

Applications and infrastructure services publish events to Noticoel over a single HTTP endpoint. Noticoel normalizes them and routes notifications to one or more destinations — Telegram, Discord, email, and more.

> Forgejo workflow events are one example. Noticoel is not tied to any single event source.

The architecture is intentionally small, making it easy to deploy, understand and extend.

---

# Why Noticoel?

Self-hosted infrastructures generate valuable events all the time:

- Deployment completed
- Workflow failed
- New subscription
- Payment failed
- Security alert
- Monitoring alert

Without a central hub, every application ends up coupled to whichever notification service it was wired to first — one app talks to Telegram, another to Discord, a third sends its own emails. Switching channels means touching every application again.

Noticoel removes that coupling. Applications publish events and know nothing about Telegram, Discord or Email — Noticoel owns that decision, in one place.

---

# Architecture

```text
  Native Event producers          Third-party systems
  (Yoostart, BookingApp,          (Forgejo, GitHub, GitLab, Gitea,
   internal apps, custom APIs)     monitoring systems...)
             │                              │
             │                              ▼
             │                          an adapter
             │                     (native payload → Event)
             │                              │
             └──────────────┬───────────────┘
                             ▼
                      HTTP REST API
                             │
                             ▼
                       Event Router
                             │
             ┌───────────────┼───────────────┐
             ▼                ▼               ▼
         Telegram          Discord           Email
```

An application that already speaks Noticoel's Event model publishes straight to the API — no adapter needed. A third-party system with its own webhook format goes through a dedicated adapter first, which converts its native payload into an Event before it reaches the same pipeline.

---

# Features

Current features:

- HTTP event ingestion
- JSON event schema
- Adapters for third-party webhooks (Forgejo, GitHub, GitLab, Gitea)
- Multiple notification channels
- YAML configuration
- SQLite persistence
- Structured logging
- Single binary, no dependencies
- Self-hosted

Planned:

- Rule-based routing
- Event history
- Dashboard
- Additional connectors

---

# Installation

Noticoel is distributed in two ways, both shipping the exact same generic application — no secrets baked into either one. Pick whichever fits your infrastructure, then head to [Configuration](#configuration) to give your installation its own secrets before running it.

## Option 1 — Standalone binary

Prebuilt binaries for Linux, macOS and Windows are published on the [Releases](https://github.com/mzeahmed/noticoel/releases) page.

```bash
tar -xzf noticoel_Linux_x86_64.tar.gz
```

This extracts the `noticoel` binary. See [Configuration](#configuration) before running it.

## Option 2 — OCI image (Docker)

Every release is also published as an OCI image, built automatically by GoReleaser using [Ko](https://ko.build) — no Dockerfile is maintained in this project.

```bash
docker pull ghcr.io/mzeahmed/noticoel:latest
```

> The OCI image contains the exact same Noticoel binary as the GitHub Release. Docker is one deployment option among others — not a project dependency.

See [Configuration](#configuration) for how to run it with your own secrets and config.

---

# Configuration

Noticoel follows the [Twelve-Factor App](https://12factor.net/config) convention: infrastructure configuration — server, database, which notifiers are enabled, their credentials — comes entirely from environment variables. A Docker deployment needs nothing but a docker-compose.yml with an `environment:` block; there is no config file to template or mount.

```text
environment variables
    ↓
server, database, notifiers, secrets — everything infrastructure-related

config/config.yaml (optional)
    ↓
future business configuration (routing rules, templates...) — unused today
```

Noticoel validates its required secrets at startup and refuses to start if one is missing. `NOTICOEL_AUTH_TOKEN` is always required; the Telegram variables are required only while `NOTICOEL_TELEGRAM_ENABLED` is `true` (the default).

## Environment variables

### Required

| Variable | Description |
|---|---|
| `NOTICOEL_AUTH_TOKEN` | Bearer token required to call the API |
| `NOTICOEL_TELEGRAM_BOT_TOKEN` | Telegram Bot API token — required while Telegram is enabled (see [Notifiers → Telegram](#telegram)) |
| `NOTICOEL_TELEGRAM_CHAT_ID` | Telegram chat or group ID to notify — required while Telegram is enabled |

### Optional

| Variable | Default | Description |
|---|---|---|
| `NOTICOEL_DEBUG` | `false` | Human-readable logs at DEBUG level instead of JSON at INFO |
| `NOTICOEL_SERVER_HOST` | `0.0.0.0` | HTTP server bind address |
| `NOTICOEL_SERVER_PORT` | `8080` | HTTP server port |
| `NOTICOEL_DATABASE_PATH` | `./data/noticoel.db` | SQLite database file path |
| `NOTICOEL_TELEGRAM_ENABLED` | `true` | Enable the Telegram notifier |
| `NOTICOEL_NTFY_ENABLED` | `false` | Enable the ntfy notifier |
| `NOTICOEL_NTFY_SERVER` | `https://ntfy.sh` | ntfy server URL |
| `NOTICOEL_NTFY_TOPIC` | *(empty)* | ntfy topic |
| `NOTICOEL_WEBHOOK_ENABLED` | `false` | Enable the generic webhook notifier |
| `NOTICOEL_WEBHOOK_URL` | *(empty)* | Webhook URL to call |
| `NOTICOEL_DISCORD_ENABLED` | `false` | Enable the Discord notifier |
| `NOTICOEL_DISCORD_WEBHOOK` | *(empty)* | Discord webhook URL |
| `NOTICOEL_EMAIL_ENABLED` | `false` | Enable the email (SMTP) notifier |
| `NOTICOEL_EMAIL_HOST` | *(empty)* | SMTP host |
| `NOTICOEL_EMAIL_PORT` | `587` | SMTP port |
| `NOTICOEL_EMAIL_USERNAME` | *(empty)* | SMTP username |
| `NOTICOEL_EMAIL_PASSWORD` | *(empty)* | SMTP password |
| `NOTICOEL_EMAIL_FROM` | *(empty)* | "From" address for outgoing emails |

> ntfy, generic webhook, Discord and email aren't implemented as notifiers yet (see [Roadmap](docs/roadmap.md)) — their configuration is already wired up so enabling them will be a config-only change once they land.

## Deployment scenarios

### Linux shell

```bash
export NOTICOEL_AUTH_TOKEN=xxxxxxxx
export NOTICOEL_TELEGRAM_BOT_TOKEN=xxxxxxxx
export NOTICOEL_TELEGRAM_CHAT_ID=-123456789

./noticoel
```

Export the variables however fits your setup — your shell, a systemd `EnvironmentFile`, or your process supervisor of choice — as long as they reach the process as real environment variables.

### Docker Compose

```yaml
services:

  # Fixes ownership of ./data before noticoel starts. The noticoel image
  # runs as a fixed non-root UID (65532) and has no Dockerfile (built via
  # ko) to chown at build time — without this, the bind mount inherits
  # the host's ownership and noticoel can't write its SQLite database.
  noticoel-init:
    image: busybox:latest
    command: chown -R 65532:65532 /data
    volumes:
      - ./data:/data

  noticoel:
    image: ghcr.io/mzeahmed/noticoel:latest
    restart: unless-stopped

    depends_on:
      noticoel-init:
        condition: service_completed_successfully

    environment:
      NOTICOEL_AUTH_TOKEN: ${NOTICOEL_AUTH_TOKEN}
      NOTICOEL_TELEGRAM_BOT_TOKEN: ${NOTICOEL_TELEGRAM_BOT_TOKEN}
      NOTICOEL_TELEGRAM_CHAT_ID: ${NOTICOEL_TELEGRAM_CHAT_ID}
      # Absolute path: the image has no WORKDIR, so a relative path would
      # resolve against "/" (the root, not writable by the image's
      # non-root user).
      NOTICOEL_DATABASE_PATH: /data/noticoel.db

    volumes:
      - ./data:/data

    ports:
      - "8080:8080"
```

`noticoel-init` only matters for a bind mount like `./data` above — Docker gives it the host directory's ownership, not the image's. Skip it if you use a named volume instead.

No config file to mount — the `environment:` block is the entire configuration. Docker Compose injects those values into the container's environment (resolved from your shell or a `.env` file next to `docker-compose.yml`). Noticoel itself has no notion of Docker or `.env` files — it only ever reads real environment variables.

Only add a `- ./config:/config` volume if you're also using the optional `config.yaml` below — same reasoning, an absolute container path.

## Local secrets (.env)

`.env` is a **local development convenience only**. Noticoel never reads a `.env` file itself — only tooling (the Makefile, [Air](https://github.com/air-verse/air)) loads it and exports its contents into the process environment for you.

```bash
cp .env.example .env
```

Fill in your own values, keeping in mind that:

- `.env` is listed in `.gitignore` and must never be committed.
- Every developer keeps their own `.env` with their own values.
- In production there is no `.env` file — set real environment variables instead.

## Optional config.yaml

`config/config.yaml` is reserved for future business configuration — routing rules, notification templates, and the like — that doesn't map cleanly onto environment variables. None of that exists yet, so most deployments don't need this file at all.

```bash
cp app/config/config.yaml.example app/config/config.yaml
```

If you're upgrading from an older Noticoel and still have a `config.yaml` with `server:`, `database:` or `notifiers:` sections, it keeps working: those values are used as a fallback for whichever environment variables you haven't set yet, so upgrading in place won't break your deployment. New deployments should just use environment variables.

---

# Development

Local development uses [Air](https://github.com/air-verse/air) for hot reloading. It is a developer convenience only — never required in production, and not involved in the release build.

```bash
go install github.com/air-verse/air@latest
```

Set up your local secrets first — see [Configuration → Local secrets (.env)](#local-secrets-env) — then run:

```bash
air
```

Air automatically rebuilds and restarts Noticoel whenever a `.go`, `.yaml` or `.sql` file changes.

---

# Notifiers

## Telegram

### 1. Create a bot

1. Open a chat with [@BotFather](https://t.me/BotFather) on Telegram.
2. Send `/newbot` and follow the prompts to pick a name and a username.
3. BotFather replies with a bot token, e.g. `123456789:AAExxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`. This is your `NOTICOEL_TELEGRAM_BOT_TOKEN`.

### 2. Get the chat ID

- **Direct message**: start a chat with your bot and send it any message.
- **Group**: add the bot to the group and send any message (mentioning the bot if the group has privacy mode restrictions).

Then open the following URL in a browser, replacing `<TOKEN>`:

```
https://api.telegram.org/bot<TOKEN>/getUpdates
```

Find `"chat":{"id": ...}` in the JSON response — that number is your `NOTICOEL_TELEGRAM_CHAT_ID` (group IDs are negative).

### 3. Configure Noticoel

Set both as environment variables — see [Configuration](#configuration) for how, depending on your deployment. `NOTICOEL_TELEGRAM_ENABLED` must also be `true` (the default).

---

# Adapters

A web application, SaaS platform or any other system that already speaks Noticoel's Event model doesn't need an adapter — it publishes straight to `POST /api/v1/events/create` (see [Examples](#examples)).

Third-party systems with their own webhook format go through a dedicated adapter instead, which converts their native payload into an Event:

| Adapter | Route | Native webhook |
|---|---|---|
| Forgejo | `POST /api/v1/adapters/forgejo` | release |
| GitHub | `POST /api/v1/adapters/github` | Actions `workflow_run` |
| GitLab | `POST /api/v1/adapters/gitlab` | pipeline |
| Gitea | `POST /api/v1/adapters/gitea` | push |

Point the corresponding webhook setting at `https://<your-noticoel-host>/api/v1/adapters/<name>`. These routes require the same bearer token as the Event API — whether you can supply it depends on what your provider's webhook UI allows (a custom header, a secret field...); check its docs.

See [Architecture → Adapters](docs/architecture.md#adapters) for how they're implemented.

---

# Examples

Forgejo publishing its native release webhook, through the Forgejo adapter:

```bash
curl -X POST http://localhost:8080/api/v1/adapters/forgejo \
  -H "Authorization: Bearer $NOTICOEL_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "published",
    "release": {
      "tag_name": "v1.4.0",
      "html_url": "https://git.example.com/example/example-app/releases/tag/v1.4.0",
      "author": { "login": "Ahmed" }
    },
    "repository": { "full_name": "example/example-app" }
  }'
```

A business application, as a native Event producer, reporting a domain event directly:

```bash
curl -X POST http://localhost:8080/api/v1/events/create \
  -H "Authorization: Bearer $NOTICOEL_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "source":"yoostart",
    "category":"billing",
    "type":"subscription.created",
    "severity":"info",
    "title":"New Premium subscription",
    "message":"A new Premium subscription has been purchased."
  }'
```

Neither application knows Noticoel notifies over Telegram — that decision lives entirely in Noticoel's configuration.

---

# Roadmap

The long-term vision is to gradually evolve Noticoel into a generic event routing platform.

The first milestone focuses on solving one problem well:

- Receive events
- Dispatch notifications
- Stay simple

Future versions will introduce:

- Rule-based routing
- Additional connectors
- Dashboard
- Event history
- More notification channels

---

# Documentation

- [Architecture](docs/architecture.md)
- [Roadmap](docs/roadmap.md)
- [Troubleshooting](docs/troubleshooting.md)

---

# Contributing

Contributions are welcome.

Feel free to open an issue or submit a pull request.

---

# Project Status

⚠️ Noticoel is under active development.

The API is experimental and may change before the first stable release.

---

# License

MIT License.