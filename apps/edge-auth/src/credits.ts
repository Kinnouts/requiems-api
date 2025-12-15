import type { WorkerBindings, CreditCheckResult } from "./types";

/**
 * Get current credit usage from D1
 *
 * For free tier: queries today's usage (resets at midnight UTC)
 * For paid tier: queries current billing month's usage
 */
export async function getUsage(
	bindings: WorkerBindings,
	apiKey: string,
	period: "daily" | "monthly",
	billingCycleStart?: string,
): Promise<number> {
	const startDate =
		period === "daily" ? getTodayStart() : billingCycleStart || getMonthStart();

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
 * Record credit usage in D1
 */
export async function recordUsage(
	bindings: WorkerBindings,
	apiKey: string,
	endpoint: string,
	credits: number,
): Promise<void> {
	await bindings.DB.prepare(`
    INSERT INTO credit_usage (api_key, endpoint, credits_used, used_at)
    VALUES (?, ?, ?, datetime('now'))
  `)
		.bind(apiKey, endpoint, credits)
		.run();
}

/**
 * Check credits and get current status
 */
export async function checkCredits(
	bindings: WorkerBindings,
	apiKey: string,
	period: "daily" | "monthly",
	limit: number,
	billingCycleStart?: string,
): Promise<CreditCheckResult> {
	const usage = await getUsage(bindings, apiKey, period, billingCycleStart);
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
 * Free tier credits reset here
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
 * Get when credits will reset
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
