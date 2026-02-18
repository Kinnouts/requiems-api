import { Hono } from "hono";
import { jsonError, jsonResponse } from "../../shared/http";
import { createLogger } from "../../shared/logger";
import type { WorkerBindings } from "../../shared/types";
import type { UsageExportResponse, UsageRecord } from "./types";

const app = new Hono<{ Bindings: WorkerBindings }>();

/**
 * GET /usage/export
 * Export usage data from D1 for sync to PostgreSQL
 *
 * Query parameters:
 * - since: ISO timestamp - get records after this time (required)
 * - limit: number - max records to return (default: 1000, max: 5000)
 * - cursor: string - pagination cursor (offset)
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

  // Parse cursor (offset)
  const offset = Number.parseInt(cursorParam || "0", 10);
  if (Number.isNaN(offset) || offset < 0) {
    return jsonError(400, "Invalid cursor parameter");
  }

  try {
    // Get total count for hasMore calculation
    const countResult = await c.env.DB.prepare(`
      SELECT COUNT(*) as total
      FROM credit_usage
      WHERE used_at >= ?
    `)
      .bind(since)
      .first<{ total: number }>();

    const total = countResult?.total || 0;

    // Fetch usage records with pagination
    const result = await c.env.DB.prepare(`
      SELECT
        api_key,
        endpoint,
        credits_used,
        used_at
      FROM credit_usage
      WHERE used_at >= ?
      ORDER BY used_at ASC
      LIMIT ? OFFSET ?
    `)
      .bind(since, limit, offset)
      .all<UsageRecord>();

    const usage = result.results || [];
    const hasMore = offset + usage.length < total;
    const nextCursor = hasMore ? (offset + usage.length).toString() : undefined;

    log.info("Usage export successful", {
      total,
      returned: usage.length,
      hasMore,
    });

    const response: UsageExportResponse = {
      usage,
      total,
      hasMore,
      nextCursor,
    };

    return jsonResponse(response);
  } catch (error) {
    log.error("Error exporting usage data", {
      error,
      params: { since, limit, offset },
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
