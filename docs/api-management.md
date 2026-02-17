# API Management Service

Internal Cloudflare Worker service for managing API keys, exporting usage data, and providing analytics.

**URL:** https://api-management.requiems.xyz
**Local:** http://localhost:6001
**Framework:** Hono (TypeScript)
**Runtime:** Cloudflare Workers (Bun for local dev)

## Overview

The API Management service is an internal-only service that handles:

- **API Key Management** - Create, revoke, and update API keys in KV + D1
- **Usage Export** - Export usage data from D1 for Rails PostgreSQL sync
- **Analytics** - Query usage statistics by endpoint, date, and user

This service is separate from the public-facing Auth Gateway, allowing independent scaling and clearer security boundaries.

## Authentication

All endpoints require the `X-API-Management-Key` header. Only the Rails dashboard has this key.

```bash
curl https://api-management.requiems.xyz/healthz \
  -H "X-API-Management-Key: your-64-char-secret-key"
```

**Security Notes:**
- This key is different from `X-Backend-Secret` (used by auth-gateway → Go backend)
- Store the key in environment variables (`API_MANAGEMENT_API_KEY`)
- Rotate the key if compromised using `wrangler secret put API_MANAGEMENT_API_KEY`

## Endpoints

### Health Check

**GET /healthz**

Check if the service is running. No authentication required.

```bash
curl http://localhost:6001/healthz
```

Response:
```json
{
  "status": "ok",
  "service": "api-management"
}
```

---

### API Key Management

**POST /api-keys**

Create, revoke, or update API keys.

**Headers:**
- `X-API-Management-Key` (required)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "action": "create" | "revoke" | "update",
  "key": "rq_live_abc123...",
  "userId": "user-uuid",
  "plan": "free" | "developer" | "business" | "professional" | "enterprise",
  "billingCycleStart": "2025-01-01T00:00:00Z" // Optional, ISO 8601
}
```

**Examples:**

Create API key:
```bash
curl -X POST http://localhost:6001/api-keys \
  -H "X-API-Management-Key: $API_MANAGEMENT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "create",
    "key": "rq_live_abc123xyz",
    "userId": "user-123",
    "plan": "developer",
    "billingCycleStart": "2025-02-01T00:00:00Z"
  }'
```

Revoke API key:
```bash
curl -X POST http://localhost:6001/api-keys \
  -H "X-API-Management-Key: $API_MANAGEMENT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "revoke",
    "key": "rq_live_abc123xyz",
    "userId": "user-123",
    "plan": "developer"
  }'
```

Update API key (plan change):
```bash
curl -X POST http://localhost:6001/api-keys \
  -H "X-API-Management-Key: $API_MANAGEMENT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "update",
    "key": "rq_live_abc123xyz",
    "userId": "user-123",
    "plan": "business",
    "billingCycleStart": "2025-02-01T00:00:00Z"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "API key created successfully"
}
```

**Error Responses:**
- `401` - Invalid or missing API management key
- `404` - API key not found (for revoke/update)
- `409` - API key already exists (for create)
- `400` - Missing required fields

---

### Usage Export

**GET /usage/export**

Export usage data from D1 for syncing to Rails PostgreSQL. Supports pagination.

**Headers:**
- `X-API-Management-Key` (required)

**Query Parameters:**
- `since` (required) - ISO 8601 timestamp, get records after this time
- `limit` (optional) - Max records to return (default: 1000, max: 5000)
- `cursor` (optional) - Pagination cursor (offset)

**Example:**
```bash
curl "http://localhost:6001/usage/export?since=2025-02-01T00:00:00Z&limit=1000" \
  -H "X-API-Management-Key: $API_MANAGEMENT_API_KEY"
```

**Response:**
```json
{
  "usage": [
    {
      "api_key": "rq_live_abc123xyz",
      "endpoint": "/v1/text/advice",
      "credits_used": 1,
      "used_at": "2025-02-17T10:30:00Z"
    }
  ],
  "total": 5432,
  "hasMore": true,
  "nextCursor": "1000"
}
```

**Pagination:**
To get the next page, use the `nextCursor` value:
```bash
curl "http://localhost:6001/usage/export?since=2025-02-01T00:00:00Z&limit=1000&cursor=1000" \
  -H "X-API-Management-Key: $API_MANAGEMENT_API_KEY"
```

---

### Analytics: By Endpoint

**GET /analytics/by-endpoint**

Get usage breakdown by endpoint for a user.

**Headers:**
- `X-API-Management-Key` (required)

**Query Parameters:**
- `userId` (required) - User ID
- `since` (optional) - ISO 8601 timestamp (defaults to billing cycle start)
- `until` (optional) - ISO 8601 timestamp (defaults to now)
- `limit` (optional) - Max top endpoints to return (default: 10, max: 100)

**Example:**
```bash
curl "http://localhost:6001/analytics/by-endpoint?userId=user-123&limit=5" \
  -H "X-API-Management-Key: $API_MANAGEMENT_API_KEY"
```

**Response:**
```json
{
  "endpoints": [
    {
      "endpoint": "/v1/text/words/define",
      "requests": 1250,
      "credits": 2500
    },
    {
      "endpoint": "/v1/text/advice",
      "requests": 980,
      "credits": 980
    }
  ],
  "dateRange": {
    "since": "2025-02-01T00:00:00Z",
    "until": "2025-02-17T18:00:00Z"
  }
}
```

---

### Analytics: By Date

**GET /analytics/by-date**

Get usage trends over time for a user.

**Headers:**
- `X-API-Management-Key` (required)

**Query Parameters:**
- `userId` (required) - User ID
- `since` (optional) - ISO 8601 timestamp (defaults to 30 days ago)
- `until` (optional) - ISO 8601 timestamp (defaults to now)
- `groupBy` (optional) - Grouping: "day" or "hour" (default: "day")

**Example:**
```bash
curl "http://localhost:6001/analytics/by-date?userId=user-123&groupBy=day" \
  -H "X-API-Management-Key: $API_MANAGEMENT_API_KEY"
```

**Response:**
```json
{
  "timeSeries": [
    {
      "date": "2025-02-15",
      "requests": 150,
      "credits": 220
    },
    {
      "date": "2025-02-16",
      "requests": 175,
      "credits": 250
    }
  ],
  "dateRange": {
    "since": "2025-01-18T00:00:00Z",
    "until": "2025-02-17T18:00:00Z"
  },
  "groupBy": "day"
}
```

---

### Analytics: Summary

**GET /analytics/summary**

Get overall usage summary for a user with top endpoints.

**Headers:**
- `X-API-Management-Key` (required)

**Query Parameters:**
- `userId` (required) - User ID
- `since` (optional) - ISO 8601 timestamp (defaults to billing cycle start)
- `until` (optional) - ISO 8601 timestamp (defaults to now)

**Example:**
```bash
curl "http://localhost:6001/analytics/summary?userId=user-123" \
  -H "X-API-Management-Key: $API_MANAGEMENT_API_KEY"
```

**Response:**
```json
{
  "userId": "user-123",
  "totalRequests": 2450,
  "totalCredits": 3680,
  "dateRange": {
    "since": "2025-02-01T00:00:00Z",
    "until": "2025-02-17T18:00:00Z"
  },
  "topEndpoints": [
    {
      "endpoint": "/v1/text/words/define",
      "requests": 1250,
      "credits": 2500
    },
    {
      "endpoint": "/v1/text/advice",
      "requests": 980,
      "credits": 980
    }
  ]
}
```

---

### Swagger Documentation

**GET /docs**

Interactive API documentation using Swagger UI.

**Authentication (Production):**
- Protected with HTTP Basic Auth in production
- Credentials: `SWAGGER_USERNAME` and `SWAGGER_PASSWORD` (from secrets)
- No auth required in development

**Access:**
```bash
# Local development
open http://localhost:6001/docs

# Production (basic auth prompt will appear)
open https://api-management.requiems.xyz/docs
```

**OpenAPI Spec:**
```bash
curl http://localhost:6001/openapi.json
```

---

## Local Development

### Setup

1. Install dependencies:
```bash
cd apps/api-management
bun install
```

2. Set environment variables:
```bash
# Create .env file
echo "API_MANAGEMENT_API_KEY=your-local-dev-key-min-32-chars" > .env
echo "ENVIRONMENT=development" >> .env
```

3. Start dev server:
```bash
bun dev
```

Server runs at http://localhost:6001

### Testing

Run tests:
```bash
bun run test
```

Run tests with coverage:
```bash
bun run test:coverage
```

Watch mode:
```bash
bun run test:watch
```

TypeScript check:
```bash
bun run typecheck
```

Linting:
```bash
bun run lint
bun run lint:fix  # Auto-fix
```

---

## Deployment

### Prerequisites

- Wrangler CLI installed: `bun install -g wrangler`
- Cloudflare account with Workers enabled
- KV namespace and D1 database created (shared with auth-gateway)

### Set Secrets

```bash
cd apps/api-management

# Required: API Management key (64+ character random string)
wrangler secret put API_MANAGEMENT_API_KEY

# Optional: Swagger basic auth (production only)
wrangler secret put SWAGGER_USERNAME
wrangler secret put SWAGGER_PASSWORD
```

### Deploy to Production

```bash
bun run deploy:prod
```

This deploys to `requiem-api-management-production` worker.

### Verify Deployment

```bash
# Health check (no auth)
curl https://api-management.requiems.xyz/healthz

# Test with API key
curl https://api-management.requiems.xyz/healthz \
  -H "X-API-Management-Key: $API_MANAGEMENT_API_KEY"
```

### Monitor

View logs:
```bash
wrangler tail requiem-api-management-production
```

---

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `API_MANAGEMENT_API_KEY` | Yes | Secret key for API authentication (min 32 chars) |
| `SWAGGER_USERNAME` | No | Basic auth username for /docs (production only) |
| `SWAGGER_PASSWORD` | No | Basic auth password for /docs (production only) |
| `ENVIRONMENT` | No | Environment name (development/staging/production) |

**Binding Variables (wrangler.toml):**
- `KV` - Cloudflare KV namespace (shared with auth-gateway)
- `DB` - Cloudflare D1 database (shared with auth-gateway)

---

## Architecture

### Data Flow

```
Rails Dashboard
    |
    | 1. Create API key via POST /api-keys
    v
API Management
    |
    ├─> 2. Write to KV: key:{api_key} → {userId, plan, ...}
    └─> 3. Write to D1: api_keys table (audit trail)

Auth Gateway (separate service)
    |
    | 4. User makes request with API key
    ├─> 5. Read from KV: validate key
    └─> 6. Write to D1: credit_usage table (async)

API Management
    |
    | 7. Rails pulls usage via GET /usage/export
    └─> 8. Query D1: credit_usage WHERE used_at >= ?
```

### Shared Resources

Both auth-gateway and api-management share:
- **Same KV namespace** - For API key storage
- **Same D1 database** - For usage tracking and audit

**KV Schema:**
- `key:{api_key}` → `{userId, plan, createdAt, billingCycleStart}`
- `rl:m:{api_key}:{minute}` → request count (managed by auth-gateway)

**D1 Schema:**
- `credit_usage` - Usage records (written by auth-gateway, read by api-management)
- `api_keys` - API key metadata (written by api-management)

---

## Troubleshooting

### "Unauthorized" Error

**Problem:** API returns 401 Unauthorized

**Solutions:**
1. Check header: `X-API-Management-Key` (not `X-Backend-Secret`)
2. Verify key value matches secret: `wrangler secret list`
3. Check key length (must be 32+ characters)

### Usage Export Returns Empty

**Problem:** `/usage/export` returns empty array

**Solutions:**
1. Check `since` parameter is not in the future
2. Verify D1 database has data: `wrangler d1 execute requiem-usage --command "SELECT COUNT(*) FROM credit_usage"`
3. Check auth-gateway is recording usage

### Swagger Docs Show 401

**Problem:** `/docs` prompts for auth in development

**Solutions:**
1. Check `ENVIRONMENT` is set to "development" in .env
2. Basic auth only applies when all three are true:
   - `ENVIRONMENT === "production"`
   - `SWAGGER_USERNAME` is set
   - `SWAGGER_PASSWORD` is set

---

## Related Documentation

- [Auth Gateway](./auth-gateway.md) - Public API gateway documentation
- [Architecture](./architecture.md) - System architecture overview
- [CLAUDE.md](../CLAUDE.md) - Full project development guide
