import type { RequestCheckResult, WorkerBindings } from "./types";

/**
 * Get current request usage from D1
 *
 * Note: Database tables still use "credit_" naming for historical reasons,
 * but we treat these as request counts in the code.
 */
export async function getRequestUsage(
  bindings: WorkerBindings,
  apiKey: string,
  period: "daily" | "monthly",
  billingCycleStart?: string,
): Promise<number> {
  const startDate = period === "daily"
    ? getTodayStart()
    : billingCycleStart || getMonthStart();

  const result = await bindings.DB.prepare(`
    SELECT COALESCE(SUM(credits_used), 0) as total
    FROM credit_usage
    WHERE api_key = ? AND used_at >= ?
  `)
    .bind(apiKey, startDate)
    .first<{ total: number }>();

  return result?.total || 0;
}

/**
 * Record request usage in D1
 */
export async function recordRequestUsage(
  bindings: WorkerBindings,
  apiKey: string,
  endpoint: string,
  requests: number,
): Promise<void> {
  await bindings.DB.prepare(`
    INSERT INTO credit_usage (api_key, endpoint, credits_used, used_at)
    VALUES (?, ?, ?, datetime('now'))
  `)
    .bind(apiKey, endpoint, requests)
    .run();
}

/**
 * Check request usage and get current status
 */
export async function checkRequestUsage(
  bindings: WorkerBindings,
  apiKey: string,
  period: "daily" | "monthly",
  limit: number,
  billingCycleStart?: string,
): Promise<RequestCheckResult> {
  const usage = await getRequestUsage(bindings, apiKey, period, billingCycleStart);
  const remaining = Math.max(0, limit - usage);
  const resetAt = getResetTime(period, billingCycleStart);

  return {
    usage,
    remaining,
    limit,
    resetAt,
  };
}

/**
 * Get start of today (midnight UTC)
 * Used for daily reset periods
 */
export function getTodayStart(): string {
  const now = new Date();
  now.setUTCHours(0, 0, 0, 0);
  return now.toISOString();
}

/**
 * Get start of current month (1st at midnight UTC)
 * Default billing cycle start for paid users
 */
export function getMonthStart(): string {
  const now = new Date();
  now.setUTCDate(1);
  now.setUTCHours(0, 0, 0, 0);
  return now.toISOString();
}

/**
 * Get when request quota will reset
 */
export function getResetTime(
  period: "daily" | "monthly",
  billingCycleStart?: string,
): string {
  const now = new Date();

  if (period === "daily") {
    // Tomorrow at midnight UTC
    now.setUTCDate(now.getUTCDate() + 1);
    now.setUTCHours(0, 0, 0, 0);
    return now.toISOString();
  }

  // Monthly: next billing cycle
  if (billingCycleStart) {
    // Calculate next billing date based on cycle start
    const cycleStart = new Date(billingCycleStart);
    const dayOfMonth = cycleStart.getUTCDate();

    const nextReset = new Date(now);
    nextReset.setUTCDate(dayOfMonth);
    nextReset.setUTCHours(0, 0, 0, 0);

    // If we're past this month's reset date, go to next month
    if (nextReset <= now) {
      nextReset.setUTCMonth(nextReset.getUTCMonth() + 1);
    }

    return nextReset.toISOString();
  }

  // Default: first of next month
  now.setUTCMonth(now.getUTCMonth() + 1);
  now.setUTCDate(1);
  now.setUTCHours(0, 0, 0, 0);
  return now.toISOString();
}
