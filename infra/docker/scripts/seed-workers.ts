/**
 * Seeds local wrangler KV and D1 for the full-stack dev environment.
 * Must be run from the auth-gateway working directory (/workers/auth-gateway)
 * so wrangler.toml bindings and schema.sql are resolved correctly.
 */

import { execSync } from "node:child_process";

interface KeyEntry {
  apiKey: string;
  userId: string;
  plan: string;
}

const DEV_KEYS: KeyEntry[] = [
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
run("bunx wrangler d1 execute requiem-usage --local --yes --file=./schema.sql");

console.log("\nSeeding KV with dev test keys...");
for (const { apiKey, userId, plan } of DEV_KEYS) {
  const kvKey = `key:${apiKey}`;
  const value = JSON.stringify({ userId, plan, createdAt: CREATED_AT });
  run(`bunx wrangler kv key put '${kvKey}' '${value}' --binding=KV --local`);
}

console.log("\nDev test keys seeded (header: requiems-api-key):");
for (const { apiKey, plan } of DEV_KEYS) {
  console.log(`  ${plan.padEnd(12)} -> ${apiKey}`);
}
console.log(
  `\nExample: curl -H 'requiems-api-key: ${DEV_KEYS[0].apiKey}' http://localhost:6000/v1/text/advice\n`,
);
