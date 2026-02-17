import { Hono } from "hono";
import { jsonError, jsonResponse } from "../../shared/http";
import { createLogger } from "../../shared/logger";
import { ApiKeyGenerator } from "../../shared/api-key-generator";
import type { ApiKeyData, WorkerBindings } from "../../shared/types";

const app = new Hono<{ Bindings: WorkerBindings }>();

/**
 * Request body for creating a new API key
 */
interface CreateApiKeyRequest {
	userId: string;
	plan: "free" | "developer" | "business" | "professional" | "enterprise";
	name: string;
	billingCycleStart?: string;
}

/**
 * Response when creating a new API key
 */
interface CreateApiKeyResponse {
	apiKey: string; // Full key (only returned once)
	keyPrefix: string;
	userId: string;
	plan: string;
	createdAt: string;
}

/**
 * POST /api-keys
 * Create a new API key
 * Auth: X-API-Management-Key header (only Rails dashboard has this)
 */
app.post("/", async (c) => {
	const log = createLogger(c.req.raw);

	// Parse request body
	let body: CreateApiKeyRequest;
	try {
		body = await c.req.json();
	} catch (error) {
		log.error("Invalid JSON in create API key request", { error });
		return jsonError(400, "Invalid JSON body");
	}

	// Validate required fields
	if (!body.userId || !body.plan || !body.name) {
		log.error("Missing required fields", { body });
		return jsonError(400, "Missing required fields: userId, plan, name");
	}

	try {
		// Generate new API key
		const fullKey = ApiKeyGenerator.generate();
		const keyPrefix = ApiKeyGenerator.extractPrefix(fullKey);
		const keyName = `key:${fullKey}`;

		// Check if key already exists (extremely unlikely with nanoid, but good practice)
		const existing = await c.env.KV.get(keyName);
		if (existing) {
			log.warn("Generated key already exists (collision)", { keyPrefix });
			return jsonError(409, "Key collision detected, please retry");
		}

		// Prepare key data for KV
		const now = new Date().toISOString();
		const keyData: ApiKeyData = {
			userId: body.userId,
			plan: body.plan,
			createdAt: now,
			billingCycleStart: body.billingCycleStart || now,
		};

		// Write to KV
		await c.env.KV.put(keyName, JSON.stringify(keyData));

		// Write metadata to D1 for audit trail
		await c.env.DB.prepare(
			`INSERT INTO api_keys (key_prefix, user_id, plan, created_at, billing_cycle_start, active)
       VALUES (?, ?, ?, ?, ?, 1)`,
		)
			.bind(keyPrefix, body.userId, body.plan, now, keyData.billingCycleStart)
			.run();

		log.info("API key created successfully", {
			userId: body.userId,
			plan: body.plan,
			keyPrefix,
		});

		const response: CreateApiKeyResponse = {
			apiKey: fullKey, // Return full key (Rails will store hash)
			keyPrefix,
			userId: body.userId,
			plan: body.plan,
			createdAt: now,
		};

		return jsonResponse(response, 201);
	} catch (error) {
		log.error("Failed to create API key", { error });
		return jsonError(500, "Failed to create API key");
	}
});

export default app;
