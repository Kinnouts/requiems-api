/**
 * Per-endpoint response-time statistics.
 *
 * A single global singleton is used so every test file contributes data.
 * Timing data is persisted to a temp file so the custom reporter (which runs
 * in the main vitest process) can read it after the test worker finishes.
 */

import fs from "node:fs";
import os from "node:os";
import path from "node:path";

export interface EndpointStats {
  path: string;
  samples: number;
  min: number;
  max: number;
  avg: number;
  p50: number;
  p95: number;
  p99: number;
  /** All raw measurements (ms) in insertion order */
  raw: number[];
}

/** Temp file where the worker process writes its timing data */
export const STATS_TEMP_FILE = path.join(
  os.tmpdir(),
  "requiem-integration-stats.json",
);

export class Stats {
  private readonly _data = new Map<string, number[]>();

  /** Record a single latency sample for a given path */
  record(path: string, durationMs: number): void {
    let arr = this._data.get(path);
    if (!arr) {
      arr = [];
      this._data.set(path, arr);
    }
    arr.push(durationMs);
  }

  /** Return stats for every path that has at least one sample */
  summarise(data?: Map<string, number[]>): EndpointStats[] {
    const source = data ?? this._data;
    const result: EndpointStats[] = [];

    for (const [path, raw] of source) {
      const sorted = [...raw].sort((a, b) => a - b);
      const n = sorted.length;

      result.push({
        path,
        samples: n,
        min: sorted[0] ?? 0,
        max: sorted[n - 1] ?? 0,
        avg: Math.round(raw.reduce((s, v) => s + v, 0) / n),
        p50: percentile(sorted, 50),
        p95: percentile(sorted, 95),
        p99: percentile(sorted, 99),
        raw,
      });
    }

    // Sort by path for deterministic output
    result.sort((a, b) => a.path.localeCompare(b.path));
    return result;
  }

  /**
   * Persist timing data to the temp file for the reporter to read.
   *
   * Merges with any data already in the file so that each test file's
   * afterAll hook accumulates into a single shared record (vitest isolates
   * modules per file even in singleFork mode, giving each file its own
   * Stats instance).
   */
  persist(): void {
    // Load previously written data from earlier test files in this run
    const accumulated: Record<string, number[]> = {};
    if (fs.existsSync(STATS_TEMP_FILE)) {
      try {
        const prior = JSON.parse(
          fs.readFileSync(STATS_TEMP_FILE, "utf8"),
        ) as Record<string, number[]>;
        for (const [k, v] of Object.entries(prior)) {
          accumulated[k] = v;
        }
      } catch {
        // Corrupt file — start fresh
      }
    }

    // Merge this file's data on top
    for (const [k, v] of this._data) {
      const existing = accumulated[k];
      accumulated[k] = existing ? [...existing, ...v] : v;
    }

    fs.writeFileSync(STATS_TEMP_FILE, JSON.stringify(accumulated), "utf8");
  }

  /** Load timing data from the temp file (used by the reporter) */
  static load(): EndpointStats[] {
    if (!fs.existsSync(STATS_TEMP_FILE)) return [];
    try {
      const raw = JSON.parse(
        fs.readFileSync(STATS_TEMP_FILE, "utf8"),
      ) as Record<string, number[]>;
      const map = new Map(Object.entries(raw));
      return new Stats().summarise(map);
    } catch {
      return [];
    }
  }

  /** Reset all data (used between runs if needed) */
  reset(): void {
    this._data.clear();
  }
}

function percentile(sorted: number[], p: number): number {
  if (sorted.length === 0) return 0;
  const idx = Math.floor((p / 100) * (sorted.length - 1));
  return sorted[Math.max(0, Math.min(idx, sorted.length - 1))] ?? 0;
}

/** Global singleton — imported by client.ts and setup.ts */
export const stats = new Stats();
