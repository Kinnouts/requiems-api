import { Hono } from "hono";
import { jsonError, jsonResponse } from "@requiem/workers-shared";
import { createLogger } from "@requiem/workers-shared";
import type { WorkerBindings } from "../../env";
import type { EndpointStats } from "./types";

const app = new Hono<{ Bindings: WorkerBindings }>();

/**
 * GET /analytics/by-endpoint
 * Get usage breakdown by endpoint for a user
 *
 * Query parameters:
 * - userId: string (required)
 * - since: ISO timestamp (optional, defaults to billing cycle start)
 * - until: ISO timestamp (optional, defaults to now)
 * - limit: number (optional, max top endpoints to return, default: 10)
 */
app.get("/by-endpoint", async (c) => {
  const log = createLogger(c.req.raw);

  const userId = c.req.query("userId");
  const since = c.req.query("since");
  const until = c.req.query("until") || new Date().toISOString();
  const limit = Math.min(Number.parseInt(c.req.query("limit") || "10", 10), 100);

  if (!userId) {
    return jsonError(400, "Missing required parameter: userId");
  }

  try {
    // If no "since" provided, get the earliest billing cycle start for this user
    let sinceDate = since;
    if (!sinceDate) {
      const billingResult = await c.env.DB.prepare(`
        SELECT MIN(billing_cycle_start) as earliest
        FROM api_keys
        WHERE user_id = ? AND active = 1
      `)
        .bind(userId)
        .first<{ earliest: string }>();

      sinceDate =
        billingResult?.earliest || new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();
    }

    // Query usage by endpoint
    const result = await c.env.DB.prepare(`
      SELECT
        endpoint,
        COUNT(*) as requests,
        SUM(credits_used) as credits
      FROM credit_usage
      WHERE user_id = ?
        AND used_at >= ?
        AND used_at <= ?
      GROUP BY endpoint
      ORDER BY credits DESC
      LIMIT ?
    `)
      .bind(userId, sinceDate, until, limit)
      .all<EndpointStats>();

    log.info("Analytics by endpoint fetched", {
      userId,
      endpoints: result.results?.length || 0,
    });

    return jsonResponse({
      endpoints: result.results || [],
      dateRange: { since: sinceDate, until },
    });
  } catch (error) {
    log.error("Error fetching endpoint analytics", {
      error,
      params: { userId, since, until, limit },
    });

    if (c.env.ENVIRONMENT === "development") {
      return jsonError(
        500,
        `Failed to fetch analytics: ${error instanceof Error ? error.message : String(error)}`,
      );
    }

    return jsonError(500, "Failed to fetch analytics");
  }
});

export default app;
