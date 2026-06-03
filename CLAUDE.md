# dispatch — Claude Context

Dispatch is a lightweight multi-machine agent hub. Workers (claude-monitor instances) phone home; the dashboard shows every machine and every session in one place. Terminal streams go **directly** from the browser to the worker — dispatch is never in the data path.

---

## Architecture

```
Worker (claude-monitor on any machine)
  │  POST /api/register every 30s (heartbeat)
  ▼
Dispatch Hub  ←──  Browser dashboard (session list, spawn, kill, resume)
                │
                │  GET /api/workers/{id}/ws/{name}  →  { ws_url: "ws://worker:7777/ws/name" }
                │
                └──────────────────────────────────────────────────────┐
                                                                       ▼
                                                          Browser WebSocket → Worker PTY
                                                          (hub NOT in the data path)
```

Key design choices:
- Hub is **stateless for data** — it only holds a worker registry in memory
- Worker IDs are stable 8-char hex derived from the worker's URL (SHA-256 truncated)
- Workers go **offline** after 90s without a heartbeat; evicted after 1 hour
- `registration_host` lets one port serve both a **public registration endpoint** (workers) and a **private dashboard** (SSO-protected)

---

## How to rebuild & deploy

The hub runs as a Docker container on port 8888. No systemd service — Docker manages restarts.

```bash
cd ~/apps/dispatch

# After any code change:
docker compose up --build -d

# View logs:
docker compose logs -f

# Restart without rebuild:
docker compose restart
```

`make install` builds the binary to `~/.local/bin/dispatch` — this is only useful if you want to run dispatch as a local binary (e.g., for testing). The production instance uses `docker compose`.

### Cross-compile for ARM64 (Raspberry Pi)
```bash
make arm64          # produces ./dispatch-arm64
```

---

## Config

Config lives in a Docker volume — not on the host filesystem.

```bash
# Read current config:
docker exec dispatch cat /root/.config/dispatch/config.json

# Write config:
docker exec dispatch sh -c 'cat > /root/.config/dispatch/config.json' << 'EOF'
{
  "port": 8888,
  "auth_token": "your-secret-token",
  "registration_host": "register.dispatch.kingdomofluna.com"
}
EOF
docker compose restart
```

| Field | Description |
|-------|-------------|
| `port` | HTTP port (default 8888) |
| `auth_token` | Workers must send `Authorization: Bearer <token>` to `/api/register`. Set the same value as `hub_auth_token` on each worker. |
| `registration_host` | When set, requests arriving on this hostname are restricted to `/api/register` and `/health` only — the dashboard is completely hidden from that host. This is the Pangolin public-registration resource hostname. |

---

## Pangolin setup (two resources, one port)

Dispatch uses a single port (8888) but two Pangolin resources with different SSO settings:

| Pangolin resource | SSO | Purpose |
|-------------------|-----|---------|
| `dispatch.kingdomofluna.com` | **On** | Dashboard — SSO-protected, browser only |
| `register.dispatch.kingdomofluna.com` | **Off** (public) | Worker registration — `auth_token` is the gate |

Both point to `localhost:8888` on aibo. The `registration_host` config tells dispatch to 404 any non-register request arriving from the public hostname.

---

## Connecting a new worker

A "worker" is any machine running claude-monitor. To add a machine:

1. Install claude-monitor on the machine (`make install` in `~/dev/claude-monitor/`)

2. Configure `~/.config/claude-monitor/config.json` on that machine:
   ```json
   {
     "port": 7777,
     "hub_url":        "https://register.dispatch.kingdomofluna.com",
     "hub_auth_token": "<dispatch auth_token>",
     "worker_url":     "https://agent.<machine>.kingdomofluna.com",
     "api_url":        "https://api.agent.<machine>.kingdomofluna.com"
   }
   ```
   - `worker_url` — the public URL the browser uses to open the session page (SSO-protected)
   - `api_url` — SSO-free hostname the hub uses for API calls and the browser uses for WebSocket connections; must be publicly reachable without SSO cookies

3. Add Pangolin resources for that worker:
   | Resource | SSO | Target |
   |----------|-----|--------|
   | `agent.<machine>.kingdomofluna.com` | **On** | `localhost:7777` |
   | `api.agent.<machine>.kingdomofluna.com` | **Off** | `localhost:7777` |

4. The worker appears in the dispatch dashboard within 30 seconds of starting claude-monitor.

---

## File structure

| File | What it does |
|------|-------------|
| `main.go` | Config, HTTP server, routing, `hostMiddleware` |
| `registry.go` | In-memory worker registry (upsert, list, get, mark-offline, evict) |
| `handlers.go` | HTTP handlers: register, workers, proxy, ws-info, dashboard, session page, health |
| `proxy.go` | `forwardToWorker` — HTTP proxy from hub to worker `/api/v1/` |
| `contract.go` | Shared types: `Registration`, `Worker`, `Session`, `workerView` |
| `ui.go` | Embedded HTML templates for dashboard, session page, multi-view |
| `dispatch_test.go` | Unit tests for registry and handler logic |
| `docker-compose.yml` | Production deployment — port 8888, named volume for config |
| `Dockerfile` | Multi-stage Go build |
| `Makefile` | `make build` (local binary), `make install` (copy to ~/.local/bin), `make arm64` |

---

## API reference

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/register` | Bearer token | Worker heartbeat |
| `GET` | `/api/workers` | — | List all workers |
| `GET` | `/api/workers/{id}` | — | Single worker detail |
| `ANY` | `/api/workers/{id}/{action}[/{name}]` | — | Proxy to worker `/api/v1/{action}/{name}` |
| `GET` | `/api/workers/{id}/ws/{name}` | — | Get direct WebSocket URL for PTY |
| `GET` | `/health` | none | Worker count (no auth required) |
| `GET` | `/` | — | Dashboard |
| `GET` | `/session/{workerid}/{name}` | — | Single session terminal page |
| `GET` | `/multi` | — | Multi-session split view |

Proxy actions forwarded to workers: `spawn`, `kill`, `restart`, `resume`, `output`, `instances`.

---

## Worker registration payload

```json
{
  "label":        "aibo",
  "url":          "https://agent.aibo.kingdomofluna.com",
  "api_url":      "https://api.agent.aibo.kingdomofluna.com",
  "worker_token": "<worker api token>",
  "version":      "2.0.0",
  "capabilities": ["claude", "agy", "terminal"],
  "sessions": [
    {
      "name":       "swift-fox",
      "status":     "running",
      "cli":        "claude",
      "dir":        "/home/user/dev/myproject",
      "summary":    "Refactoring auth module",
      "started_at": "2026-06-01T11:47:45Z"
    }
  ]
}
```
