# Noticeal Roadmap

> This roadmap outlines the planned evolution of Noticeal.
>
> Development is organized into milestones rather than releases.
> Each milestone delivers a coherent piece of functionality while keeping the project stable and usable.

---

# Current Status

- [x] Project vision
- [x] Architecture
- [x] Repository created
- [ ] Development started

---

# Milestone 1 — Foundation

The goal of this milestone is to create a solid foundation for the project.

## Repository

- [ ] Initialize Go module
- [ ] Configure project layout
- [ ] Create Makefile
- [ ] Configure Docker
- [ ] Configure Docker Compose
- [ ] Configure GitHub Actions
- [ ] Configure GoReleaser

---

## Documentation

- [x] README
- [x] Architecture
- [x] Roadmap
- [ ] Contributing Guide
- [ ] Development Guide

---

## Core Packages

- [ ] api
- [ ] auth
- [ ] channel
- [ ] config
- [ ] connector
- [ ] event
- [ ] logger
- [ ] router
- [ ] rule
- [ ] storage
- [ ] version

---

## Configuration

- [ ] Viper
- [ ] YAML configuration
- [ ] Environment variables
- [ ] Configuration validation

---

## Logging

- [ ] Zap
- [ ] Request logger
- [ ] Structured logs

---

# Milestone 2 — HTTP Server

The objective is to expose a production-ready HTTP server.

## Server

- [ ] Chi router
- [ ] Graceful shutdown
- [ ] Recovery middleware
- [ ] Request logging
- [ ] Compression
- [ ] CORS

---

## Endpoints

- [ ] GET /health
- [ ] GET /version
- [ ] GET /metrics

---

# Milestone 3 — Event Model

The Event becomes the central object of the application.

## Event

- [ ] Event model
- [ ] Event validation
- [ ] JSON serialization

---

## REST Connector

- [ ] POST /api/v1/events
- [ ] Request validation
- [ ] Error handling

---

## CLI Connector

- [ ] noticeal send
- [ ] noticeal validate
- [ ] noticeal version

---

# Milestone 4 — Routing Engine

Build the core routing engine.

## Rules

- [ ] Match by source
- [ ] Match by event type
- [ ] Match by status
- [ ] Match by metadata

---

## Router

- [ ] Sequential delivery
- [ ] Parallel delivery
- [ ] Error propagation

---

# Milestone 5 — Channels

Implement the first delivery channels.

## Webhook

- [ ] HTTP POST
- [ ] Custom headers

---

## Discord

- [ ] Webhooks
- [ ] Rich embeds

---

## ntfy

- [ ] Topics
- [ ] Priority
- [ ] Tags

---

## Email

- [ ] SMTP
- [ ] HTML
- [ ] Plain text

---

# Milestone 6 — Persistence

Persist routing information.

## SQLite

- [ ] Goose
- [ ] sqlc
- [ ] Repositories

---

## Storage

- [ ] Events
- [ ] Deliveries
- [ ] Rules
- [ ] Channels

---

# Milestone 7 — Reliability

Improve delivery robustness.

## Retry

- [ ] Retry policy
- [ ] Exponential backoff

---

## Dead Letter Queue

- [ ] Failed deliveries
- [ ] Replay

---

## Delivery Reports

- [ ] Success
- [ ] Failure
- [ ] Latency

---

# Milestone 8 — Templates

Support reusable notification templates.

## Templates

- [ ] Go templates
- [ ] Variables
- [ ] Markdown
- [ ] Conditional rendering

Example:

```
🚀 {{ .Title }}

Repository:
{{ .Metadata.repository }}

Status:
{{ .Status }}
```

---

# Milestone 9 — Connectors

Expand event ingestion.

## GitHub

- [ ] Webhooks
- [ ] Workflow events
- [ ] Release events

---

## Forgejo

- [ ] Webhooks
- [ ] Workflow events
- [ ] Release events

---

## GitLab

- [ ] Pipeline events
- [ ] Release events

---

## Custom Connectors

- [ ] MQTT
- [ ] Kafka
- [ ] NATS

---

# Milestone 10 — Dashboard

Build a web interface.

## Monitoring

- [ ] Events
- [ ] Deliveries
- [ ] Failures
- [ ] Statistics

---

## Administration

- [ ] Rules
- [ ] Channels
- [ ] Connectors

---

# Milestone 11 — Security

Secure the platform.

## Authentication

- [ ] API Keys
- [ ] JWT

---

## Authorization

- [ ] RBAC

---

## Security

- [ ] Rate limiting
- [ ] Request signatures

---

# Milestone 12 — Multi-Project

Support multiple independent projects.

## Projects

- [ ] Projects
- [ ] Project configuration
- [ ] Isolation

---

# Milestone 13 — Enterprise

## Storage

- [ ] PostgreSQL
- [ ] MySQL

---

## Queue

- [ ] Background workers

---

## Scheduling

- [ ] Delayed events
- [ ] Scheduled deliveries

---

## SDKs

- [ ] Go
- [ ] JavaScript

---

## Kubernetes

- [ ] Helm Chart
- [ ] Kubernetes Operator

---

# Future Channels

- [ ] Slack
- [ ] Telegram
- [ ] Microsoft Teams
- [ ] Mattermost
- [ ] Matrix
- [ ] Signal
- [ ] WhatsApp

---

# Future Connectors

- [ ] Jenkins
- [ ] Drone CI
- [ ] ArgoCD
- [ ] Docker
- [ ] Kubernetes
- [ ] Cron Jobs

---

# Future Integrations

- [ ] Prometheus
- [ ] Grafana
- [ ] OpenTelemetry
- [ ] Sentry
- [ ] PagerDuty
- [ ] Opsgenie

---

# Long-Term Vision

Noticeal aims to become the universal event routing layer for self-hosted infrastructures.

Applications should no longer integrate directly with notification systems.

Instead, they emit Events.

Noticeal becomes responsible for:

- Receiving events
- Normalizing events
- Evaluating routing rules
- Delivering events
- Tracking deliveries
- Providing observability

The architecture should remain simple, modular, and extensible while keeping the core independent from any external platform.