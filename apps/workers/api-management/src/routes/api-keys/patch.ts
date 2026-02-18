import { Hono } from "hono";
import { jsonError, jsonResponse, createLogger, type ApiKeyData } from "@requiem/workers-shared";
import type { WorkerBindings } from "../../shared/env";

const app = new Hono<{ Bindings: WorkerBindings }>();

/**
 * Request body for updating an API key
 */
interface UpdateApiKeyRequest {
  plan?: "free" | "developer" | "business" | "professional" | "enterprise";
  billingCycleStart?: string;
}

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

  // Parse request body
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
    // Find the full key in KV
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

    // Get existing key data
    const existingData = await c.env.KV.get<ApiKeyData>(keyName, "json");
    if (!existingData) {
      log.warn("API key data not found in KV", { keyPrefix });
      return jsonError(404, "API key not found");
    }

    // Update key data
    const updatedData: ApiKeyData = {
      ...existingData,
      ...(body.plan && { plan: body.plan }),
      ...(body.billingCycleStart && { billingCycleStart: body.billingCycleStart }),
    };

    // Write updated data to KV
    await c.env.KV.put(keyName, JSON.stringify(updatedData));

    // Update in D1
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
    log.error("Failed to update API key", {
      error,
      params: { keyPrefix, updates: body },
    });

    if (c.env.ENVIRONMENT === "development") {
      return jsonError(
        500,
        `Failed to update API key: ${error instanceof Error ? error.message : String(error)}`,
      );
    }

    return jsonError(500, "Failed to update API key");
  }
});

export default app;
