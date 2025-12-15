/**
 * Requiem API Gateway
 *
 * This Cloudflare Worker is the public API (api.requiems-api.xyz).
 * It handles authentication, rate limiting, and credit tracking,
 * then forwards approved requests to the internal Go backend.
 *
 * The backend verifies BACKEND_SECRET on every request.
 */

import type { ApiKeyData, WorkerBindings } from "./types";
import { env } from "./env";
import { getEndpointCost, PLANS } from "./config";
import { checkRateLimit } from "./rate-limit";
import { checkCredits, recordUsage } from "./credits";
import {
  addUsageHeaders,
  filterHeaders,
  jsonError,
  jsonResponse,
} from "./http";

export default {
  async fetch(request: Request, bindings: WorkerBindings): Promise<Response> {
    const url = new URL(request.url);
    const pathname = url.pathname;

    // ========================================================================
    // PUBLIC ENDPOINTS (no auth required)
    // ========================================================================

    if (pathname === "/healthz") {
      return jsonResponse({ status: "ok" });
    }

    // CORS preflight
    if (request.method === "OPTIONS") {
      return new Response(null, {
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
          "Access-Control-Allow-Headers": "Content-Type, x-api-key",
          "Access-Control-Max-Age": "86400",
        },
      });
    }

    // ========================================================================
    // 1. VALIDATE API KEY
    // ========================================================================

    const apiKey = request.headers.get("x-api-key");

    if (!apiKey) {
      return jsonError(
        401,
        "Missing x-api-key header. Get your key at requiems-api.xyz",
      );
    }

    // Lookup key in KV
    const keyData = await bindings.KV.get<ApiKeyData>(`key:${apiKey}`, "json");

    if (!keyData) {
      return jsonError(401, "Invalid API key");
    }

    const plan = PLANS[keyData.plan];
    if (!plan) {
      return jsonError(500, "Invalid plan configuration");
    }

    // ========================================================================
    // 2. RATE LIMITING
    // ========================================================================

    const rateLimit = await checkRateLimit(bindings, apiKey, plan);

    if (!rateLimit.allowed) {
      return jsonError(429, "Rate limit exceeded", {
        "X-RateLimit-Limit": plan.ratePerMinute.toString(),
        "X-RateLimit-Remaining": "0",
        "X-RateLimit-Reset": Math.ceil(rateLimit.resetAt / 1000).toString(),
        "Retry-After": Math.ceil(
          (rateLimit.resetAt - Date.now()) / 1000,
        ).toString(),
      });
    }

    // ========================================================================
    // 3. CHECK CREDITS
    // ========================================================================

    const credits = await checkCredits(
      bindings,
      apiKey,
      plan.creditPeriod,
      plan.creditLimit,
      keyData.billingCycleStart,
    );

    const endpointCost = getEndpointCost(request.method, pathname);

    // Free tier: hard limit (reject if over)
    if (plan.overageRate === null && credits.usage >= plan.creditLimit) {
      return jsonError(
        429,
        "Daily credit limit exceeded. Upgrade at requiems-api.xyz",
        {
          "X-Credits-Used": "0",
          "X-Credits-Remaining": "0",
          "X-Credits-Reset": credits.resetAt,
          "X-Plan": keyData.plan,
        },
      );
    }

    // ========================================================================
    // 4. FORWARD TO BACKEND
    // ========================================================================

    const backendUrl = new URL(pathname + url.search, env.BACKEND_URL);

    // Prepare headers with backend secret
    const backendHeaders = filterHeaders(request.headers);
    backendHeaders.set("X-Backend-Secret", env.BACKEND_SECRET);

    let backendResponse: Response;
    try {
      backendResponse = await fetch(backendUrl.toString(), {
        method: request.method,
        headers: backendHeaders,
        body: request.body,
      });
    } catch (error) {
      console.error("Backend error:", error);
      return jsonError(502, "Backend unavailable");
    }

    // ========================================================================
    // 5. RECORD USAGE (only for successful responses)
    // ========================================================================

    if (backendResponse.ok && endpointCost > 0) {
      // Fire and forget - don't block response
      recordUsage(bindings, apiKey, pathname, endpointCost);
    }

    // ========================================================================
    // 6. ADD HEADERS AND RETURN
    // ========================================================================

    const creditsUsed = backendResponse.ok ? endpointCost : 0;
    const creditsRemaining = backendResponse.ok
      ? Math.max(0, plan.creditLimit - credits.usage - endpointCost)
      : credits.remaining;

    const response = addUsageHeaders(backendResponse, {
      creditsUsed,
      creditsRemaining,
      creditsReset: credits.resetAt,
      plan: keyData.plan,
      rateLimitLimit: plan.ratePerMinute,
      rateLimitRemaining: rateLimit.remaining,
    });

    // Add CORS headers
    response.headers.set("Access-Control-Allow-Origin", "*");

    return response;
  },
};

// Re-export types for external use
export type { ApiKeyData, WorkerBindings } from "./types";
