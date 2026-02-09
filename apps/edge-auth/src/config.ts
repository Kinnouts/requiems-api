import type { PlanConfig, PlanName } from "./types";

/**
 * Plan configurations
 *
 * - Free: 50 credits/day, hard limit, resets daily at midnight UTC
 * - Paid: Monthly pool, use anytime, hard limit (no overage)
 */
export const PLANS: Record<PlanName, PlanConfig> = {
  free: {
    creditLimit: 50,
    creditPeriod: "daily",
    ratePerMinute: 30,
  },
  developer: {
    creditLimit: 500_000,
    creditPeriod: "monthly",
    ratePerMinute: 5000,
  },
  business: {
    creditLimit: 500_000,
    creditPeriod: "monthly",
    ratePerMinute: 5000,
  },
  professional: {
    creditLimit: Infinity,
    creditPeriod: "monthly",
    ratePerMinute: 50000,
  },
};

/**
 * Endpoint costs - credit cost per API call
 *
 * IMPORTANT: Keep this in sync with your Go backend routes!
 * When you add a new endpoint to the backend, add it here too.
 */
export const ENDPOINT_COSTS: Record<string, number> = {
  "GET /v1/text/advice": 1,
  "GET /v1/text/quotes/random": 1,
  "GET /v1/text/words/random": 1,
  "GET /v1/text/words/define": 2,
  "GET /v1/text/words/synonyms": 2,
};

/**
 * Default cost for unknown endpoints
 * Backend will return 404, but we still need a cost for the attempt
 */
export const DEFAULT_ENDPOINT_COST = 1;

/**
 * Get the credit cost for an endpoint
 * Matches exact path first, then tries prefix matching for dynamic routes
 */
export function getEndpointCost(method: string, pathname: string): number {
  const exactKey = `${method} ${pathname}`;

  if (ENDPOINT_COSTS[exactKey] !== undefined) {
    return ENDPOINT_COSTS[exactKey];
  }

  // Try prefix matching (for routes like /v1/finance/stocks/:symbol)
  for (const [route, cost] of Object.entries(ENDPOINT_COSTS)) {
    const [routeMethod, routePath] = route.split(" ");

    if (method === routeMethod && pathname.startsWith(routePath)) {
      return cost;
    }
  }

  return DEFAULT_ENDPOINT_COST;
}
