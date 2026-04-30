/**
 * Vitest global setup — runs once in the main process before any test workers
 * start.
 *
 * Clears the stats temp file so that a fresh run never inherits timing data
 * from a previous run. Without this, the merge logic in stats.persist() would
 * accumulate data across multiple invocations.
 */

import fs from "node:fs";
import { STATS_TEMP_FILE } from "./stats.js";

export function setup(): void {
  if (fs.existsSync(STATS_TEMP_FILE)) {
    fs.unlinkSync(STATS_TEMP_FILE);
  }
}
