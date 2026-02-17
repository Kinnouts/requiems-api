import { Hono } from "hono";
import { jsonError, jsonResponse, requireBackendSecret } from "../shared/http";
import { createLogger } from "../shared/logger";
import type { WorkerBindings } from "../shared/types";

const app = new Hono<{ Bindings: WorkerBindings }>();

/**
 * Endpoint usage stats
 */
export interface EndpointStats {
  endpoint: string;
  requests: number;
  credits: number;
}

/**
 * Date-based usage stats
 */
export interface DateStats {
  date: string;
  requests: number;
  credits: number;
}

/**
 * Usage summary for a user
 */
export interface UsageSummary {
  userId: string;
  totalRequests: number;
  totalCredits: number;
  dateRange: {
    since: string;
    until: string;
  };
  topEndpoints: EndpointStats[];
}

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

  // Check authentication
  const authError = requireBackendSecret(c.req.raw, c.env.BACKEND_SECRET);
  if (authError) {
    log.warn("Unauthorized analytics request");
    return authError;
  }

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

      sinceDate = billingResult?.earliest || new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();
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
    log.error("Error fetching endpoint analytics", { error });
    return jsonError(500, "Failed to fetch analytics");
  }
});

/**
 * GET /analytics/by-date
 * Get usage trends over time for a user
 *
 * Query parameters:
 * - userId: string (required)
 * - since: ISO timestamp (optional, defaults to 30 days ago)
 * - until: ISO timestamp (optional, defaults to now)
 * - groupBy: "day" | "hour" (optional, defaults to "day")
 */
app.get("/by-date", async (c) => {
  const log = createLogger(c.req.raw);

  // Check authentication
  const authError = requireBackendSecret(c.req.raw, c.env.BACKEND_SECRET);
  if (authError) {
    log.warn("Unauthorized analytics request");
    return authError;
  }

  const userId = c.req.query("userId");
  const until = c.req.query("until") || new Date().toISOString();
  const since = c.req.query("since") || new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();
  const groupBy = c.req.query("groupBy") || "day";

  if (!userId) {
    return jsonError(400, "Missing required parameter: userId");
  }

  if (groupBy !== "day" && groupBy !== "hour") {
    return jsonError(400, "Invalid groupBy parameter. Must be 'day' or 'hour'");
  }

  try {
    // SQLite date formatting: strftime for day or hour grouping
    const dateFormat = groupBy === "day" ? "%Y-%m-%d" : "%Y-%m-%d %H:00:00";

    const result = await c.env.DB.prepare(`
      SELECT
        strftime('${dateFormat}', used_at) as date,
        COUNT(*) as requests,
        SUM(credits_used) as credits
      FROM credit_usage
      WHERE user_id = ?
        AND used_at >= ?
        AND used_at <= ?
      GROUP BY date
      ORDER BY date ASC
    `)
      .bind(userId, since, until)
      .all<DateStats>();

    log.info("Analytics by date fetched", {
      userId,
      dataPoints: result.results?.length || 0,
      groupBy,
    });

    return jsonResponse({
      timeSeries: result.results || [],
      dateRange: { since, until },
      groupBy,
    });
  } catch (error) {
    log.error("Error fetching date analytics", { error });
    return jsonError(500, "Failed to fetch analytics");
  }
});

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

  // Check authentication
  const authError = requireBackendSecret(c.req.raw, c.env.BACKEND_SECRET);
  if (authError) {
    log.warn("Unauthorized analytics request");
    return authError;
  }

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

      sinceDate = billingResult?.earliest || new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();
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
    log.error("Error fetching analytics summary", { error });
    return jsonError(500, "Failed to fetch analytics");
  }
});

export default app;
