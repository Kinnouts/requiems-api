import { Hono } from "hono";
import { jsonError, jsonResponse, requireApiManagementKey } from "../shared/http";
import { createLogger, maskApiKey } from "../shared/logger";
import type {
  ApiKeyData,
  ApiKeyManagementRequest,
  ApiKeyManagementResponse,
  WorkerBindings,
} from "../shared/types";

const app = new Hono<{ Bindings: WorkerBindings }>();

/**
 * POST /api-keys
 * Create, revoke, or update API keys
 * Auth: X-API-Management-Key header (only Rails dashboard has this)
 */
app.post("/", async (c) => {
  const log = createLogger(c.req.raw);

  // Check authentication
  const authError = requireApiManagementKey(c.req.raw, c.env.API_MANAGEMENT_API_KEY);
  if (authError) {
    log.warn("Unauthorized API key management request");
    return authError;
  }

  // Parse request body
  let body: ApiKeyManagementRequest;
  try {
    body = await c.req.json();
  } catch (error) {
    log.error("Invalid JSON in API key management request", { error });
    return jsonError(400, "Invalid JSON body");
  }

  // Validate required fields
  if (!body.action || !body.key || !body.userId || !body.plan) {
    log.error("Missing required fields in API key management request", { body });
    return jsonError(400, "Missing required fields: action, key, userId, plan");
  }

  log.info("API key management request", {
    action: body.action,
    userId: body.userId,
    plan: body.plan,
    keyPrefix: maskApiKey(body.key),
  });

  try {
    switch (body.action) {
      case "create":
        return await createApiKey(c.env, body, log);
      case "revoke":
        return await revokeApiKey(c.env, body, log);
      case "update":
        return await updateApiKey(c.env, body, log);
      default:
        return jsonError(400, `Invalid action: ${body.action}`);
    }
  } catch (error) {
    log.error("API key management error", { error, action: body.action });
    return jsonError(500, "Internal server error");
  }
});

/**
 * Create a new API key in KV and D1
 */
async function createApiKey(
  env: WorkerBindings,
  body: ApiKeyManagementRequest,
  log: ReturnType<typeof createLogger>,
): Promise<Response> {
  const keyName = `key:${body.key}`;

  // Check if key already exists
  const existing = await env.KV.get(keyName);
  if (existing) {
    log.warn("API key already exists", { userId: body.userId });
    return jsonError(409, "API key already exists");
  }

  // Prepare key data
  const keyData: ApiKeyData = {
    userId: body.userId,
    plan: body.plan,
    createdAt: new Date().toISOString(),
    billingCycleStart: body.billingCycleStart || new Date().toISOString(),
  };

  // Write to KV
  await env.KV.put(keyName, JSON.stringify(keyData));

  // Write to D1 for metadata/audit
  await env.DB.prepare(
    `INSERT INTO api_keys (key_prefix, user_id, plan, created_at, billing_cycle_start)
     VALUES (?, ?, ?, ?, ?)`,
  )
    .bind(
      body.key.substring(0, 12), // key_prefix
      body.userId,
      body.plan,
      keyData.createdAt,
      keyData.billingCycleStart,
    )
    .run();

  log.info("API key created successfully", {
    userId: body.userId,
    plan: body.plan,
  });

  const response: ApiKeyManagementResponse = {
    success: true,
    message: "API key created successfully",
  };

  return jsonResponse(response);
}

/**
 * Revoke an API key by deleting from KV and marking in D1
 */
async function revokeApiKey(
  env: WorkerBindings,
  body: ApiKeyManagementRequest,
  log: ReturnType<typeof createLogger>,
): Promise<Response> {
  const keyName = `key:${body.key}`;

  // Check if key exists
  const existing = await env.KV.get(keyName);
  if (!existing) {
    log.warn("API key not found for revocation", { userId: body.userId });
    return jsonError(404, "API key not found");
  }

  // Delete from KV (revokes access immediately)
  await env.KV.delete(keyName);

  // Mark as revoked in D1
  await env.DB.prepare(
    `UPDATE api_keys
     SET revoked_at = ?, active = 0
     WHERE key_prefix = ?`,
  )
    .bind(new Date().toISOString(), body.key.substring(0, 12))
    .run();

  log.info("API key revoked successfully", { userId: body.userId });

  const response: ApiKeyManagementResponse = {
    success: true,
    message: "API key revoked successfully",
  };

  return jsonResponse(response);
}

/**
 * Update an API key (plan change, billing cycle)
 */
async function updateApiKey(
  env: WorkerBindings,
  body: ApiKeyManagementRequest,
  log: ReturnType<typeof createLogger>,
): Promise<Response> {
  const keyName = `key:${body.key}`;

  // Get existing key data
  const existingData = await env.KV.get<ApiKeyData>(keyName, "json");
  if (!existingData) {
    log.warn("API key not found for update", { userId: body.userId });
    return jsonError(404, "API key not found");
  }

  // Update key data
  const updatedData: ApiKeyData = {
    ...existingData,
    plan: body.plan,
    billingCycleStart: body.billingCycleStart || existingData.billingCycleStart,
  };

  // Write updated data to KV
  await env.KV.put(keyName, JSON.stringify(updatedData));

  // Update in D1
  await env.DB.prepare(
    `UPDATE api_keys
     SET plan = ?, billing_cycle_start = ?, updated_at = ?
     WHERE key_prefix = ?`,
  )
    .bind(
      updatedData.plan,
      updatedData.billingCycleStart,
      new Date().toISOString(),
      body.key.substring(0, 12),
    )
    .run();

  log.info("API key updated successfully", {
    userId: body.userId,
    newPlan: body.plan,
  });

  const response: ApiKeyManagementResponse = {
    success: true,
    message: "API key updated successfully",
  };

  return jsonResponse(response);
}

export default app;
