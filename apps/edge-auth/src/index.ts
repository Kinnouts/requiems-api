import { getEndpointCost, PLANS } from "./config";
import { checkCredits, recordUsage } from "./credits";
import { env } from "./env";
import {
  addUsageHeaders,
  corsResponse,
  fetchBackend,
  filterHeaders,
  jsonError,
  jsonResponse,
} from "./http";
import { createLogger, maskApiKey } from "./logger";
import { checkRateLimit, getCreditLimitMessage } from "./rate-limit";

import type { ApiKeyData, WorkerBindings } from "./types";

export default {
  async fetch(request: Request, bindings: WorkerBindings): Promise<Response> {
    const log = createLogger(request);
    const url = new URL(request.url);
    const pathname = url.pathname;

    if (pathname === "/healthz") {
      return jsonResponse({ status: "ok" });
    }

    if (request.method === "OPTIONS") {
      return corsResponse;
    }

    const apiKey = request.headers.get("x-api-key");

    if (!apiKey) {
      return jsonError(401, "Get your key at requiems-api.xyz");
    }

    const keyData = await bindings.KV.get<ApiKeyData>(`key:${apiKey}`, "json");

    if (!keyData) {
      return jsonError(401, "Invalid API key");
    }

    const plan = PLANS[keyData.plan];

    if (!plan) {
      return jsonError(500, "Invalid plan configuration");
    }

    const rateLimit = await checkRateLimit(bindings, apiKey, plan);

    if (!rateLimit.allowed) {
      log.info("Rate limit exceeded", {
        key: maskApiKey(apiKey),
        plan: keyData.plan,
      });

      return jsonError(429, "Rate limit exceeded", {
        "X-RateLimit-Limit": plan.ratePerMinute.toString(),
        "X-RateLimit-Remaining": "0",
        "X-RateLimit-Reset": Math.ceil(rateLimit.resetAt / 1000).toString(),
        "Retry-After": Math.ceil(
          (rateLimit.resetAt - Date.now()) / 1000,
        ).toString(),
      });
    }

    const credits = await checkCredits(
      bindings,
      apiKey,
      plan.creditPeriod,
      plan.creditLimit,
      keyData.billingCycleStart,
    );

    const endpointCost = getEndpointCost(request.method, pathname);

    if (credits.usage >= plan.creditLimit) {
      log.info("Credit limit exceeded", {
        key: maskApiKey(apiKey),
        plan: keyData.plan,
        usage: credits.usage,
      });

      return jsonError(429, getCreditLimitMessage(plan.creditPeriod), {
        "X-Credits-Used": "0",
        "X-Credits-Remaining": "0",
        "X-Credits-Reset": credits.resetAt,
        "X-Plan": keyData.plan,
      });
    }

    const backendUrl = new URL(pathname + url.search, env.BACKEND_URL);

    const backendHeaders = filterHeaders(request.headers);

    const result = await fetchBackend(backendUrl, {
      method: request.method,
      headers: backendHeaders,
      body: request.body,
    });

    if (!result.ok) {
      log.error("Backend fetch failed", { error: result.error });
      return jsonError(502, result.error);
    }

    const backendResponse = result.response;

    if (!backendResponse.ok) {
      log.warn("Backend error response", {
        status: backendResponse.status,
        path: pathname,
      });

      const response = addUsageHeaders(backendResponse, {
        creditsUsed: 0,
        creditsRemaining: credits.remaining,
        creditsReset: credits.resetAt,
        plan: keyData.plan,
        rateLimitLimit: plan.ratePerMinute,
        rateLimitRemaining: rateLimit.remaining,
      });

      return response;
    }

    void recordUsage(bindings, apiKey, pathname, endpointCost);

    const response = addUsageHeaders(backendResponse, {
      creditsUsed: endpointCost,
      creditsRemaining: Math.max(0, credits.remaining - endpointCost),
      creditsReset: credits.resetAt,
      plan: keyData.plan,
      rateLimitLimit: plan.ratePerMinute,
      rateLimitRemaining: rateLimit.remaining,
    });

    return response;
  },
};

export type { ApiKeyData, WorkerBindings } from "./types";
