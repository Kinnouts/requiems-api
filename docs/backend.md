## Backend (Go) Overview

### Goals

- Single monolithic service that can support 500+ endpoints.
- Simple, feature-oriented structure that is easy for multiple developers to
  extend.

### Project layout (backend-relevant)

- `apps/api/main.go`
  - Loads config (`internal/config`).
  - Builds the application (`internal/app`).
  - Starts the HTTP server.
- `internal/app`
  - Runs migrations with retry.
  - Creates the `chi` router.
  - Registers feature modules (advice, quotes, words, etc.).
- `internal/config`
  - `Config` struct + `Load()` reading `PORT` and `DATABASE_URL`.
- `internal/db`
  - `Connect` – creates a `pgxpool.Pool`.
  - `Migrate` – runs SQL migrations from `infra/migrations` using
    `golang-migrate`.
- `internal/httpx`
  - `JSON` and `Error` helpers for writing JSON responses.
- `internal/<feature>` (e.g. `advice`, `quotes`, `words`)
  - `service.go` – encapsulates DB access and business logic.
  - `transport_http.go` – exposes `RegisterRoutes(r chi.Router, svc *Service)`.

### Adding a new API feature (pattern)

1. **Migrations**
   - Add a pair of migrations in `infra/migrations`:
     - `000X_feature_name.up.sql`
     - `000X_feature_name.down.sql`

2. **Service package**
   - Create `internal/<feature>/service.go`:
     - Define types and methods, e.g. `Random`, `GetByID`, etc.

3. **HTTP transport**
   - Create `internal/<feature>/transport_http.go` with:
     - `func RegisterRoutes(r chi.Router, svc *Service)`
     - Register your endpoints using `r.Get`, `r.Post`, etc.

4. **Wire in `internal/app/app.go`**
   - Instantiate the service and call `RegisterRoutes`:
     - `featureSvc := feature.NewService(pool)`
     - `feature.RegisterRoutes(r, featureSvc)`

This keeps each feature self-contained and scales cleanly to hundreds of
endpoints without a single huge routes file.
