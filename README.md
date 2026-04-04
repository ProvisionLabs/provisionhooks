# webhook-tester

Webhook Tester is a lightweight Go service for receiving webhook requests locally. It is meant to sit behind a tunnel such as ngrok so external services can call your machine while you inspect the payload.

## Overview

The current codebase provides a minimal HTTP receiver:

- `GET /health` for a basic health check
- `POST /hooks` for accepting webhook payloads
- In-memory storage for each received request

The broader direction described in [docs/overview.md](docs/overview.md) is to grow this into a self-hostable webhook inspection tool with realtime streaming and a browser UI.

## How It Works

1. Start the Go server on `localhost:8080`.
2. Expose it publicly with ngrok.
3. Point a webhook provider at the public URL.
4. Incoming requests are forwarded to `/hooks`.
5. The handler captures the method, headers, body, and timestamp in memory.

## Requirements

- Go 1.26 or newer
- ngrok, if you want to receive requests from external services

## Run

Start the server directly:

```bash
go run server/main.go
```

Or use the dev helper that starts the server and ngrok together:

```bash
task dev
```

On Windows, the helper uses `scripts/dev.ps1`. On Linux and macOS, it uses `scripts/dev.sh`.

## Endpoints

| Method | Path    | Description                                             |
| ------ | ------- | ------------------------------------------------------- |
| GET    | /health | Returns `200 OK`                                        |
| POST   | /hooks  | Reads the request body and stores the webhook in memory |

## Example

```bash
curl -X POST http://localhost:8080/hooks \
  -H "Content-Type: application/json" \
  -d '{"event":"test","value":123}'
```

The request is accepted with `200 OK`.

## Project Structure

```text
webhook-tester/
├── handler/        request handling logic
├── request/        webhook request model
├── server/         HTTP server entry point
├── store/          in-memory request storage
└── scripts/        local development helpers
```

## Non-Goals

- No database
- No authentication
- No persistence across restarts
- No request replay or UI in the current implementation

## Roadmap

The overview in [docs/overview.md](docs/overview.md) describes the intended next steps:

- Realtime request streaming
- A browser UI for inspecting captured webhooks
- Request history and other small quality-of-life features
