/**
 * Concurrent User Simulation — tests/load/scenarios/concurrent-users.ts
 *
 * Purpose
 * -------
 * Simulate realistic multi-user traffic patterns against the full API stack
 * (Auth Gateway → Go backend) using k6's ramping-vus executor.
 *
 * Traffic shape
 * -------------
 *  Stage 1 — Ramp-up     (0 → 10 VUs  over  30 s)  warm up the stack
 *  Stage 2 — Sustained   (10 VUs      for  1 m)     steady-state load
 *  Stage 3 — Peak        (10 → 50 VUs over  30 s)   traffic spike
 *  Stage 4 — Peak hold   (50 VUs      for  1 m)     sustained peak
 *  Stage 5 — Ramp-down   (50 → 0  VUs over  30 s)   graceful wind-down
 *
 * Total duration ≈ 4 minutes.
 *
 * Each VU randomly picks one of the sample endpoints on every iteration,
 * simulating diverse organic traffic rather than a hammer on a single route.
 *
 * Thresholds
 * ----------
 *  - Error rate   < 1 %
 *  - p95 latency  < 500 ms
 *  - p99 latency  < 1 000 ms
 *  - Throughput   recorded in requests_per_second metric
 *
 * Usage
 * -----
 *   k6 run tests/load/scenarios/concurrent-users.ts
 *
 *   # Override concurrency
 *   PEAK_VUS=100 k6 run tests/load/scenarios/concurrent-users.ts
 */

import http from "k6/http";
import { check, sleep } from "k6";
import { Trend, Rate, Counter } from "k6/metrics";
import { Options } from "k6/options";

import {
  BASE_URL,
  API_KEYS,
  DEFAULT_THRESHOLDS,
  SAMPLE_ENDPOINTS,
  SummaryData,
} from "../config.ts";

// ---------------------------------------------------------------------------
// Custom metrics
// ---------------------------------------------------------------------------

const successRate = new Rate("success_rate");
const endpointDuration = new Trend("endpoint_duration_concurrent", true);
const totalRequests = new Counter("total_requests");

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

/** Maximum (peak) virtual users.  Override with PEAK_VUS env var. */
const PEAK_VUS = parseInt(__ENV["PEAK_VUS"] ?? "50", 10);

/**
 * Round-robin through plan keys so the load is spread across multiple users,
 * preventing any single key from exhausting its monthly quota during tests.
 * All keys used here should be developer-tier or above (high rate limits).
 */
const PLAN_KEYS: string[] = [
  API_KEYS.developer,
  API_KEYS.business,
  API_KEYS.professional,
];

// ---------------------------------------------------------------------------
// k6 options
// ---------------------------------------------------------------------------

export const options: Options = {
  scenarios: {
    concurrent_users: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "30s", target: 10 }, // ramp-up
        { duration: "1m", target: 10 }, // steady
        { duration: "30s", target: PEAK_VUS }, // spike
        { duration: "1m", target: PEAK_VUS }, // peak hold
        { duration: "30s", target: 0 }, // ramp-down
      ],
      gracefulRampDown: "15s",
    },
  },

  thresholds: {
    ...DEFAULT_THRESHOLDS,
    success_rate: ["rate>0.99"],
  },
};

// ---------------------------------------------------------------------------
// Default function
// ---------------------------------------------------------------------------

export default function (): void {
  // Pick an endpoint at random for realistic diverse traffic
  const endpoint =
    SAMPLE_ENDPOINTS[Math.floor(Math.random() * SAMPLE_ENDPOINTS.length)];

  // Round-robin through high-capacity plan keys using VU id
  const apiKey = PLAN_KEYS[(__VU - 1) % PLAN_KEYS.length];
  const params = {
    headers: {
      "requiems-api-key": apiKey,
      "Content-Type": "application/json",
    },
    tags: {
      endpoint: `${endpoint.method} ${endpoint.path}`,
    },
  };

  const url = `${BASE_URL}${endpoint.path}`;

  const res =
    endpoint.method === "POST"
      ? http.post(url, endpoint.body ?? null, params)
      : http.get(url, params);

  totalRequests.add(1);
  endpointDuration.add(res.timings.duration, {
    endpoint: `${endpoint.method} ${endpoint.path}`,
  });

  const ok = check(res, {
    "status 2xx": (r) => r.status >= 200 && r.status < 300,
    "no server error": (r) => r.status < 500,
    "response body present": (r) => (r.body as string).length > 0,
  });

  successRate.add(ok ? 1 : 0);

  if (!ok) {
    // Only log non-5xx details to avoid flooding output on expected 4xx
    if (res.status >= 500) {
      console.error(
        `[concurrent] Server error ${res.status} on ${endpoint.method} ${endpoint.path}: ${(res.body as string).substring(0, 300)}`,
      );
    } else if (res.status === 429) {
      // Rate limiting is tracked but not unexpected for high-concurrency tests
      // against free-plan keys (shouldn't happen here as we use paid tiers)
      console.warn(
        `[concurrent] 429 Rate Limited on ${endpoint.method} ${endpoint.path} (VU ${__VU})`,
      );
    }
  }

  // Think time: 0.5–2 s simulates a real user pausing between requests
  sleep(0.5 + Math.random() * 1.5);
}

// ---------------------------------------------------------------------------
// Summary
// ---------------------------------------------------------------------------

export function handleSummary(data: SummaryData): Record<string, string> {
  const m = data.metrics;

  const p95 = m["http_req_duration"]?.values["p(95)"] ?? "N/A";
  const p99 = m["http_req_duration"]?.values["p(99)"] ?? "N/A";
  const rps = m["http_reqs"]?.values["rate"] ?? "N/A";
  const failRate = ((m["http_req_failed"]?.values["rate"] ?? 0) * 100).toFixed(
    2,
  );
  const successRateVal = (
    (m["success_rate"]?.values["rate"] ?? 0) * 100
  ).toFixed(2);
  const total = m["total_requests"]?.values["count"] ?? "N/A";

  const fmtMs = (v: number | string): string =>
    typeof v === "number" ? v.toFixed(2) + " ms" : v;
  const fmtRps = (v: number | string): string =>
    typeof v === "number" ? v.toFixed(2) + " req/s" : v;

  const summary = [
    "=".repeat(60),
    "  CONCURRENT USER SIMULATION RESULTS",
    "=".repeat(60),
    `  Target URL      : ${BASE_URL}`,
    `  Peak VUs        : ${PEAK_VUS}`,
    `  Total requests  : ${total}`,
    `  Throughput      : ${fmtRps(rps)}`,
    `  p95 latency     : ${fmtMs(p95)}`,
    `  p99 latency     : ${fmtMs(p99)}`,
    `  Error rate      : ${failRate} %`,
    `  Success rate    : ${successRateVal} %`,
    "=".repeat(60),
  ].join("\n");

  console.log(summary);

  return { stdout: summary + "\n" };
}
