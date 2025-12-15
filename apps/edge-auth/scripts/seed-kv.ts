#!/usr/bin/env npx tsx
/**
 * Seed KV with test API keys for local development
 *
 * Usage: npm run kv:seed
 *
 * This uses wrangler CLI under the hood.
 * For production, use the Cloudflare dashboard or wrangler directly.
 */

import { execSync } from "node:child_process";
import { randomBytes } from "node:crypto";

type PlanName = "free" | "starter" | "pro" | "business";

interface ApiKeyData {
	userId: string;
	plan: PlanName;
	createdAt: string;
	billingCycleStart?: string;
}

interface TestKey {
	key: string;
	data: ApiKeyData;
}

/**
 * Generate a random API key
 * Format: rq_test_{16 random hex chars}
 */
function generateApiKey(): string {
	const random = randomBytes(8).toString("hex");
	return `rq_test_${random}`;
}

/**
 * Generate a random user ID
 * Format: user_{12 random hex chars}
 */
function generateUserId(): string {
	const random = randomBytes(6).toString("hex");
	return `user_${random}`;
}

const now = new Date().toISOString();

const TEST_KEYS: TestKey[] = [
	{
		key: generateApiKey(),
		data: {
			userId: generateUserId(),
			plan: "free",
			createdAt: now,
		},
	},
	{
		key: generateApiKey(),
		data: {
			userId: generateUserId(),
			plan: "starter",
			createdAt: now,
			billingCycleStart: now,
		},
	},
	{
		key: generateApiKey(),
		data: {
			userId: generateUserId(),
			plan: "pro",
			createdAt: now,
			billingCycleStart: now,
		},
	},
	{
		key: generateApiKey(),
		data: {
			userId: generateUserId(),
			plan: "business",
			createdAt: now,
			billingCycleStart: now,
		},
	},
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
	const limits =
		plan === "free"
			? "50 credits/day"
			: plan === "starter"
				? "30k credits/month"
				: plan === "pro"
					? "150k credits/month"
					: "500k credits/month";

	console.log(`   ${plan.padEnd(10)} → ${key} (${limits})`);
}

console.log(
	"\n💡 Example:\n   curl -H 'x-api-key: " +
		createdKeys[0].key +
		"' http://localhost:8787/v1/text/advice",
);
