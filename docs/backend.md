## Backend (Go) Overview

### Goals

- Single monolithic service that can support 500+ endpoints.
- Simple, feature-oriented structure that is easy for multiple developers to
  extend.

### Project layout (backend-relevant)

- `apps/api/main.go`
  - Loads config (`platform/config`).
  - Builds the application (`app`).
  - Starts the HTTP server.
- `app/`
  - `app.go` – Runs migrations with retry, creates the `chi` router, mounts
    feature routers under `/v1`.
  - `routes_v1.go` – Mounts all domain routers.
  - `healthz.go` – Health check endpoint at `/healthz`.
- `platform/config/`
  - `Config` struct + `Load()` reading `PORT`, `DATABASE_URL`, and `REDIS_URL`.
- `platform/db/`
  - `db.go` – `Connect` creates a `pgxpool.Pool`.
  - `migrate.go` – Runs SQL migrations from `migrations/` using
    `golang-migrate`.
- `platform/reqredis/`
  - `redis.go` – `Connect` creates a `*redis.Client` (go-redis/v9).
- `platform/httpx/`
  - `JSON` and `Error` helpers for writing JSON responses.
- `services/text/`
  - `router.go` – Mounts advice, quotes, words sub-routers.
- `services/text/<feature>/` (e.g. `advice`, `quotes`, `words`)
  - `service.go` – encapsulates DB access and business logic.
  - `transport_http.go` – exposes `RegisterRoutes(r chi.Router, svc *Service)`.
- `services/misc/`
  - `router.go` – Wires up the counter service and mounts its routes under
    `/v1/misc`.
- `services/misc/counter/`
  - `models.go` – `Counter` response struct and namespace validation.
  - `repository.go` – `Repository` interface + PostgreSQL `Upsert`/`Get`.
  - `service.go` – `Service` interface; Redis-primary reads/writes with
    PostgreSQL fallback.
  - `sync_worker.go` – Background goroutine; `SCAN counter:*` → pipeline GET →
    batch PostgreSQL upsert every 60 seconds.
  - `handler.go` – HTTP handlers for `POST` and `GET` on `/counter/{namespace}`.

### Adding a new API feature (pattern)

1. **Migrations**
   - Add a pair of migrations in `migrations/`:
     - `000X_feature_name.up.sql`
     - `000X_feature_name.down.sql`

2. **Service package**
   - Create `services/<domain>/<feature>/service.go`:
     - Define types and methods, e.g. `Random`, `GetByID`, etc.

3. **HTTP transport**
   - Create `services/<domain>/<feature>/transport_http.go` with:
     - `func RegisterRoutes(r chi.Router, svc *Service)`
     - Register your endpoints using `r.Get`, `r.Post`, etc.

4. **Wire in `services/<domain>/router.go`**
   - Instantiate the service and call `RegisterRoutes`:
     - `featureSvc := feature.NewService(pool)`
     - `feature.RegisterRoutes(r, featureSvc)`

   Or for a new top-level domain, wire in `app/routes_v1.go`:
   - Create a new router: `domainRouter := chi.NewRouter()`
   - Register routes and mount: `r.Mount("/v1/domain", domainRouter)`

   For features that also need Redis (e.g. counters), pass `*redis.Client` into
   the service and start any background workers via `go worker.Start(ctx, ...)`.

This keeps each feature self-contained and scales cleanly to hundreds of
endpoints without a single huge routes file.
