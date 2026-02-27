# Shared

Common types, utilities, and middleware used by all Cloudflare Workers in this
monorepo. Imported via the `@requiem/workers-shared` path alias.

## Structure

```
src/
├── index.ts              # Main exports
├── types.ts              # Core TypeScript interfaces
├── config.ts             # Plan definitions and endpoint multipliers
├── http.ts               # HTTP response utilities
├── logger.ts             # Structured logging
├── api-key-generator.ts  # API key generation
└── middleware/
    ├── index.ts
    ├── basic-auth.ts     # Basic auth for Swagger docs
    ├── cors.ts           # CORS preflight handling
    └── error-handler.ts  # Global error handler
```

## Modules

### `types.ts`

Core TypeScript interfaces shared across workers:

- `BaseWorkerBindings` — KV namespace, D1 database, and environment bindings
- `ApiKeyData` — API key metadata stored in KV (userId, plan, billingCycleStart)
- `PlanConfig` — Plan limits (requestLimit, ratePerMinute)
- `PlanName` — Union type of all plan names
- `RateLimitResult` — Result of a rate limit check
- `RequestCheckResult` — Result of a quota check
- `ApiKeyManagementRequest/Response` — Types for Rails ↔ API Management
  communication

### `config.ts`

Plan definitions and endpoint cost multipliers:

- `PLANS` — Map of plan name → limits (requestLimit, ratePerMinute)
- `ENDPOINT_MULTIPLIERS` — Map of endpoint path → credit cost multiplier
- `getRequestMultiplier(path)` — Returns how many credits an endpoint costs

Most endpoints cost 1 credit. Expensive endpoints cost more:

| Endpoint                      | Credits |
| ----------------------------- | ------- |
| `GET /v1/text/words/define`   | 2       |
| `GET /v1/text/words/synonyms` | 2       |

### `http.ts`

HTTP response helpers with CORS headers:

- `jsonResponse(data, status?)` — Successful JSON response
- `jsonError(message, status)` — Error JSON response
- `corsResponse` — CORS preflight response (for OPTIONS requests)
- `CORS_HEADERS` — CORS header constants

### `logger.ts`

Structured logging for Cloudflare Workers with cf-ray tracing:

- `createLogger(cfRay?)` — Creates a logger instance bound to a request's cf-ray
  ID
- `maskApiKey(key)` — Masks API keys for safe logging (shows first 8 chars only)
- Outputs JSON formatted logs compatible with Cloudflare Workers Logs

### `api-key-generator.ts`

API key generation utilities:

- `generateApiKey()` — Creates keys in format `requiem_<24_random_chars>`
- `extractKeyPrefix(key)` — Gets first 12 characters for storage
- `isValidKeyFormat(key)` — Validates key format

### `middleware/`

Reusable Hono middleware:

- `errorHandler` — Global error handler; catches unhandled errors and returns
  structured JSON
- `basicAuth(username, password)` — Basic Auth middleware for protecting docs
  endpoints
- `corsMiddleware` — Handles CORS preflight (OPTIONS) requests

## Usage

Workers import from this package using the path alias defined in
`tsconfig.json`:

```ts
import { jsonError, jsonResponse } from "@requiem/workers-shared";
import type { ApiKeyData, PlanName } from "@requiem/workers-shared";
import { getRequestMultiplier } from "@requiem/workers-shared";
```

## Development

```bash
pnpm run typecheck   # Type check
pnpm run lint        # Lint code
pnpm run lint:fix    # Auto-fix lint issues
pnpm run format      # Auto-format code
```
