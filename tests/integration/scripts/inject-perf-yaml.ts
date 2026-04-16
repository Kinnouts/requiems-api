/**
 * inject-perf-yaml.ts
 *
 * Reads tests/integration/perf-baseline.json and upserts a top-level
 * `performance:` block into each matching api_docs YAML file in
 * apps/dashboard/config/api_docs/.
 *
 * Usage:
 *   pnpm run docs:perf:yaml
 *
 * The script is idempotent — running it twice produces the same result.
 */

import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import yaml from "js-yaml";
import { endpointYamlMap } from "./endpoint-yaml-map.js";

// ── Paths ─────────────────────────────────────────────────────────────────────

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const INTEGRATION_ROOT = path.join(__dirname, "..");
const PROJECT_ROOT = path.join(__dirname, "../../..");
const BASELINE_PATH = path.join(INTEGRATION_ROOT, "perf-baseline.json");
const API_DOCS_DIR = path.join(
  PROJECT_ROOT,
  "apps/dashboard/config/api_docs",
);

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
  generatedAt: string | null;
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

if (!baseline.generatedAt || baseline.endpoints.length === 0) {
  console.error(
    "❌  perf-baseline.json has no data.\n   Run: pnpm run benchmark",
  );
  process.exit(1);
}

const date = baseline.generatedAt.slice(0, 10); // YYYY-MM-DD

// ── Group by api_id, keep entry with most samples ─────────────────────────────

const apiIdBest = new Map<string, EndpointEntry>();

for (const entry of baseline.endpoints) {
  const apiId = endpointYamlMap[entry.path];
  if (!apiId) continue;

  const existing = apiIdBest.get(apiId);
  if (!existing || entry.samples > existing.samples) {
    apiIdBest.set(apiId, entry);
  }
}

// ── Inject performance block into each YAML file ──────────────────────────────

let updated = 0;
let skipped = 0;

for (const [apiId, entry] of apiIdBest) {
  const yamlPath = path.join(API_DOCS_DIR, `${apiId}.yml`);

  if (!fs.existsSync(yamlPath)) {
    console.warn(`⚠️  YAML not found, skipping: ${apiId}.yml`);
    skipped++;
    continue;
  }

  const raw = fs.readFileSync(yamlPath, "utf8");
  const doc = yaml.load(raw) as Record<string, unknown>;

  const newPerf = {
    p50_ms: entry.p50,
    p95_ms: entry.p95,
    p99_ms: entry.p99,
    avg_ms: entry.avg,
    samples: entry.samples,
    measured_at: date,
  };

  // Skip if data is identical (idempotency)
  const existing = doc["performance"] as typeof newPerf | undefined;
  if (
    existing &&
    existing.p50_ms === newPerf.p50_ms &&
    existing.p95_ms === newPerf.p95_ms &&
    existing.p99_ms === newPerf.p99_ms &&
    existing.measured_at === newPerf.measured_at
  ) {
    skipped++;
    continue;
  }

  doc["performance"] = newPerf;

  fs.writeFileSync(
    yamlPath,
    yaml.dump(doc, { lineWidth: 120, quotingType: '"', forceQuotes: false }),
    "utf8",
  );

  console.log(`✅  ${apiId}.yml`);
  updated++;
}

console.log(
  `\nDone — ${updated} YAML file(s) updated, ${skipped} unchanged/skipped.`,
);
