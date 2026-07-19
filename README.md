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
    Forgejo   Yoostart   BookingApp   Monitoring   Cron Jobs
        │        │           │            │            │
        └────────┴─────┬─────┴────────────┴────────────┘
                        ▼
                 HTTP REST API
                        │
                        ▼
                  Event Router
                        │
          ┌─────────────┼─────────────┐
          ▼              ▼             ▼
      Telegram        Discord        Email
```

Any application capable of sending an HTTP request can publish events to Noticoel.

---

# Features

Current features:

- HTTP event ingestion
- JSON event schema
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

The same binary and the same `config/config.yaml` work for every installation. What differs between installations is the environment they run in.

```text
config.yaml
    ↓
application configuration (server, database, which notifiers are enabled...)

environment variables
    ↓
secret values (tokens, chat IDs...)
```

`config/config.yaml` describes *how* Noticoel behaves and is meant to be committed/shared. Secrets never go in it — Noticoel reads them from environment variables at startup instead, so nothing sensitive needs to be baked into the binary, the image, or a config file.

Noticoel validates its required secrets at startup and refuses to start if one is missing. `NOTICOEL_AUTH_TOKEN` is always required; the Telegram variables are required only when `notifiers.telegram.enabled` is `true` in `config.yaml`.

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
  noticoel:
    image: ghcr.io/mzeahmed/noticoel:latest
    restart: unless-stopped

    environment:
      NOTICOEL_AUTH_TOKEN: ${NOTICOEL_AUTH_TOKEN}
      NOTICOEL_TELEGRAM_BOT_TOKEN: ${NOTICOEL_TELEGRAM_BOT_TOKEN}
      NOTICOEL_TELEGRAM_CHAT_ID: ${NOTICOEL_TELEGRAM_CHAT_ID}

    volumes:
      - ./config:/app/config
      - ./data:/app/data

    ports:
      - "8080:8080"
```

Docker Compose automatically injects those values into the container's environment (resolved from your shell or a `.env` file next to `docker-compose.yml`). Noticoel itself has no notion of Docker or `.env` files — it only ever reads real environment variables.

## Local secrets (.env)

`.env` is a **local development convenience only**. Noticoel never reads a `.env` file itself — only tooling (the Makefile, [Air](https://github.com/air-verse/air)) loads it and exports its contents into the process environment for you.

```bash
cp .env.example .env
```

Fill in your own values, keeping in mind that:

- `.env` is listed in `.gitignore` and must never be committed.
- Every developer keeps their own `.env` with their own values.
- In production there is no `.env` file — set real environment variables instead.

## Secrets

| Variable | Description |
|---|---|
| `NOTICOEL_AUTH_TOKEN` | Bearer token required to call the events API |
| `NOTICOEL_TELEGRAM_BOT_TOKEN` | Telegram Bot API token (see [Notifiers → Telegram](#telegram)) |
| `NOTICOEL_TELEGRAM_CHAT_ID` | Telegram chat or group ID to notify |

More variables will be added here as new notifiers are implemented (see [Roadmap](docs/roadmap.md)).

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

Set both as environment variables — see [Configuration](#configuration) for how, depending on your deployment. `notifiers.telegram.enabled` in `config/config.yaml` must also be `true` (the default).

---

# Examples

A CI/CD platform reporting a deployment:

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "source":"forgejo",
    "type":"workflow",
    "status":"success",
    "title":"Deployment completed",
    "message":"BookingApp deployed successfully"
  }'
```

A business application reporting a domain event:

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "source":"yoostart",
    "type":"subscription.created",
    "status":"info",
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