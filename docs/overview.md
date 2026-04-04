# Webhook Tester — Project Overview

A lightweight, self-hostable webhook inspection tool written in Go.
Run it locally, expose it publicly via ngrok, and inspect incoming webhook requests in real time — no account or deployment required.

---

## The Problem

When building integrations (Stripe, GitHub, Slack, etc.), developers need to see what payload an external service sends. The problem: your local machine has no public URL, so external services can't reach it directly.

This tool solves the inspection side — capturing, storing, and displaying incoming webhook requests in real time. It uses ngrok to create the public tunnel, so you can receive webhooks locally without any deployment. Think of it as the missing UI layer on top of ngrok for webhook debugging.

Tools like `webhook.site` or `RequestBin` do something similar but are closed source, paywalled, or managed. This is the self-hostable alternative.

---

## How It Works

```
[External Service]
       |
       | POST https://abc123.ngrok.io/hooks
       ↓
[ngrok tunnel]                     ← forwards public traffic to localhost
       |
       | POST http://localhost:8080/hooks
       ↓
[Go HTTP Server]
       |
       | stores request in memory
       ↓
[SSE Stream] /hooks/events    ← pushes event to browser instantly
       |
       ↓
[Browser UI]                       ← user sees payload in real time
```

1. You start the Go server locally on `localhost:8080`
2. You run ngrok to get a public URL that tunnels to your local server
3. You paste the ngrok URL into any external service as the webhook destination
4. External service fires a POST request to the ngrok URL
5. ngrok forwards it to your local Go server transparently
6. Go server captures the full request (method, headers, body)
7. Server pushes the captured data to the browser via SSE (no refresh needed)
8. You inspect the payload in real time

---

## Core Concepts

### Single Receiver Endpoint

The current implementation receives webhooks at a single endpoint (`POST /hooks`).
Requests are stored in-memory and shown in the UI.

### SSE (Server-Sent Events)

SSE is a simple protocol for pushing data from server to browser over a persistent HTTP connection. When a webhook arrives, the Go server instantly notifies the browser — no polling, no WebSocket complexity.

```
GET /hooks/events
← data: {"method":"POST","body":"..."}   ← server pushes this
← data: {"method":"POST","body":"..."}   ← and this, whenever a new webhook arrives
```

### In-Memory Storage

No database needed. Requests are stored in an in-memory slice protected by mutexes. Connected SSE clients are tracked in-memory via channels. This keeps the project simple and fast.

### Concurrency

Multiple webhooks can arrive at the same time. Go handles this naturally with goroutines — each incoming HTTP request runs in its own goroutine, so nothing blocks.

---

## Project Structure

```
webhook-tester/
├── server/
│   ├── main.go              # Entry point, starts HTTP server
├── handler/
│   ├── webhook.go           # POST /hooks — receives and stores requests
│   ├── stream.go            # GET /hooks/events — SSE stream
│   └── ui.go                # GET /hooks — serves the browser UI
├── store/
│   └── store.go             # In-memory storage for requests and SSE clients
├── request/
│   └── request.go           # Structs representing captured and API response requests
└── static/
    └── index.html           # Simple frontend to display captured requests
```

---

## API Endpoints

| Method | Path            | Description                                 |
| ------ | --------------- | ------------------------------------------- |
| `POST` | `/hooks`        | Receives a webhook and stores it            |
| `GET`  | `/hooks`        | Browser UI for the requests                 |
| `GET`  | `/hooks/data`   | JSON history of captured requests           |
| `GET`  | `/hooks/events` | SSE stream — pushes new requests to browser |

---

## Go Concepts Demonstrated

| Concept                           | Where used                                     |
| --------------------------------- | ---------------------------------------------- |
| HTTP server (`net/http`)          | All handlers                                   |
| Goroutines                        | Each request handled concurrently              |
| Channels                          | Notifying SSE clients when new webhook arrives |
| Mutex-protected in-memory storage | `store/store.go`                               |
| SSE (manual HTTP streaming)       | `stream.go`                                    |
| Structs and JSON tags             | `request/request.go`                           |
| JSON encoding                     | Serializing captured requests                  |

---

## Build Milestones

### v0.1 — Receive & Store

- [x] `POST /hooks` captures method, headers, body
- [x] Store in memory
- [x] Return 200 OK to the sender

### v0.2 — Real-Time Stream

- [x] `GET /hooks/events` opens SSE connection
- [x] Server pushes new requests to connected clients via channel
- [x] Handle client disconnect gracefully

### v0.3 — Browser UI

- [x] HTML page connects to the SSE stream
- [x] Displays each captured request with syntax highlighting
- [x] Shows method, timestamp, headers, and body

### v0.4 — Polish

- [x] Request history in memory
- [ ] Copy URL button (copies the ngrok URL to clipboard)
- [ ] Print the full public URL on startup if ngrok is detected

---

## Running Locally

### 1. Install ngrok

Download from [ngrok.com/download](https://ngrok.com/download) or via package manager:

```bash
# macOS
brew install ngrok

# Linux
snap install ngrok
```

Create a free account at [ngrok.com](https://ngrok.com) and authenticate:

```bash
ngrok config add-authtoken <your-token>
```

---

### 2. Start the Go server

```bash
git clone https://github.com/youruser/webhook-tester
cd webhook-tester
go run server/main.go
# Server running at http://localhost:8080
```

---

### 3. Expose it publicly with ngrok

Open a second terminal:

```bash
ngrok http 8080
```

You'll see output like:

```
Forwarding   https://abc123.ngrok.io -> http://localhost:8080
```

That `https://abc123.ngrok.io` is your public URL. Any request to it is forwarded to your local server.

---

### 4. Use your public webhook URL

Use the ngrok URL with the webhook endpoint:

```
https://abc123.ngrok.io/hooks
```

Paste this into GitHub, Stripe, Slack — anywhere that sends webhooks. When they fire, you'll see the request appear in your browser in real time.

---

### 5. Test it yourself with curl

You don't need an external service to develop. Simulate a webhook anytime:

```bash
curl -X POST https://abc123.ngrok.io/hooks \
  -H "Content-Type: application/json" \
  -d '{"event": "payment.received", "amount": 100}'
```

> **Note:** ngrok free tier gives you a new random URL every time you restart it. For a stable URL during development, keep ngrok running in the background.

---

## Non-Goals (keep it simple)

- No database
- No authentication
- No replay functionality (yet)
- No persistent sessions across restarts

---

## Inspiration

- [webhook.site](https://webhook.site) — great tool, but closed source and paywalled
- [RequestBin](https://requestbin.com) — similar concept, also managed
- This project aims to be the self-hostable, zero-dependency alternative
