import type { PlanConfig, PlanName } from "./types";

/**
 * Plan configurations
 *
 * All plans use monthly request quotas
 * Quotas reset at the start of each billing cycle (monthly)
 *
 * Note: Some endpoints count as multiple requests (see ENDPOINT_MULTIPLIERS)
 *
 * Pricing:
 * - Free: $0/month - 500 requests
 * - Developer: $29/month - 100,000 requests
 * - Business: $74/month - 1,000,000 requests
 * - Professional: $149/month - 10,000,000 requests
 * - Enterprise: Custom - Unlimited requests
 */
export const PLANS: Record<PlanName, PlanConfig> = {
  free: {
    requestLimit: 500,
    ratePerMinute: 30,
  },
  developer: {
    requestLimit: 100_000,
    ratePerMinute: 5000,
  },
  business: {
    requestLimit: 1_000_000,
    ratePerMinute: 10000,
  },
  professional: {
    requestLimit: 10_000_000,
    ratePerMinute: 50000,
  },
  enterprise: {
    requestLimit: Infinity,
    ratePerMinute: Infinity,
  },
};

/**
 * Request multipliers - how many requests an endpoint counts as
 *
 * DEFAULT: 1 request per API call
 * This map ONLY lists endpoints that count as MORE than 1 request.
 *
 * Most endpoints = 1 request (default)
 * Some endpoints require expensive operations = 2x+ requests (listed below)
 *
 * Examples of expensive operations:
 * - API calls to external services (2-3x)
 * - Complex AI/ML inference (3-5x)
 * - Large data processing (2-3x)
 *
 * IMPORTANT: Keep this in sync with your Go backend routes!
 * When adding expensive endpoints, add them here.
 */
export const ENDPOINT_MULTIPLIERS = new Map<string, number>([
  // Dictionary operations count as 2 requests
  ["GET /v1/text/words/define", 2],
  ["GET /v1/text/words/synonyms", 2],

  // Future expensive endpoints:
  // ["GET /v1/ai/image-recognition", 5],
  // ["POST /v1/text/translate", 3],
]);

/**
 * Pre-computed prefix patterns for efficient matching
 * Format: [method, pathPrefix, multiplier]
 */
const ENDPOINT_PREFIXES: Array<[string, string, number]> = Array.from(
  ENDPOINT_MULTIPLIERS.entries(),
).map(([route, multiplier]) => {
  const [method, path] = route.split(" ", 2);
  return [method, path, multiplier];
});

/**
 * Default multiplier for endpoints not listed above
 * This is the multiplier for 90%+ of endpoints
 */
export const DEFAULT_REQUEST_MULTIPLIER = 1;

/**
 * Get the request multiplier for an endpoint
 * Matches exact path first, then tries prefix matching for dynamic routes
 */
export function getRequestMultiplier(method: string, pathname: string): number {
  const exactKey = `${method} ${pathname}`;

  const exactMatch = ENDPOINT_MULTIPLIERS.get(exactKey);
  if (exactMatch !== undefined) {
    return exactMatch;
  }

   // Prefix matching for dynamic routes (e.g., /v1/finance/stocks/:symbol)
  for (const [routeMethod, routePath, multiplier] of ENDPOINT_PREFIXES) {
    if (method === routeMethod && pathname.startsWith(routePath)) {
      return multiplier;
    }
  }

  return DEFAULT_REQUEST_MULTIPLIER;
}
