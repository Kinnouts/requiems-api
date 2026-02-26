/**
 * Seeds local wrangler KV and D1 for the full-stack dev environment.
 * Run from the auth-gateway working directory so wrangler.toml bindings
 * and schema.sql resolve correctly.
 *
 * Usage: node ./scripts/seed-dev.ts
 */

import { execSync } from "node:child_process";

interface DevKey {
  apiKey: string;
  userId: string;
  plan: string;
}

const WRANGLER = "pnpm exec wrangler";

const DEV_KEYS: DevKey[] = [
  { apiKey: "rq_test_0000000000000001", userId: "dev_user_free", plan: "free" },
  { apiKey: "rq_test_0000000000000002", userId: "dev_user_developer", plan: "developer" },
  { apiKey: "rq_test_0000000000000003", userId: "dev_user_business", plan: "business" },
  { apiKey: "rq_test_0000000000000004", userId: "dev_user_professional", plan: "professional" },
];

const CREATED_AT = "2025-01-01T00:00:00.000Z";

function run(cmd: string): void {
  execSync(cmd, { stdio: "inherit" });
}

console.log("Applying D1 schema...");
run(`${WRANGLER} d1 execute requiem-usage --local --yes --file=./schema.sql`);

console.log("\nSeeding KV with dev test keys...");
for (const { apiKey, userId, plan } of DEV_KEYS) {
  const kvKey = `key:${apiKey}`;
  const value = JSON.stringify({ userId, plan, createdAt: CREATED_AT });
  run(`${WRANGLER} kv key put '${kvKey}' '${value}' --binding=KV --local`);
}

console.log("\nDev test keys seeded (header: requiems-api-key):");
for (const { apiKey, plan } of DEV_KEYS) {
  console.log(`  ${plan.padEnd(12)} -> ${apiKey}`);
}
console.log(
  `\nExample: curl -H 'requiems-api-key: ${DEV_KEYS[0].apiKey}' http://localhost:4455/v1/text/advice\n`,
);
