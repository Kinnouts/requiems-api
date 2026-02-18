import { Hono } from "hono";
import { jsonError, jsonResponse } from "../../shared/http";
import { createLogger } from "../../shared/logger";
import type { WorkerBindings } from "../../shared/types";

const app = new Hono<{ Bindings: WorkerBindings }>();

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
    // Find the full key in KV by searching with prefix pattern
    const keys = await c.env.KV.list({ prefix: "key:" });

    let fullKey: string | null = null;
    for (const key of keys.keys) {
      if (key.name.substring(4, 16) === keyPrefix) {
        // "key:" is 4 chars, then 12 char prefix
        fullKey = key.name.substring(4); // Remove "key:" prefix
        break;
      }
    }

    if (!fullKey) {
      log.warn("API key not found for revocation", { keyPrefix });
      return jsonError(404, "API key not found");
    }

    const keyName = `key:${fullKey}`;

    // Delete from KV (revokes access immediately)
    await c.env.KV.delete(keyName);

    // Mark as revoked in D1
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
    log.error("Failed to revoke API key", {
      error,
      params: { keyPrefix },
    });

    if (c.env.ENVIRONMENT === "development") {
      return jsonError(500, `Failed to revoke API key: ${error instanceof Error ? error.message : String(error)}`);
    }

    return jsonError(500, "Failed to revoke API key");
  }
});

export default app;
