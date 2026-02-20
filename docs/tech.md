# Technology Stack

This document lists all major technologies, their versions, and where to update
them when upgrading.

## Core Languages & Runtimes

### Go 1.24

**What:** Backend API service language **Current Version:** 1.24 **Files to
Update:**

- `.github/workflows/ci.yml` - CI test runner (line 65: `go-version: "1.24"`)
- `apps/api/go.mod` - Go module definition (line 3: `go 1.24`)
- `infra/docker/docker-compose.dev.yml` - Dev container (line 31:
  `image: golang:1.24-alpine`)
- `infra/docker/api.Dockerfile` - Production build (line 1: `FROM golang:1.24-alpine AS build`)
- `apps/api/.golangci.yml` - Linter configuration (line 6: `go: "1.24"`)

**Verification:**

```bash
cd apps/api && go version  # Should show go1.24.x
```

### Ruby 3.4.8

**What:** Rails dashboard language **Current Version:** 3.4.8 **Files to
Update:**

- `.github/workflows/ci.yml` - CI test runner (line 153:
  `ruby-version: "3.4.8"`)
- `apps/dashboard/Dockerfile` - Production build (line 11:
  `ARG RUBY_VERSION=3.4.8`)
- `infra/docker/docker-compose.dev.yml` - Dev container (line 55, 106:
  `image: ruby:3.4-alpine`)
- `apps/dashboard/.ruby-version` - Ruby version manager pin (rbenv/rvm/asdf)

**Note:** Dev Docker uses `ruby:3.4-alpine` (will use latest 3.4.x), production
Dockerfile pins exact version.

**Verification:**

```bash
cd apps/dashboard && ruby -v  # Should show ruby 3.4.8
```

### TypeScript 5.3.3

**What:** Cloudflare Worker language **Current Version:** 5.3.3 **Files to
Update:**

- `apps/edge-auth/package.json` - TypeScript compiler (line 32:
  `"typescript": "^5.3.3"`)

**Note:** Caret `^` allows patch updates (5.3.x). For major updates, also check
`@cloudflare/workers-types` compatibility.

**Verification:**

```bash
cd apps/edge-auth && bun run typecheck  # Should pass
```

### Bun (Latest)

**What:** JavaScript runtime for Cloudflare Worker development **Current
Version:** Latest (not pinned) **Files to Update:**

- `.github/workflows/ci.yml` - CI setup (line 261, 282, 316:
  `bun-version: latest`)
- `apps/edge-auth/package.json` - Lock file only (no version specified)

**Note:** Intentionally unpinned. Uses latest stable from oven-sh/setup-bun.

**Verification:**

```bash
bun --version
```

## Frameworks

### Rails 8.1.2

**What:** Web framework for dashboard **Current Version:** ~> 8.1.2 (allows
8.1.x patches) **Files to Update:**

- `apps/dashboard/Gemfile` - Rails version (line 4: `gem "rails", "~> 8.1.2"`)
- `apps/dashboard/Gemfile.lock` - Auto-updated by bundler

**Update Process:**

```bash
cd apps/dashboard
# Update Gemfile version
bundle update rails
bin/rails app:update  # Apply framework updates
```

## Infrastructure

### PostgreSQL 16

**What:** Primary database for both Go API and Rails dashboard **Current
Version:** 16-alpine **Files to Update:**

- `infra/docker/docker-compose.dev.yml` - Dev database (line 6:
  `image: postgres:16-alpine`)
- `.github/workflows/ci.yml` - CI test database (lines 46, 124:
  `image: postgres:16-alpine`)

**Note:** Shared database with separate migration systems (Go:
`infra/migrations/*.sql`, Rails: `apps/dashboard/db/migrate/*.rb`)

**Credentials (Dev/Test):**

- User: `requiem`
- Password: `requiem`
- Database: `requiem` (dev), `requiem_test` (CI)

### Redis 7

**What:** Background job queue for Rails and real-time counter storage for the
Go API **Current Version:** 7-alpine
**Files to Update:**

- `infra/docker/docker-compose.dev.yml` - Dev Redis (line 23:
  `image: redis:7-alpine`)
- `.github/workflows/ci.yml` - CI Redis (line 138: `image: redis:7-alpine`)

### Cloudflare Workers Runtime

**What:** Edge execution environment for authentication gateway **Version:**
Managed by Cloudflare (automatically updated) **Files to Update:**

- `apps/edge-auth/wrangler.toml` - Compatibility date (controls available APIs)
- `apps/edge-auth/package.json` - `@cloudflare/workers-types` for TypeScript
  types (line 28)

**Note:** Set `compatibility_date` in wrangler.toml to opt into runtime updates.
See:
https://developers.cloudflare.com/workers/configuration/compatibility-dates/

## Development Tools

### Air 1.52.0

**What:** Hot reload for Go development **Current Version:** v1.52.0 **Files to
Update:**

- `infra/docker/docker-compose.dev.yml` - Dev container install (line 37:
  `air@v1.52.0`)
- `apps/api/.air.toml` - Configuration file

**Note:** Installed dynamically in dev container, not in production.

### golangci-lint 2.10.1

**What:** Go meta-linter **Current Version:** v2.10.1 (built with Go 1.24) **Files to Update:**

- `.github/workflows/ci.yml` - CI linter (line 110: `version: v2.10.1`)
- `apps/api/.golangci.yml` - Linter v2 configuration (use `golangci-lint migrate` to upgrade from v1)

**Verification:**

```bash
cd apps/api && golangci-lint --version
```

### Vitest 1.2.0

**What:** Test framework for Cloudflare Worker **Current Version:** ^1.2.0
**Files to Update:**

- `apps/edge-auth/package.json` - Test runner (line 33: `"vitest": "^1.2.0"`)
- `apps/edge-auth/vitest.config.ts` - Configuration

**Verification:**

```bash
cd apps/edge-auth && bun test
```

### Biome 1.5.3

**What:** Linter and formatter for Cloudflare Worker (replaces ESLint +
Prettier) **Current Version:** ^1.5.3 **Files to Update:**

- `apps/edge-auth/package.json` - Tool version (line 27:
  `"@biomejs/biome": "^1.5.3"`)
- `apps/edge-auth/biome.json` - Configuration

**Verification:**

```bash
cd apps/edge-auth && bun run lint && bun run format:check
```

### RuboCop (via rails-omakase)

**What:** Ruby linter for Rails dashboard **Current Version:** Managed by
rubocop-rails-omakase gem **Files to Update:**

- `apps/dashboard/Gemfile` - Style guide gem (line 81:
  `gem "rubocop-rails-omakase"`)
- `apps/dashboard/.rubocop.yml` - Configuration file

**Verification:**

```bash
cd apps/dashboard && bundle exec rubocop
```

## Update Checklist

When updating a core technology:

1. **Check dependencies:** Run tests locally after updating
2. **Update CI:** Ensure `.github/workflows/ci.yml` matches new version
3. **Update Docker:** Update both dev (`docker-compose.dev.yml`) and prod
   (`Dockerfile`) if applicable
4. **Run tests:** All three apps must pass CI before merging
5. **Update this file:** Keep `docs/tech.md` current with new versions

## Version Pinning Philosophy

- **Go/Ruby/TypeScript:** Pin major + minor (e.g., 1.24, 3.4.8, 5.3.3) for
  stability
- **Rails/Gems:** Use `~>` for patch updates (e.g., `~> 8.1.2` allows 8.1.x)
- **Node packages:** Use `^` for minor updates (e.g., `^1.2.0` allows 1.x.x)
- **Bun:** Use `latest` (fast-moving tool, not in production)
- **PostgreSQL/Redis:** Pin major version only (e.g., `16-alpine`, `7-alpine`)

## Testing Matrix (CI)

Current CI runs:

- **Go 1.24** + PostgreSQL 16 + Redis 7
- **Ruby 3.4.8** + Rails 8.1 + PostgreSQL 16 + Redis 7
- **Bun latest** + TypeScript 5.3.3 + Vitest 1.2.0

See `.github/workflows/ci.yml` for full test configuration.

## Quick Reference

```bash
# Check all versions locally
go version                    # Go 1.24.x
ruby -v                       # Ruby 3.4.8
cd apps/dashboard && bin/rails -v  # Rails 8.1.2.x
cd apps/edge-auth && bun --version # Latest
cd apps/edge-auth && bun run typecheck # TypeScript 5.3.3

# Run all tests
cd apps/api && go test ./...
cd apps/dashboard && bin/rails test
cd apps/edge-auth && bun test
```

## External Service Versions

- **Cloudflare Workers:** Runtime version managed by Cloudflare (set
  `compatibility_date` in `wrangler.toml`)
- **Cloudflare KV:** Key-value store, no version (API versioned)
- **Cloudflare D1:** SQLite edge database, no version (API versioned)

## Documentation

- Go: https://go.dev/doc/
- Ruby: https://www.ruby-lang.org/en/documentation/
- Rails: https://guides.rubyonrails.org/
- TypeScript: https://www.typescriptlang.org/docs/
- Cloudflare Workers: https://developers.cloudflare.com/workers/
