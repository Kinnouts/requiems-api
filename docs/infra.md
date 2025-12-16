## Infrastructure & Deployment

For full details, see `infra/readme.md`. Key points:

- Local dev:
  - `cd infra/docker && docker compose up --build` – runs API + Postgres.
  - API available at `http://localhost:6969`
  - Postgres available at `localhost:5432`
  - Hybrid dev: run only `db` in Docker, run Go API locally with hot reload.
- Docker build:
  - Uses multi-stage build for small images (~15MB)
  - Layer caching: dependencies cached separately from source code
  - Runs as non-root user (`appuser`) for security
  - `.dockerignore` at repo root excludes unnecessary files
- Production:
  - Single VPS running Docker Compose.
  - Caddy terminates HTTPS and proxies to the API container.
  - Cloudflare Worker sits in front and handles auth and rate limits.

> **Note:** Redis is provisioned in docker-compose.yml for future use
> (queues/cache) but is not currently used by the Go backend.
