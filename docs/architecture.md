# Architecture

## Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              INTERNET                                       │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     api.requiems-api.xyz (Cloudflare Worker)                │
│                                                                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                     │
│  │ Auth (KV)   │    │ Rate Limit  │    │ Credits     │                     │
│  │ API Keys    │    │ (KV)        │    │ (D1)        │                     │
│  └─────────────┘    └─────────────┘    └─────────────┘                     │
│                                                                             │
│  Validates x-api-key → Checks rate limit → Checks credits → Forward        │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ X-Backend-Secret header
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Internal Backend (Go + Chi)                             │
│                     (Private URL, not public DNS)                           │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         PostgreSQL                                   │   │
│  │                    (Business Data Only)                              │   │
│  │     advice, quotes, words, finance, places, etc.                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Data Stores Explained

We use **3 data stores**, each optimized for its specific use case:

### 1. Cloudflare KV (Key-Value)

**What:** Globally distributed key-value store with <10ms reads worldwide.

**Used for:**

- **API Key lookup** (`key:{api_key}` → user data, plan, etc.)
- **Rate limiting counters** (`rl:s:{key}:{second}` and `rl:m:{key}:{minute}`)

**Why KV:**

- Extremely fast reads (critical for every API request)
- TTL expiration (rate limit keys auto-delete)
- Simple key-value, no complex queries needed
- Data is small (just key → JSON blob)

**Example data:**

```json
// key:rq_live_abc123xyz
{
  "userId": "user_456",
  "plan": "starter",
  "createdAt": "2024-01-15T00:00:00Z",
  "billingCycleStart": "2024-12-01T00:00:00Z"
}
```

### 2. Cloudflare D1 (SQLite at Edge)

**What:** SQLite database running on Cloudflare's edge network.

**Used for:**

- **Credit usage tracking** (INSERT on every API call, SUM queries for totals)
- **Usage analytics** (which endpoints, when, how much)

**Why D1 (not KV):**

- Need SQL aggregations (`SUM`, `GROUP BY`, date ranges)
- Need to query "usage since X date" for billing periods
- Historical data for analytics/billing
- KV can't do `WHERE used_at >= '2024-12-01'`

**Schema:**

```sql
CREATE TABLE credit_usage (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  api_key TEXT NOT NULL,
  endpoint TEXT NOT NULL,
  credits_used INTEGER NOT NULL,
  used_at TEXT NOT NULL
);
```

**Example query:**

```sql
-- Get this month's usage for billing
SELECT SUM(credits_used) FROM credit_usage
WHERE api_key = ? AND used_at >= '2024-12-01T00:00:00Z'
```

### 3. PostgreSQL (Backend Database)

**What:** Traditional relational database running alongside the Go backend.

**Used for:**

- **All business data** (advice, quotes, words, etc.)
- **User accounts** (future)
- **Anything the API actually returns**

**Why PostgreSQL (not D1):**

- Complex queries, joins, full-text search
- Larger datasets (millions of rows)
- Backend runs on dedicated server, not edge
- D1 is edge-only, can't be accessed from Go backend

**Example tables:**

```sql
-- Business data in PostgreSQL
CREATE TABLE advice (id SERIAL, content TEXT, category TEXT, ...);
CREATE TABLE quotes (id SERIAL, text TEXT, author TEXT, ...);
CREATE TABLE words (id SERIAL, word TEXT, definition TEXT, ...);
```

## Why This Split?

| Store          | Location                | Latency  | Use Case                |
| -------------- | ----------------------- | -------- | ----------------------- |
| **KV**         | Edge (global)           | <10ms    | Auth, rate limits       |
| **D1**         | Edge (global)           | ~50ms    | Usage tracking, billing |
| **PostgreSQL** | Backend (single region) | ~100ms\* | Business data           |

\*Backend latency depends on user's distance to server region.

**The Gateway (Worker) handles:**

- Auth + rate limits (KV) - must be ultra-fast
- Credit tracking (D1) - needs SQL, still at edge

**The Backend (Go) handles:**

- Business logic - complex queries on PostgreSQL
- No auth overhead - trusts the gateway

## What We DON'T Use

- **R2 (Object Storage):** Not needed. We don't store files/images.
- **Redis:** PostgreSQL handles caching needs. Could add later for heavy
  caching.
- **Durable Objects:** Overkill. KV + D1 covers our needs.

## Request Flow

```
1. User calls: GET api.requiems-api.xyz/v1/text/advice
   └─ Header: x-api-key: rq_live_abc123

2. Worker receives request
   └─ KV.get("key:rq_live_abc123") → { plan: "starter", ... }

3. Check rate limit
   └─ KV.get("rl:s:rq_live_abc123:1702656000") → "5"
   └─ KV.get("rl:m:rq_live_abc123:28377600") → "150"
   └─ Under limit? Continue. Over? Return 429.

4. Check credits
   └─ D1: SELECT SUM(credits_used) WHERE api_key = ? AND used_at >= ?
   └─ Under limit? Continue. Over? Return 429 (free) or allow (paid).

5. Forward to backend
   └─ fetch(BACKEND_URL + "/v1/text/advice", { headers: { X-Backend-Secret } })

6. Backend processes
   └─ PostgreSQL: SELECT * FROM advice ORDER BY RANDOM() LIMIT 1
   └─ Returns JSON response

7. Record usage
   └─ D1: INSERT INTO credit_usage (api_key, endpoint, credits_used, ...)

8. Return response with headers
   └─ X-Credits-Used: 1
   └─ X-Credits-Remaining: 29849
   └─ X-RateLimit-Remaining: 149
```

## Code Layout

```
apps/
├── api/                    # Go backend entrypoint
│   └── main.go
└── edge-auth/              # Cloudflare Worker (gateway)
    ├── src/
    │   ├── index.ts        # Main handler
    │   ├── env.ts          # t3-env validation
    │   ├── types.ts        # TypeScript types
    │   ├── config.ts       # Plans, endpoint costs
    │   ├── rate-limit.ts   # KV rate limiting
    │   ├── credits.ts      # D1 usage tracking
    │   └── http.ts         # Response helpers
    ├── schema.sql          # D1 schema
    ├── wrangler.toml       # Worker config
    └── package.json

internal/
├── app/                    # Router setup, healthz
├── platform/               # Cross-cutting (config, db, httpx)
│   ├── config/
│   ├── db/
│   └── httpx/
└── text/                   # Domain: text APIs
    ├── router.go           # Mounts advice, quotes, words
    ├── advice/
    ├── quotes/
    └── words/

infra/
├── migrations/             # PostgreSQL migrations
├── docker/
└── caddy/
```

## Security Model

1. **API Keys** stored in KV, looked up on every request
2. **Rate limits** enforced at edge before backend is touched
3. **Backend URL** is secret (not in code, set via `wrangler secret`)
4. **Backend Secret** header required - even if URL leaks, can't call without it
5. **Backend trusts gateway** - no redundant auth checks
