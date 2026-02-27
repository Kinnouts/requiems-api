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

Refer to [docs/architecture.md](../../docs/architecture.md) for the overall system architecture and how the workers fit into the ecosystem.

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

Refer to [docs/business.md](../../docs/business.md) for detailed plan limits.

## Development

Refer to [docs/getting-starter.md](../../docs/getting-started.md) for Docker Compose setup and development workflow instructions. The workers are included in the `docker-compose.dev.yml` configuration and will be built and run alongside the Go API, Rails Dashboard, PostgreSQL, and Redis services.

## Tech Stack

- **Framework:** [Hono](https://hono.dev/) — lightweight web framework for
  Cloudflare Workers
- **Runtime:** Cloudflare Workers (V8 isolates)
- **Storage:** Cloudflare KV + D1 (SQLite at the edge)
- **Validation:** [Zod](https://zod.dev/)
- **Testing:** [Vitest](https://vitest.dev/) + `@cloudflare/vitest-pool-workers`
- **Linting/Formatting:** [Biome](https://biomejs.dev/)
- **Package Manager:** [pnpm](https://pnpm.io/)
