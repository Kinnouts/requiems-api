# GlitchTip Error Monitoring — Design Plan

**Date:** 2026-04-14

## Overview

This document describes the integration of [GlitchTip](https://glitchtip.com)
(self-hosted, Sentry-compatible error monitoring) across all four services in
the Requiems API monorepo.

---

## Why This Exists

The platform previously had no centralised error tracking. Unhandled exceptions
in the Go backend, Rails dashboard, and Cloudflare Workers were only visible via
container logs — requiring SSH access to diagnose production issues. GlitchTip
provides a single UI for all four services with stack traces, request context,
and event history.

GlitchTip was chosen over Sentry SaaS because:

- Self-hosted at `issues.bobadilla.tech` — no third-party data exposure
- Fully Sentry-SDK-compatible (same SDKs, same client APIs)
- Zero per-event cost

---

## Architecture

Each service reports to a dedicated GlitchTip project with its own DSN. DSNs are
**ingest-only** (write events, read nothing) and safe to commit to source
control.

| Service                    | GlitchTip Project | SDK                              |
| -------------------------- | ----------------- | -------------------------------- |
| Go API                     | `#1`              | `github.com/getsentry/sentry-go` |
| Rails Dashboard            | `#2`              | `sentry-ruby` + `sentry-rails`   |
| Auth Gateway (CF Worker)   | `#3`              | `@sentry/cloudflare`             |
| API Management (CF Worker) | `#4`              | `@sentry/cloudflare`             |

---

## Implementation Per Service

### Go API (`apps/api/`)

- **Dependency:** `github.com/getsentry/sentry-go`
- **Init:** `main.go` — initializes Sentry after config load, gated by
  `cfg.Environment != "development"`. DSN is baked in as the default value for
  the `SENTRY_DSN` config field (overridable via env var).
- **Capture:** `platform/httpx/handler.go` — both `Handle` and `HandleBatch`
  call `sentry.CaptureException(err)` before returning a 500. `*AppError`
  responses (expected domain errors) are intentionally not captured.
- **Environment flag:** `ENVIRONMENT` env var (default `"development"`). Set to
  `"production"` in `infra/docker/docker-compose.yml` for the `api` service.

### Rails Dashboard (`apps/dashboard/`)

- **Gems:** `sentry-ruby`, `sentry-rails`
- **Init:** `config/initializers/sentry.rb` — DSN baked in via `ENV.fetch`
  fallback. Reporting is gated to
  `enabled_environments = %w[production staging]` which uses `RAILS_ENV` —
  already set to `"production"` in the production Docker image.
- **Capture:** automatic via `sentry-rails` — all unhandled controller
  exceptions are reported without any manual instrumentation.

### Auth Gateway & API Management (`apps/workers/`)

Both Cloudflare Workers follow the same pattern:

- **Package:** `@sentry/cloudflare` (not `@sentry/node` — CF Workers run on the
  V8 isolate runtime, not Node.js)
- **Init:** `wrapRequestHandler` is called on every `fetch` invocation. It
  handles SDK initialization internally (only truly initializes once per
  isolate).
- **Capture:** `app.onError` is overridden per-worker to call
  `captureException(err)` before delegating to the shared error handler. This is
  necessary because Hono catches thrown errors via `onError` — they never
  propagate out to `wrapRequestHandler`'s automatic capture.
- **DSN in wrangler.toml:** empty string in `[vars]` (dev/local), real DSN in
  `[env.production.vars]`. No secrets mechanism needed since DSNs are
  ingest-only.

---

## Environment Strategy

Reporting is **off in development** across all services without any local `.env`
configuration:

| Service        | Dev off because                             | Production on because               |
| -------------- | ------------------------------------------- | ----------------------------------- |
| Go API         | `ENVIRONMENT` defaults to `"development"`   | `ENVIRONMENT=production` in compose |
| Rails          | `enabled_environments` gates on `RAILS_ENV` | `RAILS_ENV=production` already set  |
| Auth Gateway   | `SENTRY_DSN = ""` in `[vars]`               | DSN set in `[env.production.vars]`  |
| API Management | Same                                        | Same                                |

---

## Files Changed

| File                                           | Change                                                 |
| ---------------------------------------------- | ------------------------------------------------------ |
| `apps/api/go.mod`                              | Added `github.com/getsentry/sentry-go`                 |
| `apps/api/platform/config/config.go`           | Added `Environment`, `SentryDSN` fields                |
| `apps/api/main.go`                             | Sentry init gated by environment                       |
| `apps/api/platform/httpx/handler.go`           | `CaptureException` on 500 errors                       |
| `apps/dashboard/Gemfile`                       | Added `sentry-ruby`, `sentry-rails`                    |
| `apps/dashboard/Gemfile.lock`                  | Updated with resolved gem versions                     |
| `apps/dashboard/config/initializers/sentry.rb` | New — Sentry init for Rails                            |
| `apps/workers/auth-gateway/package.json`       | Added `@sentry/cloudflare`                             |
| `apps/workers/auth-gateway/src/env.ts`         | Added optional `SENTRY_DSN` binding                    |
| `apps/workers/auth-gateway/src/index.ts`       | `wrapRequestHandler` + `captureException` in `onError` |
| `apps/workers/auth-gateway/wrangler.toml`      | `SENTRY_DSN` in dev and production vars                |
| `apps/workers/api-management/package.json`     | Added `@sentry/cloudflare`                             |
| `apps/workers/api-management/src/env.ts`       | Added optional `SENTRY_DSN` binding                    |
| `apps/workers/api-management/src/index.ts`     | `wrapRequestHandler` + `captureException` in `onError` |
| `apps/workers/api-management/wrangler.toml`    | `SENTRY_DSN` in dev and production vars                |
| `infra/docker/docker-compose.yml`              | Added `ENVIRONMENT=production` to `api` service        |
