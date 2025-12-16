## Infrastructure Guide

This document explains how to run Requiem API locally for development and how to
deploy it to a single VPS.

---

## 1. Local Development Setup

### 1.1 Requirements

- Docker and Docker Compose
- Go 1.22+ (optional, only needed if you want to run the API directly)
- Node/PNPM/Yarn (optional, for Cloudflare Worker tooling later)

### 1.2 Run everything with Docker (recommended)

From the project root:

```bash
cd infra/docker
docker compose up --build
```

This starts:

- `api` – Go backend on internal port `8080` (exposed as `localhost:6969`)
- `db` – PostgreSQL (`requiem` / `requiem` / `requiem`, exposed as
  `localhost:5432`)
- `redis` – Redis for future queues/cache

Once the stack is up:

- Health check: `http://localhost:6969/healthz`
- Advice endpoint: `http://localhost:6969/v1/advice`

> Note: Caddy is mainly for the VPS setup. For local development you can hit the
> API directly on `localhost:6969`.

### 1.3 Run the API directly (without Docker)

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

### 1.4 Hybrid dev workflow (Docker infra + local Go with hot reload)

For the best developer experience, run **Postgres and Redis in Docker**, and the
Go API locally with a watcher:

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

- `http://localhost:6969/healthz`
- `http://localhost:6969/v1/advice`

### 1.5 Cloudflare Worker (edge auth) – dev notes

The Worker lives in `apps/edge-auth/index.ts`. A typical dev setup will:

- Use `wrangler dev` to run the worker locally.
- Set environment variables:
  - `BACKEND_ORIGIN` – e.g. `http://localhost:8080`
  - `API_KEY_SECRET` – your shared API key used by clients in `x-api-key`.

> The exact `wrangler.toml` configuration can be added once you hook up your
> Cloudflare account.

---

## 2. VPS Deployment Guide

Target: single VPS (Docker + Docker Compose) running API, Postgres, Redis, and
Caddy. Cloudflare Worker sits in front and forwards authorized traffic.

### 2.1 Prepare the VPS

1. Create a new VPS (any provider).
2. SSH into the server.
3. Install Docker and Docker Compose (or Docker Compose plugin).
4. Clone the repository:

```bash
git clone <your-repo-url> requiems-api
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

### 2.4 DNS setup

- Create an `A` record for `api.yourdomain.com` pointing to your VPS public IP.
- Wait for DNS to propagate.
- With DNS + Caddy running, `https://api.yourdomain.com/healthz` should respond
  successfully.

### 2.5 Cloudflare Worker configuration (edge auth)

1. Deploy `apps/edge-auth/index.ts` as a Worker in your Cloudflare account.
2. Configure environment variables/secrets:
   - `BACKEND_ORIGIN` – `https://api.yourdomain.com`
   - `BACKEND_SECRET` – a strong secret, used by clients in `x-api-key`.
3. Route your public API endpoint to the Worker (e.g.
   `https://v1.yourdomain.com/*` → Worker).

Request flow in production:

1. Client → Worker (with `x-api-key`).
2. Worker validates `x-api-key`.
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
  - `BACKEND_SECRET` – Shared secret checked against `x-api-key`.
