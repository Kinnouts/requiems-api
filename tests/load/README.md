# Load Testing Suite

Performance tests for the Requiem API using [k6](https://k6.io/), written in
TypeScript.

## Overview

The suite covers four tasks from the original issue:

| Script                                                             | Purpose                                                  |
| ------------------------------------------------------------------ | -------------------------------------------------------- |
| [`scenarios/baseline.ts`](./scenarios/baseline.ts)                 | Baseline benchmarks — single VU, all service groups      |
| [`scenarios/rate-limit.ts`](./scenarios/rate-limit.ts)             | Rate limit validation — enforces per-plan req/min caps   |
| [`scenarios/concurrent-users.ts`](./scenarios/concurrent-users.ts) | Concurrent user simulation — ramping VUs up to 50 (peak) |

Configuration shared across all scripts lives in [`config.ts`](./config.ts).

## Prerequisites

1. **k6** — `brew install k6` (macOS) or see
   [k6 installation docs](https://k6.io/docs/get-started/installation/). k6
   v0.46+ supports TypeScript natively via its built-in esbuild bundler.
2. **Dev stack running** — the Auth Gateway must be reachable at
   `localhost:4455`:

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml up
```

3. **Dev API keys seeded** — this is done automatically when the stack starts,
   or manually:

```bash
cd apps/workers/auth-gateway
node ./scripts/seed-dev.ts
```

## Running Tests

```bash
# Run all scenarios in sequence
./tests/load/run.sh

# Run a single scenario
./tests/load/run.sh baseline
./tests/load/run.sh rate-limit
./tests/load/run.sh concurrent-users

# Override target URL (e.g. staging)
BASE_URL=https://api-staging.example.com ./tests/load/run.sh all

# Increase concurrent users for the concurrent-users scenario
PEAK_VUS=100 ./tests/load/run.sh concurrent-users

# Forward extra k6 flags
./tests/load/run.sh baseline --out json=/tmp/results.json
```

Or call k6 directly:

```bash
BASE_URL=http://localhost:4455 k6 run tests/load/scenarios/baseline.ts
```

## Environment Variables

| Variable               | Default                 | Description                                 |
| ---------------------- | ----------------------- | ------------------------------------------- |
| `BASE_URL`             | `http://localhost:4455` | Auth Gateway base URL                       |
| `API_KEY_FREE`         | `rq_free_000001`        | Free-plan dev key                           |
| `API_KEY_DEVELOPER`    | `rq_devl_000001`        | Developer-plan dev key                      |
| `API_KEY_BUSINESS`     | `rq_bizz_000001`        | Business-plan dev key                       |
| `API_KEY_PROFESSIONAL` | `rq_prof_000001`        | Professional-plan dev key                   |
| `PEAK_VUS`             | `50`                    | Peak concurrent users in `concurrent-users` |

## Scenarios

### `baseline.ts` — Baseline Benchmarks

- **Executor:** `constant-vus`, 1 VU, 2 minutes
- **Coverage:** All API service groups (text, ai, email, entertainment, misc,
  places, tech, finance, fitness, convert)
- **Checks:** HTTP 2xx, non-empty body, valid JSON
- **Thresholds:**
  - `http_req_failed < 1 %`
  - `http_req_duration p(95) < 500 ms`
  - `http_req_duration p(99) < 1 000 ms`
  - `check_success_rate > 99 %`

Use this output as the reference number when comparing future releases.

### `rate-limit.ts` — Rate Limit Validation

Tests three behaviours in sequence:

| Stage                      | Key                       | Requests | Expected                  |
| -------------------------- | ------------------------- | -------- | ------------------------- |
| `free_burst` (0–20 s)      | free (30 req/min)         | 35 rapid | First ~30 → 2xx, then 429 |
| `developer_safe` (25–45 s) | developer (5 000 req/min) | 60 rapid | All 2xx, no 429           |
| `recovery` (95–115 s)      | free (window reset)       | 5        | All 2xx again             |

**Thresholds:**

- `rate_limit_hits > 0` — the free burst must have triggered at least one 429
- `unexpected_errors == 0` — no non-429 errors

### `concurrent-users.ts` — Concurrent User Simulation

Traffic shape:

```
VUs
50 │          ┌─────────────┐
   │         /               \
10 │   ┌────┘                 └─────┐
 0 │───┘                             └
   0   30s   90s  120s  180s  210s  240s
```

- Each VU picks a random endpoint on every iteration (organic traffic mix).
- Uses developer / business / professional keys (round-robin) to avoid
  exhausting the free quota.
- Think time: 0.5–2 s between requests.

**Thresholds:** same as baseline plus `success_rate > 99 %`.

## Interpreting Results

k6 prints a summary table after each run. Key metrics to watch:

| Metric                    | Good                      | Investigate               |
| ------------------------- | ------------------------- | ------------------------- |
| `http_req_failed`         | < 1 %                     | ≥ 1 %                     |
| `http_req_duration p(95)` | < 500 ms                  | ≥ 500 ms                  |
| `http_req_duration p(99)` | < 1 000 ms                | ≥ 1 000 ms                |
| `rate_limit_hits`         | > 0 (rate-limit scenario) | 0 (gateway not enforcing) |

## Architecture Context

```
k6 VU
  │
  ▼
Auth Gateway :4455  (Cloudflare Worker — auth, rate limit, usage)
  │
  ▼
Go API :8080         (business logic, PostgreSQL, Redis)
```

Rate limits are enforced in the Auth Gateway before the request reaches the Go
backend. See [`docs/core/auth-gateway.md`](../../docs/core/auth-gateway.md) for
details.
