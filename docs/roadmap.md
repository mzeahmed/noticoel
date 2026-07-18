# Noticeal Roadmap

> This roadmap reflects the planned evolution of Noticeal.
>
> Noticeal is a lightweight, self-hosted Go application distributed as a single binary. It starts as a notification service for self-hosted infrastructures and will gradually evolve into a generic event routing platform.

---

# Current Status

- [x] Project vision
- [x] Architecture
- [x] Documentation
- [x] Development started

---

# Phase 1 — Foundation

Build a solid foundation for the project.

## Repository

- [x] Project structure
- [x] Documentation
- [x] Logo
- [x] Go module
- [x] Makefile
- [x] GoReleaser
- [x] GitHub Actions

---

## Core

- [x] Configuration
- [x] Logger
- [x] HTTP server
- [ ] SQLite
- [x] Goose
- [ ] sqlc

---

# Phase 2 — Notification Service

Deliver the first usable version.

- [ ] Event model
- [ ] Dispatcher
- [ ] Discord
- [ ] ntfy
- [ ] Email
- [ ] Webhook
- [ ] Forgejo integration

---

# Phase 3 — Production Ready

Improve stability and production readiness.

- [ ] Authentication
- [ ] Retry
- [ ] Metrics
- [ ] Graceful shutdown
- [ ] Logging improvements

---

# Phase 4 — Event Routing

Introduce routing capabilities.

- [ ] Rules
- [ ] Filters
- [ ] Multiple destinations
- [ ] Templates

---

# Phase 5 — Dashboard

Provide a web interface.

- [ ] Notification history
- [ ] Statistics
- [ ] Configuration UI

---

# Phase 6 — Event Platform

Expand Noticeal into a broader event platform.

- [ ] Additional connectors
- [ ] Plugins
- [ ] Public SDK

---

# Long-Term Vision

Noticeal begins as a lightweight notification service focused on Forgejo, distributed as a single self-hosted binary.

Over time, it will evolve into a generic event routing platform capable of receiving events from multiple systems, applying routing rules, and delivering notifications through a wide range of channels while remaining simple, lightweight and easy to run — no container required.