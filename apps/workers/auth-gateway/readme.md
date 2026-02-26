# Auth Gateway

Public-facing Cloudflare Worker at `api.requiems.xyz`. Validates API keys,
enforces rate limits and monthly quotas, then proxies authenticated requests to
the Go backend.

## Request Flow

```
User → Auth Gateway
         ├── 1. Extract `requiems-api-key` header
         ├── 2. Lookup key in KV → ApiKeyData
         ├── 3. Check per-minute rate limit (KV counter)
         ├── 4. Check monthly quota (D1 query)
         └── 5. Proxy to Go backend (X-Backend-Secret)
                    ↓
              Go Backend response
                    ↓
         ├── Record usage async (D1 insert, fire-and-forget)
         └── Add usage headers → return to user
```

## Endpoints

| Method | Path       | Description                                  |
| ------ | ---------- | -------------------------------------------- |
| `GET`  | `/healthz` | Health check                                 |
| `*`    | `/*`       | Proxy to Go backend (requires valid API key) |

## Authentication

All requests (except `/healthz`) require the `requiems-api-key` header:

```
requiems-api-key: requiem_xxxxxxxxxxxxxxxxxxxxxxxx
```

Error responses:

| Status | Reason                         |
| ------ | ------------------------------ |
| `401`  | Missing or invalid API key     |
| `429`  | Per-minute rate limit exceeded |
| `429`  | Monthly quota exceeded         |

## Response Headers

Every proxied response includes usage information:

| Header                  | Description                               |
| ----------------------- | ----------------------------------------- |
| `X-Requests-Used`       | Credits used this billing period          |
| `X-Requests-Remaining`  | Credits remaining this billing period     |
| `X-Requests-Reset`      | Billing period reset timestamp (ISO 8601) |
| `X-Plan`                | User's current plan                       |
| `X-RateLimit-Limit`     | Per-minute rate limit for the plan        |
| `X-RateLimit-Remaining` | Remaining requests in the current minute  |

## Structure

```
src/
├── index.ts                  # Hono app entry point
├── http.ts                   # filterHeaders, addUsageHeaders, fetchBackend
├── rate-limit.ts             # Per-minute rate limiting (KV counters)
├── requests.ts               # Quota checks and usage recording (D1)
├── shared/
│   └── env.ts                # WorkerBindings type + env schema
├── middleware/
│   └── api-key-auth.ts       # Auth middleware (validates key, checks limits)
└── routes/
    └── proxy.ts              # Wildcard proxy handler
```

## Data Storage

### Cloudflare KV

| Key                      | Value                                               | TTL       |
| ------------------------ | --------------------------------------------------- | --------- |
| `key:{apiKey}`           | `ApiKeyData` JSON (userId, plan, billingCycleStart) | No expiry |
| `rl:m:{apiKey}:{minute}` | Request count                                       | 60s       |

### Cloudflare D1 (`credit_usage` table)

Records every API call:

```sql
CREATE TABLE credit_usage (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  api_key      TEXT    NOT NULL,
  user_id      TEXT    NOT NULL,
  endpoint     TEXT    NOT NULL,
  credits_used INTEGER NOT NULL DEFAULT 1,
  used_at      TEXT    NOT NULL  -- ISO 8601
);
```

## Environment Variables

Set via `wrangler secret put`:

| Variable         | Description                                    |
| ---------------- | ---------------------------------------------- |
| `BACKEND_URL`    | Internal URL of the Go backend                 |
| `BACKEND_SECRET` | Secret value sent as `X-Backend-Secret` header |

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

# Or run migrations
pnpm run db:migrate
```

### 4. Set Secrets

```bash
wrangler secret put BACKEND_URL
# Enter your Go backend URL (e.g., https://api-internal.requiems.xyz)

wrangler secret put BACKEND_SECRET
# Enter a strong secret: openssl rand -base64 32
```

### 5. Seed KV with Test API Keys

```bash
pnpm run kv:seed
```

## Development

```bash
pnpm dev                        # Start local dev server (port 4455)
pnpm exec vitest run            # Run tests
pnpm exec vitest run --coverage # Tests with coverage
pnpm run typecheck              # TypeScript type check
pnpm run lint                   # Lint code
pnpm run lint:fix               # Auto-fix lint issues
pnpm run format                 # Format code
pnpm run format:check           # Check formatting
```

### Database

```bash
pnpm run db:migrate       # Run D1 migrations (local)
pnpm run db:migrate:prod  # Run D1 migrations (production)
```

## Deployment

```bash
pnpm run deploy       # Deploy to staging
pnpm run deploy:prod  # Deploy to production
```

## Tests

```
src/__tests__/
├── config.test.ts      # Plan limits and endpoint multipliers
├── http.test.ts        # Header filtering and usage headers
└── rate-limit.test.ts  # Rate limiting logic
```
