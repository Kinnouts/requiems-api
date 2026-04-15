/**
 * Custom Vitest reporter that prints a per-endpoint performance table at the
 * end of the run and persists the results as a JSON report.
 *
 * Timing data is read from the temp file written by the test worker process
 * (via stats.persist() in the global setup teardown).
 *
 * Vitest's reporter interface:
 *   https://vitest.dev/advanced/reporters.html
 */

import type { File, Reporter } from "vitest";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { Stats } from "./stats.js";

const REPORT_DIR = path.join(
  path.dirname(fileURLToPath(import.meta.url)),
  "../../reports",
);

export default class PerformanceReporter implements Reporter {
  onFinished(_files?: File[]): void {
    const summary = Stats.load();
    if (summary.length === 0) return;

    // Pretty-print table to stdout
    const colWidths = {
      path: Math.max(4, ...summary.map((s) => s.path.length)),
      n: 7,
      min: 7,
      avg: 7,
      p50: 7,
      p95: 7,
      p99: 7,
      max: 7,
    };

    const header = [
      "Path".padEnd(colWidths.path),
      "Samples".padStart(colWidths.n),
      "Min(ms)".padStart(colWidths.min),
      "Avg(ms)".padStart(colWidths.avg),
      "P50(ms)".padStart(colWidths.p50),
      "P95(ms)".padStart(colWidths.p95),
      "P99(ms)".padStart(colWidths.p99),
      "Max(ms)".padStart(colWidths.max),
    ].join("  ");

    const sep = "-".repeat(header.length);

    console.log("\n\n📊  Integration Test Performance Summary");
    console.log(sep);
    console.log(header);
    console.log(sep);

    for (const s of summary) {
      const row = [
        s.path.padEnd(colWidths.path),
        String(s.samples).padStart(colWidths.n),
        String(s.min).padStart(colWidths.min),
        String(s.avg).padStart(colWidths.avg),
        String(s.p50).padStart(colWidths.p50),
        String(s.p95).padStart(colWidths.p95),
        String(s.p99).padStart(colWidths.p99),
        String(s.max).padStart(colWidths.max),
      ].join("  ");
      console.log(row);
    }
    console.log(sep);
    console.log();

    // Write JSON report
    try {
      fs.mkdirSync(REPORT_DIR, { recursive: true });
      const reportPath = path.join(
        REPORT_DIR,
        `perf-${new Date().toISOString().replace(/[:.]/g, "-")}.json`,
      );
      fs.writeFileSync(
        reportPath,
        JSON.stringify(
          {
            generatedAt: new Date().toISOString(),
            baseUrl: process.env["API_BASE_URL"] ?? "https://api.requiems.xyz",
            endpoints: summary,
          },
          null,
          2,
        ) + "\n",
        "utf8",
      );
      console.log(`📁  Performance report saved to: ${reportPath}`);
    } catch {
      // Non-fatal — report printing to stdout is more important
    }
  }
}
