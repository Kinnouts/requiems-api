import { Hono } from "hono";
import { jsonError, jsonResponse } from "../shared/http";
import { createLogger } from "../shared/logger";
import { ApiKeyGenerator } from "../shared/api-key-generator";
import type { ApiKeyData, WorkerBindings } from "../shared/types";

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
 * Request body for updating an API key
 */
interface UpdateApiKeyRequest {
	plan?: "free" | "developer" | "business" | "professional" | "enterprise";
	billingCycleStart?: string;
}

/**
 * POST /api-keys
 * Create a new API key
 * Auth: X-API-Management-Key header (only Rails dashboard has this)
 */
app.post("/", async (c) => {
	const log = createLogger(c.req.raw);

	let body: CreateApiKeyRequest;
	try {
		body = await c.req.json();
	} catch (error) {
		log.error("Invalid JSON in create API key request", { error });
		return jsonError(400, "Invalid JSON body");
	}

	if (!body.userId || !body.plan || !body.name) {
		log.error("Missing required fields", { body });
		return jsonError(400, "Missing required fields: userId, plan, name");
	}

	try {
		const fullKey = ApiKeyGenerator.generate();
		const keyPrefix = ApiKeyGenerator.extractPrefix(fullKey);
		const keyName = `key:${fullKey}`;

		const existing = await c.env.KV.get(keyName);
		
    if (existing) {
			log.warn("Generated key already exists (collision)", { keyPrefix });
			return jsonError(409, "Key collision detected, please retry");
		}

		const now = new Date().toISOString();
		
    const keyData: ApiKeyData = {
			userId: body.userId,
			plan: body.plan,
			createdAt: now,
			billingCycleStart: body.billingCycleStart || now,
		};

		await c.env.KV.put(keyName, JSON.stringify(keyData));

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
			apiKey: fullKey,
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

/**
 * DELETE /api-keys/:keyPrefix
 * Revoke an API key by its prefix
 * Auth: X-API-Management-Key header (only Rails dashboard has this)
 */
app.delete("/:keyPrefix", async (c) => {
	const log = createLogger(c.req.raw);
	const keyPrefix = c.req.param("keyPrefix");

	if (!keyPrefix) {
		return jsonError(400, "Missing keyPrefix parameter");
	}

	try {
		const keys = await c.env.KV.list({ prefix: "key:" });

		let fullKey: string | null = null;

		for (const key of keys.keys) {
			if (key.name.substring(4, 16) === keyPrefix) {
				fullKey = key.name.substring(4);
				break;
			}
		}

		if (!fullKey) {
			log.warn("API key not found for revocation", { keyPrefix });
			return jsonError(404, "API key not found");
		}

		const keyName = `key:${fullKey}`;

		await c.env.KV.delete(keyName);

		await c.env.DB.prepare(
			`UPDATE api_keys
       SET revoked_at = ?, active = 0
       WHERE key_prefix = ?`,
		)
			.bind(new Date().toISOString(), keyPrefix)
			.run();

		log.info("API key revoked successfully", { keyPrefix });

		return jsonResponse({
			success: true,
			message: "API key revoked successfully",
			keyPrefix,
		});
	} catch (error) {
		log.error("Failed to revoke API key", { error, keyPrefix });
		return jsonError(500, "Failed to revoke API key");
	}
});

/**
 * PATCH /api-keys/:keyPrefix
 * Update an API key (plan change, billing cycle)
 * Auth: X-API-Management-Key header (only Rails dashboard has this)
 */
app.patch("/:keyPrefix", async (c) => {
	const log = createLogger(c.req.raw);
	const keyPrefix = c.req.param("keyPrefix");

	if (!keyPrefix) {
		return jsonError(400, "Missing keyPrefix parameter");
	}

	let body: UpdateApiKeyRequest;
	try {
		body = await c.req.json();
	} catch (error) {
		log.error("Invalid JSON in update API key request", { error });
		return jsonError(400, "Invalid JSON body");
	}

	// Validate at least one field to update
	if (!body.plan && !body.billingCycleStart) {
		return jsonError(400, "Must provide at least one field to update: plan, billingCycleStart");
	}

	try {
		const keys = await c.env.KV.list({ prefix: "key:" });

		let fullKey: string | null = null;

		for (const key of keys.keys) {
			if (key.name.substring(4, 16) === keyPrefix) {
				fullKey = key.name.substring(4);
				break;
			}
		}

		if (!fullKey) {
			log.warn("API key not found for update", { keyPrefix });
			return jsonError(404, "API key not found");
		}

		const keyName = `key:${fullKey}`;

		const existingData = await c.env.KV.get<ApiKeyData>(keyName, "json");
		if (!existingData) {
			log.warn("API key data not found in KV", { keyPrefix });
			return jsonError(404, "API key not found");
		}

		const updatedData: ApiKeyData = {
			...existingData,
			...(body.plan && { plan: body.plan }),
			...(body.billingCycleStart && { billingCycleStart: body.billingCycleStart }),
		};

		await c.env.KV.put(keyName, JSON.stringify(updatedData));

		const updates: string[] = [];
		const bindings: (string | number)[] = [];

		if (body.plan) {
			updates.push("plan = ?");
			bindings.push(body.plan);
		}
		if (body.billingCycleStart) {
			updates.push("billing_cycle_start = ?");
			bindings.push(body.billingCycleStart);
		}
		updates.push("updated_at = ?");
		bindings.push(new Date().toISOString());
		bindings.push(keyPrefix);

		await c.env.DB.prepare(
			`UPDATE api_keys
       SET ${updates.join(", ")}
       WHERE key_prefix = ?`,
		)
			.bind(...bindings)
			.run();

		log.info("API key updated successfully", {
			keyPrefix,
			updates: body,
		});

		return jsonResponse({
			success: true,
			message: "API key updated successfully",
			keyPrefix,
			plan: updatedData.plan,
			billingCycleStart: updatedData.billingCycleStart,
		});
	} catch (error) {
		log.error("Failed to update API key", { error, keyPrefix });
		return jsonError(500, "Failed to update API key");
	}
});

export default app;
