/**
 * inject-perf-docs.ts
 *
 * Reads tests/integration/perf-baseline.json and upserts a "## Performance"
 * section into each API markdown doc that has matching benchmark data.
 *
 * Usage:
 *   pnpm run docs:perf
 *
 * The script is idempotent — running it twice produces the same result.
 */

import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { endpointDocMap } from "./endpoint-doc-map.js";

// ── Paths ────────────────────────────────────────────────────────────────────

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const INTEGRATION_ROOT = path.join(__dirname, "..");
const PROJECT_ROOT = path.join(__dirname, "../../..");
const BASELINE_PATH = path.join(INTEGRATION_ROOT, "perf-baseline.json");

// ── Types ─────────────────────────────────────────────────────────────────────

interface EndpointEntry {
  path: string;
  samples: number;
  min: number;
  avg: number;
  p50: number;
  p95: number;
  p99: number;
  max: number;
}

interface PerfBaseline {
  generatedAt: string;
  baseUrl: string;
  runs: number;
  endpoints: EndpointEntry[];
}

// ── Load baseline ─────────────────────────────────────────────────────────────

if (!fs.existsSync(BASELINE_PATH)) {
  console.error(
    `❌  perf-baseline.json not found at ${BASELINE_PATH}\n` +
      "   Run: pnpm run benchmark",
  );
  process.exit(1);
}

const baseline = JSON.parse(
  fs.readFileSync(BASELINE_PATH, "utf8"),
) as PerfBaseline;

const date = baseline.generatedAt.slice(0, 10); // YYYY-MM-DD

// ── Group entries by doc file, keeping the one with the most samples ──────────

const docBest = new Map<string, EndpointEntry>();

for (const entry of baseline.endpoints) {
  const docRelPath = endpointDocMap[entry.path];
  if (!docRelPath) continue;

  const existing = docBest.get(docRelPath);
  if (!existing || entry.samples > existing.samples) {
    docBest.set(docRelPath, entry);
  }
}

// ── Build the ## Performance section markdown ─────────────────────────────────

function buildPerfSection(entry: EndpointEntry, measuredAt: string): string {
  return [
    "## Performance",
    "",
    `Measured against production (\`${baseline.baseUrl}\`) with ${entry.samples} samples.`,
    "",
    "| Metric  | Value        |",
    "|---------|--------------|",
    `| p50     | ${entry.p50} ms |`,
    `| p95     | ${entry.p95} ms |`,
    `| p99     | ${entry.p99} ms |`,
    `| Average | ${entry.avg} ms |`,
    "",
    `_Last updated: ${measuredAt}_`,
  ].join("\n");
}

// ── Upsert the section into a markdown file ───────────────────────────────────

const PERF_SECTION_RE =
  /^## Performance\n[\s\S]*?(?=\n^## |\n*$)/m;

function upsertPerfSection(content: string, section: string): string {
  if (PERF_SECTION_RE.test(content)) {
    return content.replace(PERF_SECTION_RE, section);
  }
  // Append — ensure single trailing newline before section
  return content.trimEnd() + "\n\n" + section + "\n";
}

// ── Process each doc file ─────────────────────────────────────────────────────

let updated = 0;
let skipped = 0;

for (const [docRelPath, entry] of docBest) {
  const absPath = path.join(PROJECT_ROOT, docRelPath);

  if (!fs.existsSync(absPath)) {
    console.warn(`⚠️  Doc not found, skipping: ${docRelPath}`);
    skipped++;
    continue;
  }

  const original = fs.readFileSync(absPath, "utf8");
  const section = buildPerfSection(entry, date);
  const updated_content = upsertPerfSection(original, section);

  if (updated_content === original) {
    // Section already up to date
    skipped++;
    continue;
  }

  fs.writeFileSync(absPath, updated_content, "utf8");
  console.log(`✅  ${docRelPath}`);
  updated++;
}

console.log(
  `\nDone — ${updated} doc(s) updated, ${skipped} unchanged/skipped.`,
);
