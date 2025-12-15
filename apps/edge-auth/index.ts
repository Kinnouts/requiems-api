/**
 * Requiem API Gateway
 *
 * This Cloudflare Worker handles:
 * - API key validation
 * - Credit checking & tracking
 * - Rate limiting
 * - Request forwarding to backend
 */

export interface Env {
  // Backend URL (secret, internal)
  BACKEND_URL: string;

  // Cloudflare KV namespace for API keys & config
  KV: KVNamespace;

  // Cloudflare D1 database for usage tracking
  DB: D1Database;
}

// ============================================================================
// TYPES
// ============================================================================

interface ApiKeyData {
  userId: string;
  plan: "free" | "starter" | "pro" | "business" | "enterprise";
  createdAt: string;
}

interface PlanConfig {
  creditLimit: number;
  creditPeriod: "daily" | "monthly";
  ratePerSecond: number;
  ratePerMinute: number;
  overageRate: number | null; // null = hard limit (free tier)
}

// Endpoint costs - keep in sync with your API routes
const ENDPOINT_COSTS: Record<string, number> = {
  // Text domain
  "GET /v1/text/advice": 1,
  "GET /v1/text/quotes/random": 1,
  "GET /v1/text/words/random": 1,
  "GET /v1/text/words/define": 2,
  "GET /v1/text/words/synonyms": 2,

  // Finance domain (coming soon)
  "GET /v1/finance/commodities": 3,
  "GET /v1/finance/stocks": 3,
  "GET /v1/finance/crypto": 3,
  "GET /v1/finance/exchange-rates": 2,

  // Places domain (coming soon)
  "GET /v1/places/geocode": 5,
  "GET /v1/places/reverse-geocode": 5,
  "GET /v1/places/timezone": 2,
};

// Plan configurations
const PLANS: Record<string, PlanConfig> = {
  free: {
    creditLimit: 50,
    creditPeriod: "daily",
    ratePerSecond: 1,
    ratePerMinute: 30,
    overageRate: null,
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

// ============================================================================
// MAIN HANDLER
// ============================================================================

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);
    const pathname = url.pathname;

    // Health check endpoint (no auth required)
    if (pathname === "/healthz") {
      return jsonResponse({ status: "ok" });
    }

    // ========================================================================
    // 1. VALIDATE API KEY
    // ========================================================================
    const apiKey = request.headers.get("x-api-key");

    if (!apiKey) {
      return jsonError(401, "Missing x-api-key header");
    }

    // Lookup key in KV
    const keyData = await env.KV.get<ApiKeyData>(`key:${apiKey}`, "json");

    if (!keyData) {
      return jsonError(401, "Invalid API key");
    }

    const plan = PLANS[keyData.plan];

    if (!plan) {
      return jsonError(500, "Invalid plan configuration");
    }

    // ========================================================================
    // 2. RATE LIMITING (simple in-memory style using KV)
    // ========================================================================
    const rateLimitKey = `ratelimit:${apiKey}:${
      Math.floor(Date.now() / 60000)
    }`; // per minute bucket
    const currentRequests = parseInt((await env.KV.get(rateLimitKey)) || "0");

    if (currentRequests >= plan.ratePerMinute) {
      return jsonError(429, "Rate limit exceeded", {
        "X-RateLimit-Limit": plan.ratePerMinute.toString(),
        "X-RateLimit-Remaining": "0",
        "Retry-After": "60",
      });
    }

    // Increment rate limit counter
    await env.KV.put(rateLimitKey, (currentRequests + 1).toString(), {
      expirationTtl: 120,
    });

    // ========================================================================
    // 3. CHECK CREDITS
    // ========================================================================
    const usage = await getUsage(env, apiKey, plan.creditPeriod);
    const endpointCost = getEndpointCost(request.method, pathname);

    // Free tier: hard limit
    if (plan.overageRate === null && usage >= plan.creditLimit) {
      return jsonError(
        429,
        "Credit limit exceeded. Upgrade your plan for more credits.",
        {
          "X-Credits-Used": endpointCost.toString(),
          "X-Credits-Remaining": "0",
          "X-Plan": keyData.plan,
        },
      );
    }

    // ========================================================================
    // 4. FORWARD TO BACKEND
    // ========================================================================
    const backendUrl = new URL(pathname + url.search, env.BACKEND_URL);

    const backendResponse = await fetch(backendUrl.toString(), {
      method: request.method,
      headers: filterHeaders(request.headers),
      body: request.body,
    });

    // ========================================================================
    // 5. RECORD USAGE (only for successful responses)
    // ========================================================================
    if (backendResponse.ok && endpointCost > 0) {
      await recordUsage(env, apiKey, pathname, endpointCost);
    }

    // ========================================================================
    // 6. ADD HEADERS AND RETURN
    // ========================================================================
    const newUsage = backendResponse.ok ? usage + endpointCost : usage;
    const remaining = Math.max(0, plan.creditLimit - newUsage);
    const resetTime = getResetTime(plan.creditPeriod);

    // Clone response and add headers
    const response = new Response(backendResponse.body, {
      status: backendResponse.status,
      statusText: backendResponse.statusText,
      headers: backendResponse.headers,
    });

    response.headers.set("X-Credits-Used", endpointCost.toString());
    response.headers.set("X-Credits-Remaining", remaining.toString());
    response.headers.set("X-Credits-Reset", resetTime);
    response.headers.set("X-Plan", keyData.plan);
    response.headers.set("X-RateLimit-Limit", plan.ratePerMinute.toString());
    response.headers.set(
      "X-RateLimit-Remaining",
      Math.max(0, plan.ratePerMinute - currentRequests - 1).toString(),
    );

    return response;
  },
};

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

/**
 * Get endpoint cost from the config
 * Matches exact path first, then tries prefix matching for dynamic routes
 */
function getEndpointCost(method: string, pathname: string): number {
  // Try exact match first
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

  // Default cost for unknown endpoints (still forward, backend will 404)
  return 1;
}

/**
 * Get current usage from D1 database
 */
async function getUsage(
  env: Env,
  apiKey: string,
  period: "daily" | "monthly",
): Promise<number> {
  const startDate = period === "daily" ? getTodayStart() : getMonthStart();

  const result = await env.DB.prepare(`
    SELECT COALESCE(SUM(credits_used), 0) as total
    FROM credit_usage
    WHERE api_key = ? AND used_at >= ?
  `)
    .bind(apiKey, startDate)
    .first<{ total: number }>();

  return result?.total || 0;
}

/**
 * Record usage in D1 database
 */
async function recordUsage(
  env: Env,
  apiKey: string,
  endpoint: string,
  credits: number,
): Promise<void> {
  await env.DB.prepare(`
    INSERT INTO credit_usage (api_key, endpoint, credits_used, used_at)
    VALUES (?, ?, ?, datetime('now'))
  `)
    .bind(apiKey, endpoint, credits)
    .run();
}

/**
 * Get start of today (UTC)
 */
function getTodayStart(): string {
  const now = new Date();
  now.setUTCHours(0, 0, 0, 0);
  return now.toISOString();
}

/**
 * Get start of current month (UTC)
 */
function getMonthStart(): string {
  const now = new Date();
  now.setUTCDate(1);
  now.setUTCHours(0, 0, 0, 0);
  return now.toISOString();
}

/**
 * Get reset time for credit period
 */
function getResetTime(period: "daily" | "monthly"): string {
  const now = new Date();

  if (period === "daily") {
    // Tomorrow midnight UTC
    now.setUTCDate(now.getUTCDate() + 1);
    now.setUTCHours(0, 0, 0, 0);
  } else {
    // First of next month
    now.setUTCMonth(now.getUTCMonth() + 1);
    now.setUTCDate(1);
    now.setUTCHours(0, 0, 0, 0);
  }

  return now.toISOString();
}

/**
 * Filter headers before forwarding to backend
 */
function filterHeaders(headers: Headers): Headers {
  const filtered = new Headers();

  for (const [key, value] of headers.entries()) {
    // Skip Cloudflare-specific headers
    if (key.toLowerCase().startsWith("cf-")) continue;
    if (key.toLowerCase() === "x-api-key") continue; // Don't forward API key

    filtered.set(key, value);
  }

  return filtered;
}

/**
 * JSON response helper
 */
function jsonResponse(
  data: unknown,
  status = 200,
  headers: Record<string, string> = {},
): Response {
  return new Response(JSON.stringify(data), {
    status,
    headers: {
      "Content-Type": "application/json",
      ...headers,
    },
  });
}

/**
 * JSON error response helper
 */
function jsonError(
  status: number,
  message: string,
  headers: Record<string, string> = {},
): Response {
  return jsonResponse({ error: message }, status, headers);
}
