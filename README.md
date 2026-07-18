<div align="center">

<img src="assets/logo.png" alt="Noticeal Logo" width="220">

<br>

**A lightweight notification service for self-hosted infrastructures.**

Receive events. Notify anywhere.

<br>

[![Go Report Card](https://goreportcard.com/badge/github.com/mzeahmed/noticeal)](https://goreportcard.com/report/github.com/mzeahmed/noticeal)
[![Go Reference](https://pkg.go.dev/badge/github.com/mzeahmed/noticeal.svg)](https://pkg.go.dev/github.com/mzeahmed/noticeal)
[![Release](https://img.shields.io/github/v/release/mzeahmed/noticeal)](https://github.com/mzeahmed/noticeal/releases)
[![License](https://img.shields.io/github/license/mzeahmed/noticeal)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mzeahmed/noticeal)](app/go.mod)

</div>

---

# Overview

Noticeal is a lightweight notification service designed for self-hosted environments.

It receives events over HTTP and dispatches notifications to one or more channels.

The first version focuses on a simple use case:

> Receive Forgejo workflow events and send notifications.

The architecture is intentionally small, making it easy to deploy, understand and extend.

---

# Why Noticeal?

CI/CD pipelines continuously generate valuable events:

- Build succeeded
- Build failed
- Release created
- Deployment completed

Most self-hosted platforms provide webhooks, but turning those events into useful notifications often requires custom scripts.

Noticeal removes that complexity by providing a single notification endpoint.

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

Noticeal is distributed in two ways, both shipping the exact same application. Pick whichever fits your infrastructure.

## Option 1 — Standalone binary

Prebuilt binaries for Linux, macOS and Windows are published on the [Releases](https://github.com/mzeahmed/noticeal/releases) page.

```bash
tar -xzf noticeal_Linux_x86_64.tar.gz
./noticeal
```

## Option 2 — OCI image (Docker)

Every release is also published as an OCI image, built automatically by GoReleaser using [Ko](https://ko.build) — no Dockerfile is maintained in this project.

```bash
docker run \
  --rm \
  -p 8080:8080 \
  -v $(pwd)/config:/app/config \
  -v $(pwd)/data:/app/data \
  ghcr.io/mzeahmed/noticeal:latest
```

Or with Docker Compose:

```yaml
services:
  noticeal:
    image: ghcr.io/mzeahmed/noticeal:latest
    restart: unless-stopped

    volumes:
      - ./config:/app/config
      - ./data:/app/data

    ports:
      - "8080:8080"
```

> The OCI image contains the exact same Noticeal binary as the GitHub Release. Configuration and the SQLite database should be stored on mounted volumes. Docker is one deployment option among others — not a project dependency.

---

# Development

Local development uses [Air](https://github.com/air-verse/air) for hot reloading. It is a developer convenience only — never required in production, and not involved in the release build.

```bash
go install github.com/air-verse/air@latest
```

Run the application:

```bash
air
```

Air automatically rebuilds and restarts Noticeal whenever a `.go`, `.yaml` or `.sql` file changes.

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

The long-term vision is to gradually evolve Noticeal into a generic event routing platform.

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

⚠️ Noticeal is under active development.

The API is experimental and may change before the first stable release.

---

# License

MIT License.