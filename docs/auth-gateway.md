# The Auth Gateway (Cloudflare Worker)

## Overview

The Auth Gateway is a Cloudflare Worker that sits at the edge of the Requiem API
architecture, handling authentication, rate limiting, and credit tracking before
requests reach the backend.

**Location:** [apps/edge-auth/](../apps/edge-auth/) **Technology:** TypeScript
on Cloudflare Workers **Runtime:** V8 isolates running globally across
Cloudflare's network

## Architecture

The Worker acts as the first line of defense and management for all API
requests:

```
Client Request
    ↓
Worker (validates requiems-api-key header)
    ↓
KV Lookup (API key → user data)
    ↓
Rate Limit Check (KV counter)
    ↓
Credit Check (D1 query)
    ↓
Forward to Backend (with X-Backend-Secret)
    ↓
Record Usage (D1 insert)
    ↓
Add Usage Headers & Return Response
```

## Responsibilities

### 1. API Key Validation

- Extracts `requiems-api-key` header from requests
- Looks up key data in Cloudflare KV (~10ms global latency)
- Returns 401 if key is invalid or missing

### 2. Rate Limiting

- Enforces per-minute rate limits based on plan tier
- Stores counters in KV with 60-second TTL
- Returns 429 when limits exceeded

### 3. Credit Tracking

- Checks credit usage against plan limits
- Queries Cloudflare D1 (SQLite at the edge)
- Tracks per-endpoint credit costs
- Supports daily reset (free tier) and monthly reset (paid tiers)

### 4. Request Forwarding

- Adds `X-Backend-Secret` header for backend authentication
- Strips sensitive headers (requiems-api-key, Cloudflare headers)
- Forwards to internal backend URL

### 5. Usage Recording

- Records credit usage in D1 database
- Tracks endpoint, timestamp, and credits used
- Enables usage analytics and billing

## Key Files

### Core Logic

- [src/index.ts](../apps/edge-auth/src/index.ts) - Main worker entry point
- [src/config.ts](../apps/edge-auth/src/config.ts) - Plan limits and endpoint
  costs
- [src/credits.ts](../apps/edge-auth/src/credits.ts) - Credit tracking logic
- [src/rate-limit.ts](../apps/edge-auth/src/rate-limit.ts) - Rate limiting logic
- [src/http.ts](../apps/edge-auth/src/http.ts) - HTTP helpers and header
  filtering

### Configuration

- [wrangler.toml](../apps/edge-auth/wrangler.toml) - Worker configuration
- [schema.sql](../apps/edge-auth/schema.sql) - D1 database schema

### Development

- [scripts/seed-kv.ts](../apps/edge-auth/scripts/seed-kv.ts) - Test data seeding

## Development

### Local Development

```bash
cd apps/edge-auth

# Install dependencies
bun install

# Run type checking
bun run typecheck

# Start local dev server
bun run dev
# Worker runs on http://localhost:8787
```

### Environment Variables

Set these as Wrangler secrets (never commit):

- `BACKEND_URL` - Internal backend endpoint (e.g.,
  `https://internal.requiems.xyz`)
- `BACKEND_SECRET` - Shared secret for backend authentication (min 32 chars)

```bash
wrangler secret put BACKEND_URL
wrangler secret put BACKEND_SECRET
```

### Cloudflare Resources

**KV Namespace (API keys & rate limits):**

```bash
wrangler kv:namespace create KV
# Add ID to wrangler.toml
```

**D1 Database (usage tracking):**

```bash
wrangler d1 create requiem-usage
wrangler d1 execute requiem-usage --file=schema.sql
# Add ID to wrangler.toml
```

### Testing Locally

```bash
# Seed test API keys
bun run kv:seed

# Test request
curl -H "requiems-api-key: rq_test_xxxxx" http://localhost:8787/v1/text/advice
```

## Response Headers

Every successful response includes usage information:

```
X-Credits-Used: 1
X-Credits-Remaining: 499
X-Credits-Reset: 2026-03-01T00:00:00.000Z
X-Plan: developer
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 4999
```

## Plan Tiers

| Plan         | Credit Limit | Period  | Rate Limit |
| ------------ | ------------ | ------- | ---------- |
| Free         | 50           | Daily   | 30/min     |
| Developer    | 500,000      | Monthly | 5,000/min  |
| Business     | 500,000      | Monthly | 5,000/min  |
| Professional | Unlimited    | Monthly | 50,000/min |

## Endpoint Costs

Most endpoints cost **1 credit** per request (default).

Expensive endpoints (require more resources):

- Dictionary operations: 2 credits
- Future AI/ML endpoints: 3-5 credits

See [src/config.ts](../apps/edge-auth/src/config.ts) for the full list.

## Deployment

```bash
# Deploy to production
bun run deploy

# Or deploy to specific environment
bun run deploy:prod
```

## Integration with Rails Dashboard

The Rails dashboard syncs API keys to Cloudflare KV automatically:

- New API key created → synced to KV via
  [CloudflareKvSyncService](../apps/dashboard/app/services/cloudflare/kv_sync_service.rb)
- API key updated → KV updated
- API key deleted → KV entry removed

This ensures the Worker always has up-to-date key information.

## Performance

- **KV Lookups:** ~10ms globally
- **D1 Queries:** ~15-30ms
- **Total Overhead:** ~50-100ms before reaching backend
- **Global Distribution:** Runs in 300+ Cloudflare data centers worldwide

## Related Documentation

- [Architecture Overview](./architecture.md)
- [Backend Documentation](./backend.md)
- [Rails App Documentation](./rails-app.md)
