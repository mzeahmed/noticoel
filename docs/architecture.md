# Noticeal Architecture

> **Noticeal is an event router. Notifications are just one destination.**

---

# Introduction

Noticeal is an open-source event routing platform.

Its purpose is to receive events from different producers, normalize them into a common format, evaluate routing rules, and deliver them to one or more destinations.

Unlike traditional notification services, Noticeal does not focus on a specific provider or platform. Instead, it acts as the missing layer between systems producing events and systems consuming them.

```
                  +----------------------+
                  |   Event Producers    |
                  +----------------------+
                     REST API / CLI
                           │
                           ▼
                   +---------------+
                   |   Noticeal    |
                   +---------------+
                           │
                 Rule Evaluation Engine
                           │
                           ▼
                   +---------------+
                   |    Router     |
                   +---------------+
                     │     │     │
                     ▼     ▼     ▼
                 Discord ntfy Email
```

---

# Design Principles

Noticeal is built around a few core principles.

## Simplicity

Every package should have a single responsibility.

Complexity should emerge from composition rather than inheritance or deeply nested abstractions.

---

## Self-hosted First

Noticeal is designed primarily as a self-hosted application.

Typical deployments include:

- Local development
- VPS
- Home server
- Docker
- Kubernetes

A public SaaS offering is intentionally outside the scope of the core project.

---

## Platform Agnostic

Noticeal should never depend on GitHub, Forgejo, GitLab, Discord, Slack or any other vendor.

External systems are considered integrations.

The core should remain completely independent.

---

## Extensible

Every integration should be replaceable.

Adding support for a new destination should never require modifications to the routing engine.

---

## Event Driven

Everything inside Noticeal revolves around events.

An event is the universal language spoken by every package.

---

# High-Level Architecture

```
                Connectors

         REST API       CLI
             │           │
             └─────┬─────┘
                   ▼
              Event Parser
                   │
                   ▼
              Event Model
                   │
                   ▼
            Rule Evaluation
                   │
                   ▼
                 Router
                   │
        ┌──────────┼──────────┐
        ▼          ▼          ▼
    Discord      ntfy      Webhook
```

---

# Core Concepts

## Connector

A connector is responsible for receiving events from the outside world.

The initial version of Noticeal provides only two connectors:

- REST API
- CLI

Future versions may introduce additional connectors such as:

- GitHub Webhooks
- Forgejo Webhooks
- GitLab Webhooks
- MQTT
- Kafka
- NATS

Connectors should never contain business logic.

Their only responsibility is to convert external data into an Event.

---

## Event

The Event is the central model of the application.

Every connector must convert its incoming payload into a normalized Event.

Every package inside Noticeal works exclusively with Events.

Example:

```go
type Event struct {
	ID         string
	Source     string
	Type       string
	Status     string
	Title      string
	Message    string
	Metadata   map[string]string
	Timestamp  time.Time
}
```

This guarantees that the routing engine never depends on external payload formats.

---

## Rule Engine

The Rule Engine decides **what should happen**.

Example:

```
If:

status == failure

Then:

Discord
Email
```

Rules should remain independent from channel implementations.

The Rule Engine only determines which channels should receive an event.

---

## Router

The Router is responsible for execution.

It receives:

- one Event
- a list of Channels

Its responsibilities include:

- sequential or parallel delivery
- error propagation
- retry policy
- delivery lifecycle

The Router should never evaluate business rules.

---

## Channel

A Channel is responsible for delivering an event to an external system.

Examples:

- Discord
- ntfy
- Email
- Slack
- Teams
- Webhook

Every channel implements the same interface.

```go
type Channel interface {
	Name() string
	Send(ctx context.Context, event Event) error
}
```

---

# Project Structure

```
noticeal/

cmd/
    noticeal/

internal/

    api/

    auth/

    channel/

    config/

    event/

    logger/

    router/

    rule/

    storage/

    connector/

    version/

migrations/

docs/

examples/

scripts/

Dockerfile

docker-compose.yml

Makefile

README.md
```

---

# Package Responsibilities

## api

HTTP server.

Responsibilities:

- Routing
- Request validation
- JSON serialization
- Authentication
- Error responses

No business logic.

---

## connector

Receives events.

Examples:

- REST API
- CLI

Future:

- GitHub
- Forgejo
- GitLab

---

## event

Defines the Event model.

No persistence.

No routing.

No networking.

---

## rule

Evaluates routing rules.

Input:

```
Event
```

Output:

```
[]Channel
```

---

## router

Executes deliveries.

Input:

```
Event

+

Channels
```

Responsibilities:

- Send events
- Retry
- Collect results

---

## channel

Contains every delivery implementation.

Each channel is completely independent.

---

## storage

Persistence layer.

Responsibilities:

- repositories
- database access

Never contains business logic.

---

## config

Application configuration.

Supported sources:

- YAML
- Environment variables
- CLI flags

---

## logger

Centralized structured logging.

---

## auth

Authentication.

Initially:

API Keys.

Future:

JWT

RBAC

---

## version

Application version.

Injected during build.

---

# Configuration

Version 1 uses a YAML configuration.

Example:

```yaml
server:
  host: 0.0.0.0
  port: 8080

channels:

  discord:
    webhook: https://...

  ntfy:
    topic: noticeal

rules:

  - match:
      status: failure

    send:
      - discord
      - ntfy
```

---

# Persistence

Noticeal stores only operational data.

Examples:

- Events
- Deliveries
- Channels
- Rules

Business logic must never depend on persistence.

Removing the database should only remove history, never the routing capability.

---

# Deployment

Noticeal is designed to run in multiple environments.

Examples:

```
Developer machine

↓

Docker Compose
```

```
VPS

↓

Docker
```

```
Kubernetes

↓

Deployment
```

The application should not assume where it is running.

---

# REST API

Version 1 exposes a minimal API.

```
POST /api/v1/events
```

Receive an Event.

---

```
GET /health
```

Health check.

---

```
GET /version
```

Version information.

---

```
GET /metrics
```

Prometheus metrics.

---

# Development Principles

- Small packages.
- Single responsibility.
- Composition over inheritance.
- No framework.
- Keep dependencies minimal.
- Interfaces only at architectural boundaries.
- Prefer explicit code over magic.
- The Event is the only shared model.
- Connectors never deliver events.
- Channels never evaluate rules.
- The Router never knows external services.
- The Rule Engine never sends notifications.

---

# Long-Term Vision

Noticeal aims to become the standard event routing layer for self-hosted infrastructures.

Applications should no longer integrate directly with dozens of notification providers.

Instead, they emit Events.

Noticeal becomes responsible for:

- Event normalization
- Rule evaluation
- Routing
- Delivery
- Retry policies
- Observability

Notifications are only one possible outcome.

The architecture should remain flexible enough to support future integrations such as:

- Message brokers
- Workflow engines
- Serverless platforms
- Internal APIs
- Custom extensions

without changing the core design.