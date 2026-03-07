import { Hono } from "hono";
import { sValidator } from "@hono/standard-validator";
import * as z from "zod";
import { jsonError, jsonResponse, createLogger } from "@requiem/workers-shared";
import type { WorkerBindings } from "../../env";
import type { DateStats } from "./types";

const app = new Hono<{ Bindings: WorkerBindings }>();

const byDateQuerySchema = z.object({
  userId: z.string().min(1, "Missing required parameter: userId"),
  since: z.string().optional(),
  until: z.string().optional(),
  groupBy: z.enum(["day", "hour"]).default("day"),
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
app.get(
  "/by-date",
  sValidator("query", byDateQuerySchema, (result, _c) => {
    if (!result.success) {
      return jsonError(400, result.error[0]?.message ?? "Validation error");
    }
  }),
  async (c) => {
    const log = createLogger(c.req.raw);
    const { userId, groupBy } = c.req.valid("query");
    const since =
      c.req.valid("query").since ?? new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();
    const until = c.req.valid("query").until ?? new Date().toISOString();

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
        return jsonError(
          500,
          `Failed to fetch analytics: ${error instanceof Error ? error.message : String(error)}`,
        );
      }

      return jsonError(500, "Failed to fetch analytics");
    }
  },
);

export default app;
