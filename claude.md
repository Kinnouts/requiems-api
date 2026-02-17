# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
code in this repository.

## Project Overview

Requiems API is a production-ready API platform providing unified access to
multiple enterprise-grade APIs (email validation, text utilities, etc.). Built
as a multi-language monorepo with:

- **Go 1.23 API** (apps/api) - Internal business logic backend
- **Rails 8 Dashboard** (apps/dashboard) - Public web UI, user management, admin
  panel
- **Cloudflare Worker Gateway** (apps/edge-auth) - Global edge auth, rate
  limiting, credit tracking

## Development Commands

### Full Stack Development (Recommended)

Start all services with hot reload:

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml up
```

Services:

- Go API: <http://localhost:8080/healthz> (Air hot reload)
- Rails Dashboard: <http://localhost:3000> (Rails hot reload)
- PostgreSQL: localhost:5432 (user: `requiem`, password: `requiem`)
- Redis: localhost:6379

Rebuild after dependency changes:

```bash
docker compose -f docker-compose.dev.yml up --build
```

### Go API (apps/api)

Run locally (requires PostgreSQL):

```bash
cd apps/api
go run main.go
```

Run tests:

```bash
cd apps/api
go test ./...
```

Run specific test:

```bash
cd apps/api
go test ./internal/email/disposable -v -run TestCheckDisposable
```

Build:

```bash
cd apps/api
go build -o bin/api main.go
```

Hot reload (via Air):

```bash
cd apps/api
air
```

### Rails Dashboard (apps/dashboard)

Rails console:

```bash
docker compose -f docker-compose.dev.yml exec dashboard rails console
```

Or locally:

```bash
cd apps/dashboard
bin/rails console
```

Run tests:

```bash
cd apps/dashboard
bin/rails test
```

Run specific test:

```bash
cd apps/dashboard
bin/rails test test/models/user_test.rb
```

Migrations:

```bash
cd apps/dashboard
bin/rails db:migrate
bin/rails db:rollback
```

### Cloudflare Worker (apps/edge-auth)

Dev mode:

```bash
cd apps/edge-auth
npm run dev
```

Deploy to staging:

```bash
cd apps/edge-auth
npm run deploy
```

Deploy to production:

```bash
cd apps/edge-auth
npm run deploy:prod
```

Type check:

```bash
cd apps/edge-auth
npm run typecheck
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

## Continuous Integration

The monorepo uses a unified GitHub Actions workflow at
`.github/workflows/ci.yml` that automatically tests and validates code across
all three applications.

### CI Philosophy

- **Tests are blocking** - PRs cannot merge if tests fail
- **Security scans are blocking** - Critical for production safety
- **Linting is advisory** - Shows warnings but doesn't block PRs
- **Path-based execution** - Only runs jobs for changed apps (efficient)

### What Gets Tested

**Go Backend (`apps/api`):**

- ✅ **Tests** (blocking) - `go test -race -coverprofile=coverage.out ./...`
- ⚠️ **golangci-lint** (advisory) - 21 linters enabled (errcheck, gosimple,
  govet, staticcheck, etc.)

**Rails Dashboard (`apps/dashboard`):**

- ✅ **Tests** (blocking) - Minitest suite
- ✅ **Security Scans** (blocking) - Brakeman, bundler-audit, importmap audit
- ⚠️ **RuboCop** (advisory) - rails-omakase style guide

**Cloudflare Worker (`apps/edge-auth`):**

- ✅ **TypeScript Check** (blocking) - `tsc --noEmit`
- ✅ **Tests** (blocking) - Vitest with 29% coverage (71 tests)
- ⚠️ **Biome Lint** (advisory) - Modern fast linter + formatter

### Local Testing Before Push

Run these commands locally to catch issues before CI:

```bash
# Go Backend
cd apps/api
go test ./...                    # Tests (must pass)
golangci-lint run                # Linting (advisory - requires golangci-lint installed)

# Rails Dashboard
cd apps/dashboard
bin/rails test                   # Tests (must pass)
bundle exec brakeman --no-pager  # Security scan (must pass)
bundle exec bundler-audit        # Dependency audit (must pass)
bin/importmap audit              # JS dependency audit (must pass)
bundle exec rubocop              # Linting (advisory)

# Cloudflare Worker
cd apps/edge-auth
bunx vitest run                  # Tests (must pass - 71 tests)
bun run typecheck                # TypeScript (must pass)
bun run lint                     # Linting (advisory)
bun run format:check             # Formatting (advisory)
```

### CI Workflow Behavior

The CI automatically detects which apps changed and runs only relevant tests:

- **Change Go files** → Only Go tests and lint run
- **Change Rails files** → Only Rails tests, security, and lint run
- **Change Worker files** → Only Worker tests and lint run
- **Change multiple apps** → All relevant jobs run in parallel

This makes CI fast - typically under 5 minutes for single-app changes.

### Branch Protection

Configure GitHub branch protection to require the **"CI Success"** job for
merging:

1. Go to Settings → Branches → Branch protection rules
2. Select `main` branch
3. Enable "Require status checks to pass before merging"
4. Select "CI Success" as required check
5. Enable "Require branches to be up to date before merging"

### CI Configuration Files

**Tool Configurations:**

- `apps/api/.golangci.yml` - Go linter config (21 linters)
- `apps/edge-auth/vitest.config.ts` - Test framework config
- `apps/edge-auth/biome.json` - Linter + formatter config

**Workflow:**

- `.github/workflows/ci.yml` - Main CI orchestrator (all apps)

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

**When to extract partials:**

- Page exceeds ~100-150 lines
- Content is reused across multiple pages
- Sections have distinct responsibilities
- Need to test sections independently

**Rendering partials:**

```erb
<!-- Page-specific partials -->
<%= render 'partials/contact/info_sections' %>
<%= render 'partials/contact/contact_form' %>

<!-- Shared partials -->
<%= render 'partials/shared/footer' %>
```

**Why this structure:**

- Keeps controller view directories clean (one file per page)
- Clear separation between views and components
- Easy to find partials organized by page name
- Scales well as the application grows
- Similar to React's component-per-folder structure

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
cd apps/dashboard
bin/rails generate migration AddFieldToTable
bin/rails db:migrate
```

### Cloudflare Worker Development

1. Local development uses miniflare (simulates KV, D1)
2. Secrets must be set via `wrangler secret put`:
   - `BACKEND_URL`
   - `BACKEND_SECRET`
3. KV/D1 bindings configured in `wrangler.toml`
4. Seed KV with: `bun run scripts/seed-kv.ts`

### Hot Reload

- **Go**: Air watches `.go` files, rebuilds (~2-3s)
- **Rails**: Native Rails reloader, instant for most changes
- **Cloudflare Worker**: Wrangler dev mode with instant reload

### Running Single Tests

Go:

```bash
cd apps/api
go test ./internal/text/advice -v -run TestGetAdvice
```

Rails:

```bash
cd apps/dashboard
bin/rails test test/models/api_key_test.rb:15  # Line number
```

## Common Development Tasks

View service logs:

```bash
docker compose -f docker-compose.dev.yml logs -f api
docker compose -f docker-compose.dev.yml logs -f dashboard
```

Reset database:

```bash
docker compose -f docker-compose.dev.yml down -v
docker compose -f docker-compose.dev.yml up
```

Connect to database:

- Host: localhost
- Port: 5432
- Database: requiem
- User: requiem
- Password: requiem
