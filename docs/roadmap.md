# Noticoel Roadmap

> This roadmap reflects the planned evolution of Noticoel.
>
> Noticoel is a lightweight, self-hosted Go application distributed as a single binary.
> It is an event hub: applications publish events, and Noticoel routes notifications to one or more channels. It will gradually evolve into a generic event routing platform.

---

# Current Status

- [x] Project vision
- [x] Architecture
- [x] Documentation
- [x] Development started

---

# Phase 1 — Foundation

Build a solid foundation.

## Repository

- [x] Project structure
- [x] Documentation
- [x] Logo
- [x] Go module
- [x] Makefile
- [x] GitHub Actions
- [x] GoReleaser

## Core

- [x] Configuration
- [x] Logger
- [x] HTTP server
- [x] SQLite
- [x] Goose
- [x] sqlc

---

# Phase 2 — Notification Engine

Deliver the first production-ready notification service.

## Core

- [x] Event model
- [x] Event validation
- [x] Notifier interface
- [x] Dispatcher

## Persistence

- [x] Store events
- [ ] Store deliveries

## Notifiers

- [x] Telegram
- [ ] Discord
- [ ] ntfy
- [ ] Email
- [ ] Webhook

## Integrations

- [ ] Forgejo webhook
- [x] Manual event endpoint

---

# Phase 3 — Production Ready

Improve reliability and operational readiness.

- [x] Authentication
- [ ] Retry strategy
- [ ] Delivery status
- [ ] Metrics
- [x] Health checks
- [ ] Graceful shutdown
- [ ] Structured logging improvements

---

# Phase 4 — Event Routing

Transform Noticoel into an event router.

- [ ] Routing rules
- [ ] Filters
- [ ] Multiple destinations
- [ ] Templates
- [ ] Event transformations

---

# Phase 5 — Dashboard

Provide a lightweight web interface.

- [ ] Notification history
- [ ] Delivery history
- [ ] Statistics
- [ ] Configuration UI

---

# Phase 6 — Event Platform

Expand Noticoel into a generic event platform.

## Connectors

- [ ] GitHub
- [ ] GitLab
- [ ] Gitea
- [ ] Jenkins
- [ ] Generic Webhooks

## Extensibility

- [ ] Plugin system
- [ ] Public Go SDK
- [ ] REST API improvements

---

# Long-Term Vision

Noticoel is a lightweight event hub, distributed as a single self-hosted binary. Applications publish events over HTTP; Noticoel decides what happens next.

Over time, it will evolve into a generic event routing platform capable of receiving events from multiple systems, applying routing rules, storing delivery history, and delivering notifications through a wide range of channels while remaining lightweight, dependency-free, and easy to deploy.