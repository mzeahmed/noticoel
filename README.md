<div align="center">

<img src="assets/logo.png" alt="Noticoel Logo" width="220">

<br>

**A lightweight notification service for self-hosted infrastructures.**

Receive events. Notify anywhere.

<br>

[![Go Report Card](https://goreportcard.com/badge/github.com/mzeahmed/noticoel)](https://goreportcard.com/report/github.com/mzeahmed/noticoel)
[![Go Reference](https://pkg.go.dev/badge/github.com/mzeahmed/noticoel.svg)](https://pkg.go.dev/github.com/mzeahmed/noticoel)
[![Release](https://img.shields.io/github/v/release/mzeahmed/noticoel)](https://github.com/mzeahmed/noticoel/releases)
[![License](https://img.shields.io/github/license/mzeahmed/noticoel)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mzeahmed/noticoel)](app/go.mod)

</div>

---

# Overview

Noticoel is a lightweight notification service designed for self-hosted environments.

It receives events over HTTP and dispatches notifications to one or more channels.

The first version focuses on a simple use case:

> Receive Forgejo workflow events and send notifications.

The architecture is intentionally small, making it easy to deploy, understand and extend.

---

# Why Noticoel?

CI/CD pipelines continuously generate valuable events:

- Build succeeded
- Build failed
- Release created
- Deployment completed

Most self-hosted platforms provide webhooks, but turning those events into useful notifications often requires custom scripts.

Noticoel removes that complexity by providing a single notification endpoint.

---

# Architecture

```text
        Forgejo
           │
           ▼
     HTTP REST API
           │
           ▼
      Dispatcher
           │
     ┌─────┴─────┐
     ▼           ▼
 Discord       ntfy
               ▼
             Email
```

---

# Features

Current features:

- REST API
- JSON events
- Multiple notification channels
- YAML configuration
- SQLite storage
- Structured logging
- Single binary, no dependencies
- Self-hosted

---

# Installation

Noticoel is distributed in two ways, both shipping the exact same application. Pick whichever fits your infrastructure.

## Option 1 — Standalone binary

Prebuilt binaries for Linux, macOS and Windows are published on the [Releases](https://github.com/mzeahmed/noticoel/releases) page.

```bash
tar -xzf noticoel_Linux_x86_64.tar.gz
./noticoel
```

## Option 2 — OCI image (Docker)

Every release is also published as an OCI image, built automatically by GoReleaser using [Ko](https://ko.build) — no Dockerfile is maintained in this project.

```bash
docker run \
  --rm \
  -p 8080:8080 \
  -v $(pwd)/config:/app/config \
  -v $(pwd)/data:/app/data \
  ghcr.io/mzeahmed/noticoel:latest
```

Or with Docker Compose:

```yaml
services:
  noticoel:
    image: ghcr.io/mzeahmed/noticoel:latest
    restart: unless-stopped

    volumes:
      - ./config:/app/config
      - ./data:/app/data

    ports:
      - "8080:8080"
```

> The OCI image contains the exact same Noticoel binary as the GitHub Release. Configuration and the SQLite database should be stored on mounted volumes. Docker is one deployment option among others — not a project dependency.

---

# Development

Local development uses [Air](https://github.com/air-verse/air) for hot reloading. It is a developer convenience only — never required in production, and not involved in the release build.

```bash
go install github.com/air-verse/air@latest
```

Copy the environment template and set your own token:

```bash
cp .env.example .env
```

Run the application:

```bash
air
```

Air automatically rebuilds and restarts Noticoel whenever a `.go`, `.yaml` or `.sql` file changes.

---

# Notifiers

## Telegram

> The Telegram notifier is not wired up yet (see [Roadmap](docs/roadmap.md)), but the credentials below are already read from the environment, so you can set them up ahead of time.

### 1. Create a bot

1. Open a chat with [@BotFather](https://t.me/BotFather) on Telegram.
2. Send `/newbot` and follow the prompts to pick a name and a username.
3. BotFather replies with a bot token, e.g. `123456789:AAExxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`. This is your `TELEGRAM_BOT_TOKEN`.

### 2. Get the chat ID

- **Direct message**: start a chat with your bot and send it any message.
- **Group**: add the bot to the group and send any message (mentioning the bot if the group has privacy mode restrictions).

Then open the following URL in a browser, replacing `<TOKEN>`:

```
https://api.telegram.org/bot<TOKEN>/getUpdates
```

Find `"chat":{"id": ...}` in the JSON response — that number is your `TELEGRAM_CHAT_ID` (group IDs are negative).

### 3. Configure Noticoel

Set both values in your `.env` (local dev, see [Development](#development)) or as real environment variables in production:

```bash
TELEGRAM_BOT_TOKEN=123456789:AAExxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TELEGRAM_CHAT_ID=-123456789
```

`notifiers.telegram.enabled` in `config/config.yaml` must also be `true` (the default).

---

# Example

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