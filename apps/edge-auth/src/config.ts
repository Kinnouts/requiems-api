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
 * DEFAULT COST: 1 credit per request
 * This object ONLY lists endpoints that cost MORE than 1 credit.
 *
 * Most endpoints are simple lookups = 1 credit (default)
 * Some endpoints require expensive operations = 2+ credits (listed below)
 *
 * Examples of expensive operations:
 * - API calls to external services (2-3x cost)
 * - Complex AI/ML inference (3-5x cost)
 * - Large data processing (2-3x cost)
 *
 * IMPORTANT: Keep this in sync with your Go backend routes!
 * When adding expensive endpoints, add them here.
 */
export const ENDPOINT_COSTS: Record<string, number> = {
  // Dictionary operations are more expensive
  "GET /v1/text/words/define": 2,
  "GET /v1/text/words/synonyms": 2,

  // Future expensive endpoints:
  // "GET /v1/ai/image-recognition": 5,
  // "POST /v1/text/translate": 3,
};

/**
 * Default cost for endpoints not listed above
 * This is the cost for 90%+ of endpoints
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
