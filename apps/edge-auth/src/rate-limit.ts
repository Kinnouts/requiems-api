import type { PlanConfig, RateLimitResult, WorkerBindings } from "./types";

/**
 * Sliding window rate limiter using KV
 *
 * Uses two windows (current second + current minute) to enforce both limits.
 * Keys auto-expire so no cleanup needed.
 */
export async function checkRateLimit(
  bindings: WorkerBindings,
  apiKey: string,
  plan: PlanConfig,
): Promise<RateLimitResult> {
  const now = Date.now();
  const currentSecond = Math.floor(now / 1000);
  const currentMinute = Math.floor(now / 60000);

  // Check per-second limit
  const secondKey = `rl:s:${apiKey}:${currentSecond}`;
  const secondCount = parseInt((await bindings.KV.get(secondKey)) || "0");

  if (secondCount >= plan.ratePerSecond) {
    return {
      allowed: false,
      remaining: 0,
      resetAt: (currentSecond + 1) * 1000, // Next second
    };
  }

  // Check per-minute limit
  const minuteKey = `rl:m:${apiKey}:${currentMinute}`;
  const minuteCount = parseInt((await bindings.KV.get(minuteKey)) || "0");

  if (minuteCount >= plan.ratePerMinute) {
    return {
      allowed: false,
      remaining: 0,
      resetAt: (currentMinute + 1) * 60000, // Next minute
    };
  }

  // Increment both counters (fire and forget for speed)
  await Promise.all([
    bindings.KV.put(secondKey, (secondCount + 1).toString(), {
      expirationTtl: 2,
    }),
    bindings.KV.put(minuteKey, (minuteCount + 1).toString(), {
      expirationTtl: 120,
    }),
  ]);

  return {
    allowed: true,
    remaining: plan.ratePerMinute - minuteCount - 1,
    resetAt: (currentMinute + 1) * 60000,
  };
}
