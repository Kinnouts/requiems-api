# Edge Auth Gateway

Cloudflare Worker that handles API authentication, rate limiting, and credit
tracking for Requiem API.

## Setup

### 1. Install Wrangler

```bash
npm install -g wrangler
wrangler login
```

### 2. Create KV Namespace

```bash
wrangler kv:namespace create KV
# Copy the ID to wrangler.toml
```

### 3. Create D1 Database

```bash
wrangler d1 create requiem-usage
# Copy the ID to wrangler.toml

# Apply schema
wrangler d1 execute requiem-usage --file=schema.sql
```

### 4. Set Secrets

```bash
# Backend URL (keep this private!)
wrangler secret put BACKEND_URL
# Enter your Go backend URL (e.g., https://requiem-backend.fly.dev)

# Backend secret (min 32 chars - generate a strong random string)
wrangler secret put BACKEND_SECRET
# Enter a strong secret, e.g.: openssl rand -base64 32
```

### 5. Seed KV with API Keys

```bash
# Add a test API key
wrangler kv:key put --binding=KV "key:rq_test_123" '{"userId":"user-1","plan":"free","createdAt":"2025-01-01T00:00:00Z"}'
```

## Local Development

```bash
wrangler dev
```

Test locally:

```bash
curl -H "x-api-key: rq_test_123" http://localhost:8787/v1/text/advice
```

## Deploy

```bash
wrangler deploy
```

## Architecture

```
User Request
    │
    ▼
┌─────────────────────────────────────────┐
│  Worker (api.requiems-api.xyz)          │
│                                         │
│  1. Validate API key (KV lookup)        │
│  2. Check rate limit (KV counter)       │
│  3. Check credits (D1 query)            │
│  4. Forward to backend                  │
│  5. Record usage (D1 insert)            │
│  6. Add headers, return response        │
│                                         │
└─────────────────────────────────────────┘
    │
    ▼
  Backend (internal URL)
```

## Data Storage

### KV (Key-Value)

- **Fast reads** (~10ms globally)
- API keys: `key:{api_key}` → `{ userId, plan, createdAt }`
- Rate limits: `ratelimit:{api_key}:{minute}` → count

### D1 (SQLite)

- **Queryable** (can SUM, filter by date)
- Usage logs: `credit_usage` table

## Response Headers

Every response includes:

```
X-Credits-Used: 1
X-Credits-Remaining: 49
X-Credits-Reset: 2025-12-16T00:00:00Z
X-Plan: free
X-RateLimit-Limit: 30
X-RateLimit-Remaining: 29
```

## Endpoint Costs

Update `ENDPOINT_COSTS` in `src/config.ts` when adding new routes:

```typescript
export const ENDPOINT_COSTS: Record<string, number> = {
  "GET /v1/text/advice": 1,
  "GET /v1/text/quotes/random": 1,
  // Add new endpoints here
};
```

## Environment Variables

| Variable         | Source        | Description                                |
| ---------------- | ------------- | ------------------------------------------ |
| `BACKEND_URL`    | Secret        | Internal backend URL (not public)          |
| `BACKEND_SECRET` | Secret        | Auth header sent to backend (min 32 chars) |
| `ENVIRONMENT`    | wrangler.toml | `development`, `staging`, or `production`  |

## Project Structure

```
src/
├── index.ts        # Main fetch handler
├── env.ts          # t3-env validation (process.env)
├── types.ts        # TypeScript types
├── config.ts       # PLANS + ENDPOINT_COSTS
├── rate-limit.ts   # KV-based rate limiting
├── credits.ts      # D1 usage tracking
└── http.ts         # Response helpers
```
