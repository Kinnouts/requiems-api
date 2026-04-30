# API Management

> https://api-management.requiems.xyz/docs

Internal-only Cloudflare Worker at `api-management.requiems.xyz`. Provides a
REST API for the Rails dashboard to manage API keys, export usage data to
PostgreSQL, and query analytics from D1.

## Endpoints

### Health

| Method | Path            | Description                                     |
| ------ | --------------- | ----------------------------------------------- |
| `GET`  | `/healthz`      | Health check                                    |
| `GET`  | `/docs`         | Swagger UI (basic auth protected in production) |
| `GET`  | `/openapi.json` | OpenAPI spec                                    |

### API Keys

| Method   | Path                   | Description             |
| -------- | ---------------------- | ----------------------- |
| `POST`   | `/api-keys`            | Create a new API key    |
| `DELETE` | `/api-keys/:keyPrefix` | Revoke an API key       |
| `PATCH`  | `/api-keys/:keyPrefix` | Update API key metadata |

### Usage Export

| Method | Path            | Description                              |
| ------ | --------------- | ---------------------------------------- |
| `GET`  | `/usage/export` | Export usage records for PostgreSQL sync |

Query params for `/usage/export`:

| Param    | Required | Default | Description                                    |
| -------- | -------- | ------- | ---------------------------------------------- |
| `since`  | Yes      | —       | ISO 8601 timestamp; returns records after this |
| `limit`  | No       | 1000    | Max records per page (max: 5000)               |
| `cursor` | No       | —       | Pagination offset                              |

Returns paginated `UsageRecord[]` with `hasMore` and `nextCursor`.

### Analytics

| Method | Path                     | Description                      |
| ------ | ------------------------ | -------------------------------- |
| `GET`  | `/analytics/summary`     | Overall usage summary for a user |
| `GET`  | `/analytics/by-endpoint` | Usage breakdown by endpoint      |
| `GET`  | `/analytics/by-date`     | Usage breakdown by date          |

## Authentication

All management endpoints require the `X-API-Management-Key` header. Only the
Rails dashboard should have this key.

```
X-API-Management-Key: <secret>
```

The Swagger UI at `/docs` is protected with HTTP Basic Auth in production.

## Structure

```
src/
├── index.ts                      # Hono app entry point, Swagger setup
├── env.ts                        # WorkerBindings type + env schema
├── middleware/
│   └── api-key-auth.ts           # X-API-Management-Key validation
└── routes/
    ├── api-keys/
    │   ├── index.ts              # Route orchestrator
    │   ├── create.ts             # POST /api-keys
    │   ├── delete.ts             # DELETE /api-keys/:keyPrefix
    │   └── patch.ts              # PATCH /api-keys/:keyPrefix
    ├── usage/
    │   ├── index.ts              # Route orchestrator
    │   ├── export.ts             # GET /usage/export
    │   └── types.ts              # UsageRecord types
    ├── analytics/
    │   ├── index.ts              # Route orchestrator
    │   ├── summary.ts            # GET /analytics/summary
    │   ├── by-endpoint.ts        # GET /analytics/by-endpoint
    │   ├── by-date.ts            # GET /analytics/by-date
    │   └── types.ts              # Analytics types
    └── swagger.ts                # OpenAPI spec
```

## API Key Lifecycle

```
Rails Dashboard
    ↓  POST /api-keys
API Management creates key:
  1. Generate full key: requiem_<24 random chars>
  2. Store in KV: key:{fullKey} → ApiKeyData
  3. Store prefix in D1 for audit trail
  4. Return full key to Rails (only time it's visible)
    ↓
Rails stores key prefix + metadata in PostgreSQL
    ↓
User receives full key from Rails
```

## Usage Sync

Rails background jobs call `/usage/export` to sync D1 → PostgreSQL:

```
Rails Sidekiq Job (scheduled)
    ↓  GET /usage/export?since=<last_sync>&limit=1000
API Management queries D1:
  SELECT * FROM credit_usage WHERE used_at > since ORDER BY used_at
    ↓
Rails upserts records into PostgreSQL usage_logs table
    ↓
Repeat with cursor until hasMore = false
```

## Environment Variables

Set via `wrangler secret put`:

| Variable                 | Description                                                |
| ------------------------ | ---------------------------------------------------------- |
| `API_MANAGEMENT_API_KEY` | Secret key; Rails must send this as `X-API-Management-Key` |
| `SWAGGER_USERNAME`       | Basic auth username for `/docs`                            |
| `SWAGGER_PASSWORD`       | Basic auth password for `/docs`                            |

## Development

```bash
pnpm dev                        # Start local dev server (port 6001)
pnpm exec vitest run            # Run tests
pnpm exec vitest run --coverage # Tests with coverage
pnpm run typecheck              # TypeScript type check
pnpm run lint                   # Lint code
pnpm run lint:fix               # Auto-fix lint issues
pnpm run format                 # Format code
pnpm run format:check           # Check formatting
```

## Deployment

```bash
pnpm run deploy       # Deploy to staging
pnpm run deploy:prod  # Deploy to production
```
