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
 * Run an async action `runs` times (from config) and return successful results.
 * Network errors and timeouts are swallowed so a flaky backend doesn't abort
 * the entire benchmark run — partial data is still recorded in the stats table.
 * Throws only if every single attempt fails.
 */
export async function repeat<T>(action: () => Promise<T>): Promise<T[]> {
  const { runs } = getConfig();
  const results: T[] = [];
  let failures = 0;

  for (let i = 0; i < runs; i++) {
    try {
      results.push(await action());
    } catch (err) {
      failures++;
      const msg = err instanceof Error ? err.message : String(err);
      process.stderr.write(`    ⚠  attempt ${i + 1}/${runs} failed: ${msg}\n`);
    }
  }

  if (failures === runs) {
    throw new Error(`All ${runs} attempts failed — backend may be down.`);
  }

  return results;
}
