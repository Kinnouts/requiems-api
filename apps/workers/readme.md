# Workers

Cloudflare Workers powering the Requiems API edge layer. These workers handle
authentication, rate limiting, request proxying, API key management, and
analytics — all running at the edge without a traditional server.

## Packages

| Package                             | Port | Domain                        | Description                                            |
| ----------------------------------- | ---- | ----------------------------- | ------------------------------------------------------ |
| [auth-gateway](./auth-gateway/)     | 4455 | `api.requiems.xyz`            | Public-facing gateway: auth, rate limiting, proxying   |
| [api-management](./api-management/) | 6001 | `api-management.requiems.xyz` | Internal management: API keys, usage export, analytics |
| [shared](./shared/)                 | —    | —                             | Shared types, utilities, and middleware                |

## Architecture

```
Public API Request
        ↓
  Auth Gateway (api.requiems.xyz)
  ├── Validate API key (KV lookup)
  ├── Check per-minute rate limit (KV counter)
  ├── Check monthly quota (D1 query)
  └── Proxy to Go backend (X-Backend-Secret)
        ↓
  Go Backend (internal)
        ↓
  Auth Gateway
  ├── Record usage async (D1 insert)
  └── Add usage headers to response
        ↓
  User receives response

Internal Management (Rails only)
        ↓
  API Management (api-management.requiems.xyz)
  ├── Create / revoke / update API keys (KV + D1)
  ├── Export usage data to PostgreSQL (D1 → PostgreSQL)
  └── Query analytics (D1 aggregations)
```

## Data Storage

| Store         | Type        | Used For                                     |
| ------------- | ----------- | -------------------------------------------- |
| Cloudflare KV | Key-value   | API key lookup, per-minute rate limiting     |
| Cloudflare D1 | Edge SQLite | Usage recording, API key metadata, analytics |
| PostgreSQL    | Relational  | Synced from D1 via Rails background jobs     |

**KV key patterns:**

- `key:{apiKey}` → API key data (userId, plan, billingCycleStart)
- `rl:m:{apiKey}:{minute}` → Rate limit counter (TTL: 60s)

## Plan Limits

| Plan         | Requests/Month | Rate/Minute |
| ------------ | -------------- | ----------- |
| Free         | 500            | 30          |
| Developer    | 100,000        | 5,000       |
| Business     | 1,000,000      | 10,000      |
| Professional | 10,000,000     | 50,000      |
| Enterprise   | Unlimited      | Unlimited   |

## Development

Start all workers locally (from repo root):

```bash
# Auth Gateway
cd apps/workers/auth-gateway
bun dev  # Port 4455

# API Management
cd apps/workers/api-management
bun dev  # Port 6001
```

Run all tests:

```bash
# Auth Gateway
cd apps/workers/auth-gateway
bunx vitest run

# API Management
cd apps/workers/api-management
bunx vitest run
```

## Tech Stack

- **Framework:** [Hono](https://hono.dev/) — lightweight web framework for
  Cloudflare Workers
- **Runtime:** Cloudflare Workers (V8 isolates)
- **Storage:** Cloudflare KV + D1 (SQLite at the edge)
- **Validation:** [Zod](https://zod.dev/)
- **Testing:** [Vitest](https://vitest.dev/) + `@cloudflare/vitest-pool-workers`
- **Linting/Formatting:** [Biome](https://biomejs.dev/)
- **Package Manager:** [Bun](https://bun.sh/)
