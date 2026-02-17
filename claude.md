# Project Overview

Requiems API is a production-ready API platform providing unified access to
multiple enterprise-grade APIs (email validation, text utilities, etc.). Built
as a multi-language monorepo with:

- **Go API** (apps/api) - Internal business logic backend
- **Rails Dashboard** (apps/dashboard) - Public web UI, user management, admin
  panel
- **Cloudflare Worker Gateway** (apps/edge-auth) - Global edge auth, rate
  limiting, credit tracking

## Development Commands

### Full Stack Development

Start all services with hot reload:

```bash
docker compose -f docker-compose.dev.yml up
```

Services:

- Go API: <http://localhost:8080/healthz> (Air hot reload)
- Rails Dashboard: <http://localhost:3000> (Rails hot reload)
- PostgreSQL: localhost:5432 (user: `requiem`, password: `requiem`)
- Redis: localhost:6379

### Go API (apps/api)

**Note:** Commands assume containers are running. Container: `requiem-dev-api-1`

Run tests:

```bash
docker exec requiem-dev-api-1 go test ./...
```

Run tests with coverage:

```bash
docker exec requiem-dev-api-1 go test -race -coverprofile=coverage.out ./...
```

Run specific test:

```bash
docker exec requiem-dev-api-1 go test ./internal/text/advice -v -run TestGetAdvice
```

### Rails Dashboard (apps/dashboard)

**Note:** Commands assume containers are running. Container: `requiem-dev-dashboard-1`

Rails console:

```bash
docker exec -it requiem-dev-dashboard-1 rails console
```

Run tests:

```bash
docker exec requiem-dev-dashboard-1 bin/rails test
```

Run specific test:

```bash
docker exec requiem-dev-dashboard-1 bin/rails test test/models/user_test.rb
```

Run security scans:

```bash
docker exec requiem-dev-dashboard-1 bundle exec brakeman --no-pager
docker exec requiem-dev-dashboard-1 bundle exec bundler-audit
docker exec requiem-dev-dashboard-1 bin/importmap audit
```

Run linting:

```bash
docker exec requiem-dev-dashboard-1 bundle exec rubocop
```

Migrations:

```bash
docker exec requiem-dev-dashboard-1 bin/rails db:migrate
docker exec requiem-dev-dashboard-1 bin/rails db:rollback
```

Generate migration:

```bash
docker exec requiem-dev-dashboard-1 bin/rails generate migration AddFieldToTable
```

### Cloudflare Worker (apps/edge-auth)

Dev mode:

```bash
cd apps/edge-auth
bun dev
```

Type check:

```bash
cd apps/edge-auth
bun run typecheck
```

Run tests:

```bash
cd apps/edge-auth
bunx vitest run              # Run all tests
bunx vitest run --coverage   # With coverage report
bun run test:watch           # Watch mode
```

Run linting and formatting:

```bash
cd apps/edge-auth
bun run lint                 # Lint code
bun run lint:fix             # Auto-fix lint issues
bun run format:check         # Check formatting
bun run format               # Auto-format code
```

### Local Testing Before Push

Run these commands locally to catch issues before CI.

**Note:** Containers must be running. Use `docker compose -f docker-compose.dev.yml up` in `infra/docker` first.

```bash

# Go Backend
docker exec requiem-dev-api-1 go test ./...                    # Tests (must pass)
docker exec requiem-dev-api-1 golangci-lint run                # Linting (advisory)

# Rails Dashboard
docker exec requiem-dev-dashboard-1 bin/rails test             # Tests (must pass)
docker exec requiem-dev-dashboard-1 bundle exec brakeman --no-pager  # Security (must pass)
docker exec requiem-dev-dashboard-1 bundle exec bundler-audit  # Dependency audit (must pass)
docker exec requiem-dev-dashboard-1 bin/importmap audit        # JS audit (must pass)
docker exec requiem-dev-dashboard-1 bundle exec rubocop        # Linting (advisory)
```

For Cloudflare Worker (run locally, not in Docker):

```bash
cd apps/edge-auth
bunx vitest run                  # Tests (must pass - 71 tests)
bun run typecheck                # TypeScript (must pass)
bun run lint                     # Linting (advisory)
bun run format:check             # Formatting (advisory)
```

## Architecture Overview

### Request Flow

```
User → Cloudflare Worker (edge-auth) → Go Backend (api) → PostgreSQL
       ↓                                  ↓
       KV (auth, rate limit)              Business logic
       D1 (usage tracking)
```

1. **Cloudflare Worker** (`apps/edge-auth/src/index.ts`):
   - Validates API key from Cloudflare KV
   - Checks rate limits (KV counters)
   - Checks credit usage (D1 SQLite queries)
   - Forwards to Go backend with `X-Backend-Secret` header
   - Records usage in D1
   - Returns response with usage headers

2. **Go Backend** (`apps/api/main.go`):
   - Receives requests from gateway only
   - Executes business logic
   - Queries PostgreSQL for data
   - Returns JSON responses
   - **No authentication** - trusts the gateway

3. **Rails Dashboard** (`apps/dashboard`):
   - User registration/login
   - API key management
   - Usage statistics
   - Admin panel
   - Syncs API keys to Cloudflare KV

### Code Organization

#### Go API (apps/api/internal/)

Domain-driven design with feature modules:

- `app/` - Application initialization, routing
- `email/` - Email-related endpoints (disposable checking, etc.)
  - `disposable/service.go` - Business logic
  - `disposable/transport_http.go` - HTTP handlers
  - `disposable/type.go` - Types
  - `router.go` - Routes for `/v1/email/*`
- `text/` - Text utility endpoints (advice, lorem, quotes, words)
  - Each subdomain follows same pattern: service, transport_http, type
  - `router.go` - Routes for `/v1/text/*`
- `platform/` - Shared infrastructure
  - `config/` - Environment configuration
  - `db/` - PostgreSQL connection, migrations
  - `httpx/` - HTTP utilities

**Pattern**: Each feature has `service.go` (business logic), `transport_http.go`
(HTTP handlers), `type.go` (data types), and parent `router.go` (routes).

#### Cloudflare Worker (apps/edge-auth/src/)

- `index.ts` - Main request handler (auth, rate limit, proxy)
- `requests.ts` - Usage tracking (D1 queries)
- `rate-limit.ts` - Rate limiting logic (KV)
- `http.ts` - HTTP utilities
- `config.ts` - Plans, pricing
- `types.ts` - TypeScript types

#### Rails Dashboard (apps/dashboard/)

Standard Rails 8 structure with:

- Turbo + Stimulus for interactivity
- Tailwind CSS for styling
- Sidekiq for background jobs
- Solid Queue/Cache for Rails 8 features

**View Organization Pattern**:

Rails views are organized to keep controller directories clean with
page-specific partials in a dedicated `partials/` directory:

```
app/views/
├── {controller}/
│   ├── {action}.html.erb      # Main views only (one file per page)
│   └── {another_action}.html.erb
└── partials/
    ├── {page_name}/
    │   ├── _section_name.html.erb
    │   └── _another_section.html.erb
    └── shared/                 # Truly shared across multiple pages
        └── _component.html.erb
```

Example structure:

```
app/views/
├── home/
│   ├── contact.html.erb       # Clean! Main views only
│   ├── about.html.erb
│   └── pricing.html.erb
└── partials/
    ├── contact/
    │   ├── _info_sections.html.erb
    │   ├── _additional_links.html.erb
    │   └── _contact_form.html.erb
    └── shared/
        └── _footer.html.erb
```

## Database Architecture

### PostgreSQL (Shared)

Single database with two migration systems:

- **Go migrations**: `infra/migrations/*.sql` (business data tables: advice,
  quotes, words)
- **Rails migrations**: `apps/dashboard/db/migrate/*.rb` (user tables: users,
  api_keys, subscriptions, usage_logs)

Separate migration tracking tables prevent conflicts.

### Cloudflare KV (Edge)

Key-value store for:

- API key lookup: `key:{api_key}` → `{userId, plan, billingCycleStart, ...}`
- Rate limiting: `rl:m:{key}:{minute}` → request count (auto-expires)

### Cloudflare D1 (Edge SQLite)

Usage tracking:

- `credit_usage` table: Records every API call with credits used
- SQL queries for billing period aggregations

## Key Dependencies

### Go

- `chi` - HTTP router
- `pgx` - PostgreSQL driver
- `golang-migrate` - Database migrations
- `bobadilla-tech/is-email-disposable` - Email validation
- `bobadilla-tech/lorelai` - Lorem ipsum generation

### Rails

- Rails 8.1
- Tailwind CSS
- Turbo/Stimulus (Hotwire)
- Solid Cache/Queue
- Sidekiq

### Cloudflare Worker

- Wrangler - Deployment tool
- `@cloudflare/workers-types` - TypeScript types
- `zod` - Schema validation
- `@t3-oss/env-core` - Environment variables

## Important Notes

### Adding New Go Endpoints

1. Create feature directory in `apps/api/internal/{domain}/{feature}/`
2. Add `service.go` (business logic)
3. Add `transport_http.go` (HTTP handler)
4. Add `type.go` (request/response types)
5. Register routes in parent `router.go`
6. Mount router in `apps/api/internal/app/app.go`

Example structure:

```
internal/
  text/
    advice/
      service.go
      transport_http.go
      type.go
    router.go  # Mounts all text/* routes
```

### Go Backend Security

The Go backend trusts the Cloudflare Worker gateway completely:

- No API key validation in Go
- No rate limiting in Go
- Expects `X-Backend-Secret` header from gateway
- Only processes business logic and database queries

### Database Migrations

**Go migrations** (business data):

```bash
# Add new migration in infra/migrations/
# Named: YYYYMMDDHHMMSS_description.up.sql
# Runs automatically on app startup via app/app.go:migrateWithRetry()
```

**Rails migrations** (user data):

```bash
docker exec requiem-dev-dashboard-1 bin/rails generate migration AddFieldToTable
docker exec requiem-dev-dashboard-1 bin/rails db:migrate
```

### Cloudflare Worker Development

1. Local development uses miniflare (simulates KV, D1)
2. Secrets must be set via `wrangler secret put`:
   - `BACKEND_URL`
   - `BACKEND_SECRET`
3. KV/D1 bindings configured in `wrangler.toml`
4. Seed KV with: `bun run scripts/seed-kv.ts`

## Common Development Tasks

**Note:** These commands manage the Docker containers. Run from `infra/docker` directory.

View service logs:

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml logs -f api
docker compose -f docker-compose.dev.yml logs -f dashboard
```

Reset database:

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml down -v
docker compose -f docker-compose.dev.yml up
```

Restart a specific service:

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml restart api
docker compose -f docker-compose.dev.yml restart dashboard
```

Connect to database:

- Host: localhost
- Port: 5432
- Database: requiem
- User: requiem
- Password: requiem
