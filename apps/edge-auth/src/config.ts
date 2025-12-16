import type { PlanConfig, PlanName } from "./types";

/**
 * Plan configurations
 *
 * - Free: 50 credits/day, hard limit, resets daily at midnight UTC
 * - Paid: Monthly pool, use anytime, overage billed at end of month
 */
export const PLANS: Record<PlanName, PlanConfig> = {
  free: {
    creditLimit: 50,
    creditPeriod: "daily",
    ratePerSecond: 1,
    ratePerMinute: 30,
    overageRate: null, // Hard limit - no overage
  },
  starter: {
    creditLimit: 30_000,
    creditPeriod: "monthly",
    ratePerSecond: 10,
    ratePerMinute: 300,
    overageRate: 0.001,
  },
  pro: {
    creditLimit: 150_000,
    creditPeriod: "monthly",
    ratePerSecond: 30,
    ratePerMinute: 1000,
    overageRate: 0.0008,
  },
  business: {
    creditLimit: 500_000,
    creditPeriod: "monthly",
    ratePerSecond: 100,
    ratePerMinute: 5000,
    overageRate: 0.0005,
  },
  enterprise: {
    creditLimit: Infinity,
    creditPeriod: "monthly",
    ratePerSecond: 1000,
    ratePerMinute: 50000,
    overageRate: 0,
  },
};

/**
 * Endpoint costs - credit cost per API call
 *
 * IMPORTANT: Keep this in sync with your Go backend routes!
 * When you add a new endpoint to the backend, add it here too.
 */
export const ENDPOINT_COSTS: Record<string, number> = {
  // ==========================================================================
  // Text Domain (/v1/text/*)
  // ==========================================================================
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
