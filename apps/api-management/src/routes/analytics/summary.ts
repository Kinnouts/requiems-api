import { Hono } from "hono";
import { jsonError, jsonResponse } from "../../shared/http";
import { createLogger } from "../../shared/logger";
import type { WorkerBindings } from "../../shared/types";
import type { EndpointStats, UsageSummary } from "./types";

const app = new Hono<{ Bindings: WorkerBindings }>();

/**
 * GET /analytics/summary
 * Get overall usage summary for a user
 *
 * Query parameters:
 * - userId: string (required)
 * - since: ISO timestamp (optional, defaults to billing cycle start)
 * - until: ISO timestamp (optional, defaults to now)
 */
app.get("/summary", async (c) => {
  const log = createLogger(c.req.raw);

  const userId = c.req.query("userId");
  const since = c.req.query("since");
  const until = c.req.query("until") || new Date().toISOString();

  if (!userId) {
    return jsonError(400, "Missing required parameter: userId");
  }

  try {
    // If no "since" provided, get the earliest billing cycle start
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

    // Get total requests and credits
    const totalsResult = await c.env.DB.prepare(`
      SELECT
        COUNT(*) as totalRequests,
        SUM(credits_used) as totalCredits
      FROM credit_usage
      WHERE user_id = ?
        AND used_at >= ?
        AND used_at <= ?
    `)
      .bind(userId, sinceDate, until)
      .first<{ totalRequests: number; totalCredits: number }>();

    // Get top 5 endpoints
    const topEndpointsResult = await c.env.DB.prepare(`
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
      LIMIT 5
    `)
      .bind(userId, sinceDate, until)
      .all<EndpointStats>();

    const summary: UsageSummary = {
      userId,
      totalRequests: totalsResult?.totalRequests || 0,
      totalCredits: totalsResult?.totalCredits || 0,
      dateRange: {
        since: sinceDate,
        until,
      },
      topEndpoints: topEndpointsResult.results || [],
    };

    log.info("Analytics summary fetched", { userId });

    return jsonResponse(summary);
  } catch (error) {
    log.error("Error fetching analytics summary", {
      error,
      params: { userId, since, until },
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
