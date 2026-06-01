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

## Deploying dispatch (Docker)

```bash
git clone https://github.com/sandeshkini/dispatch
cd dispatch
docker compose up -d
```

Config lives in a Docker volume at `/root/.config/dispatch/config.json`. Write it with:

```bash
docker exec dispatch sh -c 'mkdir -p /root/.config/dispatch && cat > /root/.config/dispatch/config.json' << 'EOF'
{
  "port": 8888,
  "auth_token": "your-secret-token-here",
  "registration_host": "register.dispatch.yourdomain.com"
}
EOF
docker compose restart
```

| Field | Description |
|-------|-------------|
| `auth_token` | Required on `/api/register` — set the same value as `hub_auth_token` on each worker |
| `registration_host` | Hostname of the public registration resource (Pangolin). Requests from this host are restricted to `/api/register` and `/health` only — the dashboard is not exposed through it |

### Pangolin setup

Create **two** Pangolin resources pointing to the same port 8888:

| Resource | SSO | Purpose |
|----------|-----|---------|
| `dispatch.yourdomain.com` | **On** | Dashboard — browser access, SSO protected |
| `register.dispatch.yourdomain.com` | **Off** (public) | Worker registration only — `auth_token` is the gate |

Add both as Cloudflare A records pointing to your VPS IP.

---

## Connecting a new worker

On the machine you want to add:

1. **Install claude-monitor** (the worker):
   ```bash
   git clone https://github.com/sandeshkini/claude-monitor
   cd claude-monitor
   go build -o ~/.local/bin/claude-monitor .
   ```

2. **Configure** `~/.config/claude-monitor/config.json`:
   ```json
   {
     "port": 7777,
     "hub_url":        "https://register.dispatch.yourdomain.com",
     "hub_auth_token": "<your dispatch auth_token>",
     "worker_url":     "https://agent.thismachine.yourdomain.com",
     "allowed_dirs":   ["~"]
   }
   ```

3. **Pangolin** — add a resource for this worker:
   - Hostname: `agent.thismachine.yourdomain.com`
   - Target: `localhost:7777`
   - SSO: **public** (browsers need direct WebSocket access; `/api/v1/` is token-protected by the worker itself)

4. **Systemd service** at `~/.config/systemd/user/claude-monitor.service`:
   ```ini
   [Unit]
   Description=Claude Monitor
   After=network-online.target

   [Service]
   ExecStart=%h/.local/bin/claude-monitor
   Restart=always
   RestartSec=5

   [Install]
   WantedBy=default.target
   ```
   ```bash
   systemctl --user enable --now claude-monitor
   loginctl enable-linger $USER
   ```

The worker appears in the dispatch dashboard within 30 seconds.

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
