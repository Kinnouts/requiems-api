## Infrastructure Guide

This document explains how to run Requiem API locally for development and how to
deploy it to a single VPS.

---

## 1. Local Development Setup

### 1.1 Requirements (Quick Start)

**Recommended:**

- Docker and Docker Compose (required for the main stack)
- Bun or Node.js 18+ (only for edge-auth development)

**Optional (Advanced):**

- Go 1.22+ (only if running API without Docker)

### 1.2 Run everything with Docker (RECOMMENDED)

From the project root:

```bash
cd infra/docker
docker compose up --build
```

This starts:

- `api` – Go backend on internal port `8080` (exposed as `localhost:8080`)
- `db` – PostgreSQL (`requiem` / `requiem` / `requiem`, exposed as
  `localhost:5432`)
- `redis` – Redis for future queues/cache

Once the stack is up:

- Health check: `http://localhost:8080/healthz`

> Note: Caddy is mainly for the VPS setup. For local development you can hit the
> API directly on `localhost:8080`.

### 1.3 Edge Auth (Cloudflare Worker) - Local Development

**Note:** The Worker cannot run in Docker. Use wrangler for local development.

```bash
cd apps/edge-auth

# Install dependencies
bun install

# Start local dev server
bun run dev
# Worker runs on http://localhost:8787
```

Set environment variables in `wrangler.toml` or as secrets:

- `BACKEND_URL` – e.g., `http://localhost:8080`
- `BACKEND_SECRET` – shared secret for backend authentication

### 1.4 Run the API directly (without Docker) - Advanced

If you prefer to run the Go server directly:

1. Make sure PostgreSQL is running locally (matching the default DSN or set
   `DATABASE_URL`).
2. From the project root:

```bash
export DATABASE_URL="postgres://requiem:requiem@localhost:5432/requiem?sslmode=disable"
go run ./apps/api
```

The API listens on `:8080` by default:

- `http://localhost:8080/healthz`
- `http://localhost:8080/v1/advice`

### 1.5 Hybrid dev workflow (Docker infra + local Go) - Advanced

For advanced users who want hot reload without Docker overhead, run **Postgres
and Redis in Docker**, and the Go API locally with a watcher:

- Start infra only:

```bash
cd infra/docker
docker compose up db redis
```

- In another terminal, from the repo root:

```bash
export DATABASE_URL="postgres://requiem:requiem@localhost:5432/requiem?sslmode=disable"
go run ./apps/api
```

Or, with a hot-reload tool like `air`:

```bash
export DATABASE_URL="postgres://requiem:requiem@localhost:5432/requiem?sslmode=disable"
air
```

The API is still available on:

- `http://localhost:8080/healthz`
- `http://localhost:8080/v1/advice`

---

## 2. VPS Deployment Guide

Target: single VPS (Docker + Docker Compose) running API, Postgres, Redis, and
Caddy. Cloudflare Worker sits in front and forwards authorized traffic.

### 2.1 Prepare the VPS

#### Recommended Hetzner Configuration

- **Server Type:** CPX21 or better (3 vCPU, 4GB RAM, 80GB SSD) - ~€7.95/month
- **Location:** Choose based on target audience
  - Nuremberg (nbg1) - Europe/Germany
  - Helsinki (hel1) - Northern Europe
  - Ashburn (ash) - US East
- **Image:** Ubuntu 24.04 LTS
- **Networking:** Enable IPv4 (required), IPv6 (optional)

#### Initial Setup

1. **Create VPS** in Hetzner Cloud Console

2. **SSH into the server:**

```bash
ssh root@YOUR_VPS_IP
```

3. **Update system:**

```bash
apt update && apt upgrade -y
```

4. **Install Docker:**

```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
rm get-docker.sh

# Install Docker Compose plugin
apt install docker-compose-plugin -y

# Verify installation
docker --version
docker compose version
```

5. **Configure firewall:**

```bash
apt install ufw
ufw allow 22/tcp   # SSH
ufw allow 80/tcp   # HTTP
ufw allow 443/tcp  # HTTPS
ufw enable
ufw status
```

6. **Clone repository:**

```bash
git clone https://github.com/bobadilla-tech/requiems-api.git
cd requiems-api/infra/docker
```

### 2.2 Start the stack

From `infra/docker`:

```bash
docker compose up -d --build
```

Services:

- `api` – Go API (port `8080` in the Docker network)
- `db` – PostgreSQL
- `redis` – Redis
- `caddy` – HTTPS reverse proxy on ports `80` and `443`

You can check logs with:

```bash
docker compose logs -f api
docker compose logs -f caddy
```

### 2.3 Caddy configuration (HTTPS + reverse proxy)

Caddy is configured via `infra/caddy/Caddyfile`. Example:

```bash
api.yourdomain.com {
  encode gzip

  reverse_proxy api:8080
}
```

What this does:

- Terminates HTTPS and manages TLS certificates automatically.
- Proxies `https://api.yourdomain.com` → `api:8080` (the API container).

### 2.4 DNS Setup (Cloudflare)

#### Add Domain to Cloudflare

1. **Add your domain** to Cloudflare (if not already):
   - Sign up/login at [cloudflare.com](https://cloudflare.com)
   - Add your domain
   - Update nameservers at your registrar to Cloudflare's nameservers

2. **Create DNS records:**

| Type | Name     | Content     | Proxy Status           | TTL  |
| ---- | -------- | ----------- | ---------------------- | ---- |
| A    | api      | YOUR_VPS_IP | DNS only (grey cloud)  | Auto |
| A    | internal | YOUR_VPS_IP | DNS only (grey cloud)  | Auto |
| A    | @        | YOUR_VPS_IP | Proxied (orange cloud) | Auto |

**Important:**

- Use "DNS only" (grey cloud) for `api` and `internal` subdomains
- This allows Caddy to obtain Let's Encrypt certificates directly
- The root domain (@) can be proxied for DDoS protection

3. **Verify DNS propagation:**

```bash
dig api.yourdomain.com
dig internal.yourdomain.com
# Or use: nslookup api.yourdomain.com
```

4. **Wait for DNS** to propagate (5-30 minutes)

5. **Cloudflare Worker routing** (for public API):
   - In Cloudflare dashboard, go to Workers & Pages
   - Deploy your edge-auth worker
   - Add a route: `api.yourdomain.com/*` → your worker
   - Set `BACKEND_URL` environment variable to `https://internal.yourdomain.com`

#### Domain Architecture

```
User Request
    ↓
api.yourdomain.com (Cloudflare Worker - handles auth)
    ↓ (validates requiems-api-key, forwards with X-Backend-Secret)
internal.yourdomain.com (Direct to VPS)
    ↓
Caddy on VPS (HTTPS termination)
    ↓
Go API container (port 8080)
```

With DNS + Caddy running:

- `https://api.yourdomain.com/healthz` → should work (via Worker)
- `https://internal.yourdomain.com/healthz` → should work (direct to VPS)

### 2.5 Cloudflare Worker configuration (edge auth)

1. Deploy `apps/edge-auth/index.ts` as a Worker in your Cloudflare account.
2. Configure environment variables/secrets:
   - `BACKEND_ORIGIN` – `https://api.yourdomain.com`
   - `BACKEND_SECRET` – a strong secret, used by clients in `requiems-api-key`.
3. Route your public API endpoint to the Worker (e.g.
   `https://v1.yourdomain.com/*` → Worker).

Request flow in production:

1. Client → Worker (with `requiems-api-key`).
2. Worker validates `requiems-api-key`.
3. Worker forwards to `https://api.yourdomain.com/...`.
4. Cloudflare → Caddy on VPS → Go API.

---

## 3. Environment Variables Summary

- **API container**
  - `PORT` – API listen port (default `8080`).
  - `DATABASE_URL` – Postgres DSN (set by Docker Compose for container, or
    manually in local dev).
  - `REDIS_URL` – Redis URL for future queues/cache.

- **Cloudflare Worker**
  - `BACKEND_ORIGIN` – Base URL of the API behind Caddy.
  - `BACKEND_SECRET` – Shared secret checked against `requiems-api-key`.
