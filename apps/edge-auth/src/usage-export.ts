import { jsonError, jsonResponse } from "./http";
import type { WorkerBindings } from "./types";

/**
 * Usage data record from D1
 */
export interface UsageRecord {
  api_key: string;
  endpoint: string;
  credits_used: number;
  used_at: string;
}

/**
 * Usage export response with pagination
 */
export interface UsageExportResponse {
  usage: UsageRecord[];
  total: number;
  hasMore: boolean;
  nextCursor?: string;
}

/**
 * Export usage data from D1 for sync to PostgreSQL
 *
 * This endpoint is protected by X-Backend-Secret header
 * and should only be called by the Rails backend.
 *
 * Query parameters:
 * - since: ISO timestamp - get records after this time (required)
 * - limit: number - max records to return (default: 1000, max: 5000)
 * - cursor: string - pagination cursor (offset)
 */
export async function handleUsageExport(
  request: Request,
  bindings: WorkerBindings,
): Promise<Response> {
  // Verify backend secret
  const backendSecret = request.headers.get("X-Backend-Secret");

  if (!backendSecret || backendSecret !== bindings.BACKEND_SECRET) {
    return jsonError(401, "Unauthorized - Invalid backend secret");
  }

  const url = new URL(request.url);
  const since = url.searchParams.get("since");
  const limitParam = url.searchParams.get("limit");
  const cursorParam = url.searchParams.get("cursor");

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
    const countResult = await bindings.DB.prepare(`
      SELECT COUNT(*) as total
      FROM credit_usage
      WHERE used_at >= ?
    `)
      .bind(since)
      .first<{ total: number }>();

    const total = countResult?.total || 0;

    // Fetch usage records with pagination
    const result = await bindings.DB.prepare(`
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

    const response: UsageExportResponse = {
      usage,
      total,
      hasMore,
      nextCursor,
    };

    return jsonResponse(response);
  } catch (error) {
    console.error("Error exporting usage data:", error);
    return jsonError(500, "Failed to export usage data");
  }
}
