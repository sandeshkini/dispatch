# dispatch

A lightweight hub that aggregates AI agent workers into a unified dashboard. Workers phone home — you see all your machines and sessions in one place.

---

## How it works

```
Worker  ──heartbeat──▶  Hub  ──coordinates──▶  Browser
                                                   │
                         └──────── direct WebSocket (PTY) ─────────┘
                                        (hub not in data path)
```

1. Workers POST a heartbeat to `/api/register` on startup and every 30s
2. Hub maintains a live registry — workers go offline after 90s without a heartbeat
3. Dashboard shows all machines, their capabilities, and active sessions
4. For terminal sessions, hub returns the direct WebSocket URL; browser connects straight to the worker — no relay

---

## Quick start

```bash
git clone https://github.com/your-org/dispatch
cd dispatch
make install        # builds + copies to ~/.local/bin/dispatch
dispatch            # starts on :8888
```

Config file is auto-created at `~/.config/dispatch/config.json`:

```json
{
  "port": 8888,
  "auth_token": ""
}
```

---

## Worker registration

Workers register by POSTing to `/api/register`. If you're building a compatible worker, start a heartbeat loop on startup:

```bash
# minimal test with curl
curl -X POST http://localhost:8888/api/register \
  -H "Content-Type: application/json" \
  -d '{"label":"my-machine","url":"http://100.x.x.x:7777","sessions":[]}'
```

Full registration schema and worker API contract: **[CONTRACT.md](CONTRACT.md)**

---

## Tailscale

Works best with Tailscale. Workers register with their `100.x.x.x` Tailscale IPs — the browser can then connect directly to any worker for terminal streaming regardless of NAT, no port forwarding needed.

---

## Raspberry Pi

Cross-compile for ARM64:

```bash
make arm64          # produces ./dispatch-arm64
```

Copy to Pi, install as a systemd service. The hub is stateless and uses ~10MB RAM — a Pi Zero 2W is sufficient.

---

## API

See [CONTRACT.md](CONTRACT.md) for full schema. Quick reference:

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/register` | Worker heartbeat |
| `GET`  | `/api/workers` | List all workers |
| `GET`  | `/api/workers/{id}` | Single worker |
| `ANY`  | `/api/workers/{id}/{action}` | Proxy to worker |
| `GET`  | `/api/workers/{id}/ws/{name}` | Get direct WebSocket URL |
| `GET`  | `/health` | Hub health (no auth) |
