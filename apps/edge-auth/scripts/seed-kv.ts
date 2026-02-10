import { execSync } from "node:child_process";
import { randomBytes } from "node:crypto";
import type { ApiKeyData, PlanName } from "../src/types";
import { getPlanLimits } from "../src/rate-limit";

interface TestKey {
  key: string;
  data: ApiKeyData;
}

function generateApiKey(): string {
  const random = randomBytes(8).toString("hex");
  return `rq_test_${random}`;
}

function generateUserId(): string {
  const random = randomBytes(6).toString("hex");
  return `user_${random}`;
}

function testKey(plan: PlanName) {
  return {
    key: generateApiKey(),
    data: {
      userId: generateUserId(),
      plan: plan,
      createdAt: now,
    },
  };
}

const now = new Date().toISOString();

const TEST_KEYS: TestKey[] = [
  testKey("free"),
  testKey("starter"),
  testKey("pro"),
  testKey("business"),
  testKey("enterprise"),
];

console.log("🌱 Seeding KV with test API keys...\n");

const createdKeys: { key: string; plan: PlanName }[] = [];

for (const { key, data } of TEST_KEYS) {
  const kvKey = `key:${key}`;
  const value = JSON.stringify(data);

  console.log(`  Adding ${key} (${data.plan} plan)`);

  try {
    execSync(
      `npx wrangler kv:key put --binding=KV "${kvKey}" '${value}' --local`,
      { stdio: "pipe" },
    );
    createdKeys.push({ key, plan: data.plan });
  } catch (error) {
    console.error(`  ❌ Failed to add ${kvKey}`);
    console.error(error);
  }
}

console.log("\n✅ Done! Test keys created:\n");

for (const { key, plan } of createdKeys) {
  console.log(`   ${plan.padEnd(10)} → ${key} (${getPlanLimits(plan)})`);
}

console.log(
  "\n💡 Example:\n   curl -H 'requiems-api-key: " +
    createdKeys[0].key +
    "' http://localhost:8080/v1/text/advice",
);
