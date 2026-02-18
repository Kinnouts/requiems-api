import { Hono } from "hono";
import { jsonError, jsonResponse } from "../../shared/http";
import { createLogger } from "../../shared/logger";
import type { WorkerBindings } from "../../shared/types";
import type { DateStats } from "./types";

const app = new Hono<{ Bindings: WorkerBindings }>();

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

  const userId = c.req.query("userId");
  const until = c.req.query("until") || new Date().toISOString();
  const since =
    c.req.query("since") || new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();
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
    log.error("Error fetching date analytics", {
      error,
      params: { userId, since, until, groupBy },
    });

    if (c.env.ENVIRONMENT === "development") {
      return jsonError(500, `Failed to fetch analytics: ${error instanceof Error ? error.message : String(error)}`);
    }

    return jsonError(500, "Failed to fetch analytics");
  }
});

export default app;
