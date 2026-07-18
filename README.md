<div align="center">

<img src="assets/logo.png" alt="Noticeal Logo" width="220">

<br>

**An open-source event routing platform for self-hosted infrastructures.**

Receive events. Evaluate rules. Deliver anywhere.

<br>

[//]: # ([![CI]&#40;https://github.com/mzeahmed/noticeal/actions/workflows/ci.yml/badge.svg&#41;]&#40;https://github.com/mzeahmed/noticeal/actions/workflows/ci.yml&#41;)
[![Go Report Card](https://goreportcard.com/badge/github.com/mzeahmed/noticeal)](https://goreportcard.com/report/github.com/mzeahmed/noticeal)
[![Go Reference](https://pkg.go.dev/badge/github.com/mzeahmed/noticeal.svg)](https://pkg.go.dev/github.com/mzeahmed/noticeal)
[![Release](https://img.shields.io/github/v/release/mzeahmed/noticeal)](https://github.com/mzeahmed/noticeal/releases)
[![License](https://img.shields.io/github/license/mzeahmed/noticeal)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mzeahmed/noticeal)](app/go.mod)

[//]: # ([![Coverage]&#40;https://codecov.io/gh/mzeahmed/noticeal/branch/main/graph/badge.svg&#41;]&#40;https://codecov.io/gh/mzeahmed/noticeal&#41;)
[//]: # ([![Open Issues]&#40;https://img.shields.io/github/issues/mzeahmed/noticeal&#41;]&#40;https://github.com/mzeahmed/noticeal/issues&#41;)
[//]: # ([![Pull Requests]&#40;https://img.shields.io/github/issues-pr/mzeahmed/noticeal&#41;]&#40;https://github.com/mzeahmed/noticeal/pulls&#41;)
[//]: # ([![Downloads]&#40;https://img.shields.io/github/downloads/mzeahmed/noticeal/total&#41;]&#40;https://github.com/mzeahmed/noticeal/releases&#41;)

</div>

---

## Overview

Noticeal is an event routing platform built for self-hosted environments.

Instead of integrating every application with multiple notification providers, applications simply send events to Noticeal.

Noticeal evaluates routing rules and delivers events to the appropriate destinations.

The core remains independent from external platforms, making integrations simple, modular and extensible.

---

## Why Noticeal?

Modern infrastructures generate events everywhere:

- CI/CD pipelines
- Deployment platforms
- Monitoring systems
- Internal applications
- Automation tools

Most applications implement notification logic themselves.

Noticeal centralizes this responsibility.

Applications emit events.

Noticeal decides what happens next.

---

## Architecture

```text
            Connectors

      REST API        CLI
          │            │
          └──────┬─────┘
                 ▼
             Event Model
                 │
                 ▼
         Processing Engine
        ├─────────────────┐
        │                 │
        ▼                 ▼
 Rule Evaluation     Delivery
        │                 │
        └────────┬────────┘
                 ▼
              Channels
```

---

## Features

### Event Processing

- Event-driven architecture
- Rule-based routing
- Normalized event model
- Processing engine
- Structured logging

### Connectors

- REST API
- CLI

### Channels

- Webhook
- Discord
- ntfy
- Email

---

## Example

Send an event using the REST API.

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "source": "forgejo",
    "type": "workflow",
    "status": "success",
    "title": "Deployment completed",
    "message": "BookingApp deployed successfully"
  }'
```

---

## Philosophy

Noticeal follows a simple principle:

> Applications should emit events. Noticeal should decide what to do with them.

This keeps applications focused on their own business logic while Noticeal handles routing and delivery.

---

## Documentation

- [Architecture](docs/architecture.md)
- [Roadmap](docs/roadmap.md)

More documentation will be added as the project evolves.

---

## Contributing

Contributions are welcome.

If you find a bug, have an idea, or want to contribute, feel free to open an issue or submit a pull request.

---

## Project Status

⚠️ Noticeal is currently under active development.

The project is not production-ready yet and the API may change before the first stable release.

---

## License

Noticeal is released under the MIT License.

See the [LICENSE](LICENSE) file for details.