/**
 * Shared helpers for integration test suites.
 */

import { expect } from "vitest";
import { getConfig } from "./config.js";
import { assertShape } from "./snapshot.js";

/**
 * Assert that a response has a standard API envelope:
 *   { data: ..., metadata: { ... } }
 */
export async function assertEnvelope(
  response: Response,
  snapshotName: string,
  endpointKey: string,
): Promise<{ data: unknown; metadata: unknown }> {
  expect(response.status, `Expected 200 OK for ${endpointKey}`).toBe(200);

  const body = (await response.json()) as Record<string, unknown>;

  expect(body, "Response must have a 'data' key").toHaveProperty("data");
  expect(body, "Response must have a 'metadata' key").toHaveProperty(
    "metadata",
  );

  assertShape(snapshotName, endpointKey, body);

  return { data: body["data"], metadata: body["metadata"] };
}

/**
 * Run an async action `runs` times (from config) and return results.
 *
 * Behaviour differs by mode:
 *
 * - **Test mode** (default): any error propagates immediately so the test
 *   fails fast (~8 s with AbortSignal) rather than hanging for 30 s.
 *
 * - **Benchmark mode** (`UPDATE_PERF_BASELINE=true`): errors are caught so
 *   partial data is still collected from a flaky backend.  However, the loop
 *   bails out after MAX_CONSECUTIVE consecutive failures to avoid 50 × 8 s of
 *   wasted retries when an endpoint is completely unreachable.
 */
const MAX_CONSECUTIVE = 3;

export async function repeat<T>(action: () => Promise<T>): Promise<T[]> {
  const { runs } = getConfig();
  const isBenchmark = process.env["UPDATE_PERF_BASELINE"] === "true";
  const results: T[] = [];
  let consecutiveFailures = 0;

  for (let i = 0; i < runs; i++) {
    try {
      results.push(await action());
      consecutiveFailures = 0;
    } catch (err) {
      if (!isBenchmark) throw err;

      consecutiveFailures++;
      const msg = err instanceof Error ? err.message : String(err);
      process.stderr.write(`    ⚠  attempt ${i + 1}/${runs} failed: ${msg}\n`);

      if (consecutiveFailures >= MAX_CONSECUTIVE) break;
    }
  }

  if (results.length === 0) {
    throw new Error(`All attempts failed — backend may be down.`);
  }

  return results;
}
