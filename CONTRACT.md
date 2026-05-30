# Dispatch Hub API & Integration Contract

This document specifies the communication contract between the Dispatch Hub, workers, and clients (e.g., browsers, dashboards, and CLI tools).

---

## 1. Worker Registration (`POST /api/register`)

Workers must announce their presence on startup and send periodic heartbeats to this endpoint.

* **Endpoint:** `/api/register`
* **Method:** `POST`
* **Heartbeat Interval:** Recommended every **30 seconds** (the Hub marks workers offline if no heartbeat is received for **90 seconds**).
* **Response:** On success, returns `200 OK` with a stable, unique 8-character hex ID derived from the worker's URL.

### JSON Schema
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Registration",
  "type": "object",
  "required": ["label", "url", "sessions"],
  "properties": {
    "label": {
      "type": "string",
      "description": "Human-readable label for the worker"
    },
    "url": {
      "type": "string",
      "format": "uri",
      "description": "The base HTTP URL of the worker (e.g. http://10.0.0.5:7777)"
    },
    "version": {
      "type": "string",
      "description": "Optional version tag of the worker software"
    },
    "capabilities": {
      "type": "array",
      "items": { "type": "string" },
      "description": "Optional list of tools available on this worker (e.g. ['claude', 'pi', 'terminal'])"
    },
    "sessions": {
      "type": "array",
      "items": { "$ref": "#/definitions/Session" },
      "description": "Current sessions on the worker"
    }
  },
  "definitions": {
    "Session": {
      "type": "object",
      "required": ["name", "status", "cli", "dir"],
      "properties": {
        "name":       { "type": "string" },
        "status":     { "type": "string", "enum": ["running", "stopped"] },
        "cli":        { "type": "string" },
        "dir":        { "type": "string" },
        "summary":    { "type": "string" },
        "started_at": { "type": "string", "format": "date-time" }
      }
    }
  }
}
```

### Request Example
```json
{
  "label": "aibo",
  "url": "http://100.x.x.x:7777",
  "version": "0.1.0",
  "capabilities": ["claude", "pi", "terminal"],
  "sessions": [
    {
      "name": "swift-fox",
      "status": "running",
      "cli": "claude",
      "dir": "/home/user/dev/myproject",
      "summary": "Refactoring auth module",
      "started_at": "2026-05-30T11:47:45Z"
    }
  ]
}
```

### Response Example
```json
{ "id": "c0a80005", "status": "registered" }
```

---

## 2. Worker API Endpoints (Hub → Worker)

Each worker must implement these endpoints. The hub proxies client commands directly to them.

| Method | Path | Description |
|--------|------|-------------|
| `GET`  | `/api/v1/instances` | List all sessions |
| `POST` | `/api/v1/spawn` | Spawn a new session |
| `POST` | `/api/v1/kill/{name}` | Kill a session |
| `POST` | `/api/v1/restart/{name}` | Restart a session |
| `POST` | `/api/v1/resume/{name}` | Resume a stopped session |
| `GET`  | `/api/v1/output/{name}?lines=100` | Get last N lines of output |
| `WS`   | `/ws/{name}` | PTY WebSocket stream (browser connects directly) |

### Spawn request body
```json
{
  "dir":   "/path/to/project",
  "cli":   "claude",
  "flags": ["--dangerously-skip-permissions"],
  "name":  "swift-fox"
}
```

---

## 3. Hub API Endpoints (Dashboard / Tools → Hub)

All endpoints except `/health` require `Authorization: Bearer <token>` when `auth_token` is set in config.

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/register` | Worker heartbeat / registration |
| `GET`  | `/api/workers` | List all workers |
| `GET`  | `/api/workers/{id}` | Single worker detail |
| `ANY`  | `/api/workers/{id}/{action}[/{name}]` | Proxy action to worker |
| `GET`  | `/api/workers/{id}/ws/{name}` | Get direct WebSocket URL for a session |
| `GET`  | `/health` | Hub health (unauthenticated) |

### Worker list response
```json
[
  {
    "id": "c0a80005",
    "label": "aibo",
    "url": "http://100.x.x.x:7777",
    "version": "0.1.0",
    "capabilities": ["claude", "pi", "terminal"],
    "sessions": [ ... ],
    "online": true,
    "last_seen": "2026-05-30T11:47:50Z"
  }
]
```

### WebSocket info response
```json
{
  "ws_url":     "ws://100.x.x.x:7777/ws/swift-fox",
  "worker_url": "http://100.x.x.x:7777"
}
```

### Health response
```json
{
  "status": "ok",
  "workers_total": 4,
  "workers_online": 3,
  "workers_offline": 1
}
```

---

## 4. Direct WebSocket Flow

PTY streaming bypasses the hub entirely. The hub only provides the coordinates.

```
Browser → GET /api/workers/{id}/ws/{name}
Hub     → { "ws_url": "ws://worker:7777/ws/session-name" }
Browser → WebSocket directly to worker (no hub relay)
```

This keeps the hub stateless for data paths and removes it as a latency bottleneck for terminal sessions.

---

## 5. Minimal Heartbeat Example (Python)

```python
import json, time, urllib.request

HUB_URL    = "http://dispatch-host:8888/api/register"
WORKER_URL = "http://this-machine:7777"
AUTH_TOKEN = ""  # match hub config auth_token if set

def heartbeat():
    payload = json.dumps({
        "label": "my-worker",
        "url": WORKER_URL,
        "version": "0.1.0",
        "capabilities": ["terminal"],
        "sessions": []
    }).encode()
    headers = {"Content-Type": "application/json"}
    if AUTH_TOKEN:
        headers["Authorization"] = f"Bearer {AUTH_TOKEN}"
    req = urllib.request.Request(HUB_URL, data=payload, headers=headers, method="POST")
    with urllib.request.urlopen(req, timeout=10) as r:
        print(r.read().decode())

while True:
    try:    heartbeat()
    except Exception as e: print(f"heartbeat failed: {e}")
    time.sleep(30)
```
