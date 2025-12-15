## Architecture

### High-Level

- **Edge Auth Gateway (Cloudflare Worker)**:
  - Validates `x-api-key`.
  - Enforces rate limits and plans (future).
  - Forwards trusted requests to the backend API.
- **Backend API (Go)**:
  - Single monolithic service.
  - Exposes all public HTTP endpoints (over 500+ planned).
  - Uses PostgreSQL as the primary store and Redis for cache/queues.
- **Workers (Go, planned)**:
  - Execute heavy/async tasks.
  - Process files and long-running jobs.

### Code Layout

- `apps/api` – main Go API entrypoint.
- `apps/edge-auth` – Cloudflare Worker for edge authentication.
- `internal/` – application code (domain modules, infra helpers, etc.).
- `infra/` – Docker, migrations, and deployment-related files.
- `docs/` – this documentation.

### Backend structure (`internal/`)

- `internal/app`:
  - Builds the HTTP router (using `chi`).
  - Runs DB migrations with retry.
  - Wires feature modules (advice, quotes, words, etc.).
- `internal/config`:
  - Loads configuration from env vars (`PORT`, `DATABASE_URL`, etc.).
- `internal/db`:
  - PostgreSQL connection pool (pgx).
  - `golang-migrate` integration (migrations from `infra/migrations`).
- `internal/httpx`:
  - Small HTTP helpers for JSON responses and errors.
- `internal/<feature>` (e.g. `advice`, `quotes`, `words`):
  - `service.go` – business logic and DB queries.
  - `transport_http.go` – HTTP routes and handlers (using `chi.Router`).

This feature-first layout scales well to hundreds of endpoints: each new API
gets its own package under `internal/` with a `Service` and a `RegisterRoutes`
function.
