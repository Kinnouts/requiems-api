/**
 * Baseline Benchmark — tests/load/scenarios/baseline.ts
 *
 * Purpose
 * -------
 * Establish latency and error-rate baselines across every API service group
 * with a single virtual user producing a steady, low-volume load.  Results
 * form the reference numbers that regression / capacity tests are measured
 * against.
 *
 * What it tests
 * -------------
 * - Each sample endpoint defined in config.ts is exercised in sequence.
 * - HTTP status codes and response shapes are validated.
 * - Response-time histograms (p95 / p99) and failure rate are recorded.
 *
 * Usage
 * -----
 *   k6 run tests/load/scenarios/baseline.ts
 *
 * Override defaults:
 *   BASE_URL=https://api.example.com \
 *   API_KEY_DEVELOPER=rq_devl_xxx \
 *   k6 run tests/load/scenarios/baseline.ts
 */

import http from "k6/http";
import { check, sleep } from "k6";
import { Counter, Rate, Trend } from "k6/metrics";
import { Options } from "k6/options";

import {
  authParams,
  BASE_URL,
  DEFAULT_THRESHOLDS,
  SAMPLE_ENDPOINTS,
  SummaryData,
} from "../config.ts";

// ---------------------------------------------------------------------------
// Custom metrics
// ---------------------------------------------------------------------------

/** Per-endpoint latency breakdown */
const endpointDuration = new Trend("endpoint_duration", true);

/** Count of successful checks */
const successfulChecks = new Counter("successful_checks");

/** Count of failed checks */
const failedChecks = new Counter("failed_checks");

/** Overall check success rate */
const checkSuccessRate = new Rate("check_success_rate");

// ---------------------------------------------------------------------------
// k6 scenario options
// ---------------------------------------------------------------------------

export const options: Options = {
  /**
   * Single VU, low-frequency sweep of all endpoints.
   * Runs for 2 minutes to produce stable p95/p99 measurements.
   */
  scenarios: {
    baseline: {
      executor: "constant-vus",
      vus: 1,
      duration: "2m",
    },
  },

  thresholds: {
    ...DEFAULT_THRESHOLDS,
    check_success_rate: ["rate>0.99"],
  },
};

// ---------------------------------------------------------------------------
// Default function (VU entry point)
// ---------------------------------------------------------------------------

export default function (): void {
  for (const endpoint of SAMPLE_ENDPOINTS) {
    const url = `${BASE_URL}${endpoint.path}`;
    const params = authParams("developer");

    const res = endpoint.method === "POST"
      ? http.post(url, endpoint.body ?? null, params)
      : http.get(url, params);

    // ------------------------------------------------------------------
    // Validate response
    // ------------------------------------------------------------------

    const tag = { endpoint: `${endpoint.method} ${endpoint.path}` };

    // Record endpoint-specific latency
    endpointDuration.add(res.timings.duration, tag);

    const ok = check(res, {
      "status is 2xx": (r) => r.status >= 200 && r.status < 300,
      "response has body": (r) => (r.body as string).length > 0,
      "response is JSON": (r) => {
        try {
          JSON.parse(r.body as string);
          return true;
        } catch {
          return false;
        }
      },
    });

    if (ok) {
      successfulChecks.add(1);
      checkSuccessRate.add(1);
    } else {
      failedChecks.add(1);
      checkSuccessRate.add(0);
      console.error(
        `[baseline] FAIL ${endpoint.method} ${endpoint.path} → HTTP ${res.status}: ${
          (res.body as string).substring(0, 200)
        }`,
      );
    }

    // Brief pause between requests to avoid hammering a single endpoint
    sleep(0.1);
  }

  // Pause between full sweeps
  sleep(1);
}

// ---------------------------------------------------------------------------
// Summary hook — printed after the run
// ---------------------------------------------------------------------------

export function handleSummary(data: SummaryData): Record<string, string> {
  const metrics = data.metrics;

  const p95 = metrics["http_req_duration"]?.values["p(95)"] ?? "N/A";
  const p99 = metrics["http_req_duration"]?.values["p(99)"] ?? "N/A";
  const failRate = (
    (metrics["http_req_failed"]?.values["rate"] ?? 0) * 100
  ).toFixed(2);
  const checkRate = (
    (metrics["check_success_rate"]?.values["rate"] ?? 0) * 100
  ).toFixed(2);

  const summary = [
    "=".repeat(60),
    "  BASELINE BENCHMARK RESULTS",
    "=".repeat(60),
    `  Target URL   : ${BASE_URL}`,
    `  p95 latency  : ${
      typeof p95 === "number" ? p95.toFixed(2) + " ms" : p95
    }`,
    `  p99 latency  : ${
      typeof p99 === "number" ? p99.toFixed(2) + " ms" : p99
    }`,
    `  Error rate   : ${failRate} %`,
    `  Check pass   : ${checkRate} %`,
    "=".repeat(60),
  ].join("\n");

  console.log(summary);

  return {
    stdout: summary + "\n",
  };
}
