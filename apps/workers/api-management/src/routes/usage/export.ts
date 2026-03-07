import { Hono } from "hono";
import { jsonError, jsonResponse } from "@requiem/workers-shared";
import { createLogger } from "@requiem/workers-shared";
import type { WorkerBindings } from "../../env";
import type { UsageExportResponse, UsageRecord } from "./types";

const app = new Hono<{ Bindings: WorkerBindings }>();

/**
 * GET /usage/export
 * Export usage data from D1 for sync to PostgreSQL
 *
 * Query parameters:
 * - since: ISO timestamp - get records after this time (required)
 * - limit: number - max records to return (default: 1000, max: 5000)
 * - cursor: string - pagination cursor (last seen record id)
 *
 * Auth: X-API-Management-Key header (only Rails dashboard has this)
 */
app.get("/export", async (c) => {
  const log = createLogger(c.req.raw);

  const since = c.req.query("since");
  const limitParam = c.req.query("limit");
  const cursorParam = c.req.query("cursor");

  // Validate required parameter
  if (!since) {
    return jsonError(400, "Missing required parameter: since");
  }

  // Validate ISO timestamp
  const sinceDate = new Date(since);
  if (Number.isNaN(sinceDate.getTime())) {
    return jsonError(400, "Invalid timestamp format for 'since' parameter");
  }

  // Parse and validate limit
  const limit = Math.min(
    Number.parseInt(limitParam || "1000", 10),
    5000, // Max 5000 records per request
  );

  if (Number.isNaN(limit) || limit < 1) {
    return jsonError(400, "Invalid limit parameter");
  }

  // Parse cursor (last seen id — 0 means start from beginning)
  const afterId = Number.parseInt(cursorParam || "0", 10);
  
  if (Number.isNaN(afterId) || afterId < 0) {
    return jsonError(400, "Invalid cursor parameter");
  }

  try {
    // Fetch usage records using keyset pagination on the autoincrement id.
    // This is stable under concurrent inserts (unlike OFFSET which can skip rows).
    const result = await c.env.DB.prepare(`
      SELECT
        id,
        api_key,
        endpoint,
        credits_used,
        used_at
      FROM credit_usage
      WHERE used_at >= ? AND id > ?
      ORDER BY id ASC
      LIMIT ?
    `)
      .bind(since, afterId, limit)
      .all<UsageRecord>();

    const records = result.results || [];
    const hasMore = records.length === limit;
    const nextCursor = hasMore ? records[records.length - 1].id.toString() : undefined;

    log.info("Usage export successful", {
      returned: records.length,
      hasMore,
    });

    // Strip internal id field before returning to callers
    const usage = records.map(({ id: _id, ...rest }) => rest);

    const response: UsageExportResponse = {
      usage,
      hasMore,
      nextCursor,
    };

    return jsonResponse(response);
  } catch (error) {
    log.error("Error exporting usage data", {
      error,
      params: { since, limit, afterId },
    });

    // Return more detailed error in development
    if (c.env.ENVIRONMENT === "development") {
      return jsonError(
        500,
        `Failed to export usage data: ${error instanceof Error ? error.message : String(error)}`,
      );
    }

    return jsonError(500, "Failed to export usage data");
  }
});

export default app;
