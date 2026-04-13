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
  expect(body, "Response must have a 'metadata' key").toHaveProperty("metadata");

  assertShape(snapshotName, endpointKey, body);

  return { data: body["data"], metadata: body["metadata"] };
}

/**
 * Run an async action `runs` times (from config) and return all results.
 * Useful for warming up timing samples and detecting flakiness.
 */
export async function repeat<T>(action: () => Promise<T>): Promise<T[]> {
  const { runs } = getConfig();
  const results: T[] = [];
  for (let i = 0; i < runs; i++) {
    results.push(await action());
  }
  return results;
}
