/**
 * Lightweight snapshot utilities.
 *
 * Rather than comparing raw response bodies (which change every call for
 * random-data endpoints), we compare the *shape* of the response — i.e. the
 * set of top-level keys and the JSON types of their values.
 *
 * Snapshots are written to tests/integration/snapshots/ as JSON files the
 * first time a suite is run.  On subsequent runs they are loaded and compared
 * so you can detect regressions or unintentional shape changes.
 */

import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const SNAPSHOT_DIR = path.join(path.dirname(fileURLToPath(import.meta.url)), "../../snapshots");

/** Maps a JSON value to a human-readable type token */
function typeOf(value: unknown): string {
  if (value === null) return "null";
  if (Array.isArray(value)) {
    const elementTypes = new Set((value as unknown[]).map(typeOf));
    const inner = [...elementTypes].sort().join("|");
    return `Array<${inner || "unknown"}>`;
  }
  if (typeof value === "object") {
    return shapeOf(value as Record<string, unknown>);
  }
  return typeof value;
}

/** Recursively build a shape descriptor for a plain object */
function shapeOf(obj: Record<string, unknown>): string {
  const entries = Object.entries(obj)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([k, v]) => `${k}:${typeOf(v)}`);
  return `{${entries.join(",")}}`;
}

/**
 * Derive the shape of an arbitrary JSON value.
 * Objects are represented as `{key:type,...}`, arrays as `Array<type>`.
 */
export function deriveShape(value: unknown): string {
  return typeOf(value);
}

/** Derive the shape of every key in the top-level object independently */
export function deriveTopLevelShapes(obj: Record<string, unknown>): Record<string, string> {
  const result: Record<string, string> = {};
  for (const [k, v] of Object.entries(obj)) {
    result[k] = typeOf(v);
  }
  return result;
}

export interface SnapshotFile {
  /** ISO timestamp of when the snapshot was first created */
  createdAt: string;
  /** ISO timestamp of the last update */
  updatedAt: string;
  /** Map of endpoint path → recorded shape */
  shapes: Record<string, string>;
}

function snapshotPath(name: string): string {
  return path.join(SNAPSHOT_DIR, `${name}.snap.json`);
}

/** Load an existing snapshot file, or return undefined if none exists */
export function loadSnapshot(name: string): SnapshotFile | undefined {
  const p = snapshotPath(name);
  if (!fs.existsSync(p)) return undefined;
  return JSON.parse(fs.readFileSync(p, "utf8")) as SnapshotFile;
}

/** Write (or update) a snapshot file */
export function saveSnapshot(name: string, shapes: Record<string, string>): void {
  fs.mkdirSync(SNAPSHOT_DIR, { recursive: true });

  const existing = loadSnapshot(name);
  const now = new Date().toISOString();

  const file: SnapshotFile = {
    createdAt: existing?.createdAt ?? now,
    updatedAt: now,
    shapes: { ...(existing?.shapes ?? {}), ...shapes },
  };

  fs.writeFileSync(snapshotPath(name), JSON.stringify(file, null, 2) + "\n", "utf8");
}

/**
 * Assert that the shape of `body` for the given `endpointKey` matches the
 * saved snapshot (if one exists). On first run the shape is recorded.
 *
 * @param snapshotName  Logical name of the snapshot file (e.g. "text")
 * @param endpointKey   Identifies the specific endpoint within that file
 * @param body          Parsed JSON response body
 */
export function assertShape(
  snapshotName: string,
  endpointKey: string,
  body: Record<string, unknown>,
): void {
  const shape = deriveShape(body);
  const existing = loadSnapshot(snapshotName);

  if (!existing || !(endpointKey in existing.shapes)) {
    // First run — record the shape
    const shapes = existing?.shapes ?? {};
    shapes[endpointKey] = shape;
    saveSnapshot(snapshotName, shapes);
    return;
  }

  const recorded = existing.shapes[endpointKey];
  if (shape !== recorded) {
    throw new Error(
      [
        `Response shape mismatch for "${endpointKey}":`,
        `  Expected : ${recorded}`,
        `  Received : ${shape}`,
        "",
        "If this change is intentional, delete the snapshot and re-run.",
      ].join("\n"),
    );
  }
}
