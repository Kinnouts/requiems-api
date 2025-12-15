## Infrastructure & Deployment

For full details, see `infra/readme.md`. Key points:

- Local dev:
  - `cd infra/docker && docker compose up --build` – runs API + Postgres +
    Redis.
  - Hybrid dev: run only `db` and `redis` in Docker, run Go API locally with hot
    reload.
- Production:
  - Single VPS running Docker Compose.
  - Caddy terminates HTTPS and proxies to the API container.
  - Cloudflare Worker sits in front and handles auth and rate limits.
