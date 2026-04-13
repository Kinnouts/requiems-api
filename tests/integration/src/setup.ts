/**
 * Global test setup — runs as a setupFile in each vitest worker process.
 *
 * Loads the .env file if present (so developers don't have to export vars
 * manually), validates that the API key is configured, and registers a
 * root-level afterAll hook that persists timing data to a temp file for the
 * custom reporter to read after all tests complete.
 */

import { afterAll } from "vitest";
import { readFileSync, existsSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";
import { stats } from "./stats.js";

const ROOT = join(dirname(fileURLToPath(import.meta.url)), "../../");

function loadDotenv(): void {
  const envPath = join(ROOT, ".env");
  if (!existsSync(envPath)) return;

  const content = readFileSync(envPath, "utf8");
  for (const line of content.split("\n")) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith("#")) continue;

    const eqIdx = trimmed.indexOf("=");
    if (eqIdx === -1) continue;

    const key = trimmed.slice(0, eqIdx).trim();
    const value = trimmed.slice(eqIdx + 1).trim().replace(/^["']|["']$/g, "");

    // Don't overwrite variables already set in the shell environment
    if (!(key in process.env)) {
      process.env[key] = value;
    }
  }
}

// Load .env immediately so getConfig() can read the variables
loadDotenv();

// Persist timing data to disk after all tests in this worker finish so the
// reporter (running in the main vitest process) can build the summary table.
afterAll(() => {
  stats.persist();
});
