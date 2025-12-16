import type {
  PlanConfig,
  PlanName,
  RateLimitResult,
  WorkerBindings,
} from "./types";

/**
 * Get credit limits description for a plan
 */
export function getPlanLimits(plan: PlanName): string {
  const limits: Record<PlanName, string> = {
    free: "50 credits/day",
    starter: "30k credits/month",
    pro: "150k credits/month",
    business: "500k credits/month",
    enterprise: "custom limits",
  };

  return limits[plan];
}

export async function checkRateLimit(
  bindings: WorkerBindings,
  apiKey: string,
  plan: PlanConfig,
): Promise<RateLimitResult> {
  const now = Date.now();
  const currentMinute = Math.floor(now / 60000);

  const minuteKey = `rl:m:${apiKey}:${currentMinute}`;
  const existing = (await bindings.KV.get(minuteKey)) ?? "0";

  const minuteCount = Number.parseInt(existing, 10);
  const resetAt = (currentMinute + 1) * 60000;

  if (minuteCount >= plan.ratePerMinute) {
    return {
      allowed: false,
      remaining: 0,
      resetAt,
    };
  }

  await bindings.KV.put(minuteKey, (minuteCount + 1).toString(), {
    expirationTtl: 60,
  });

  return {
    allowed: true,
    remaining: plan.ratePerMinute - minuteCount - 1,
    resetAt,
  };
}
