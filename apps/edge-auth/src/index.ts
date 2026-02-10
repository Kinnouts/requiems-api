import { getRequestMultiplier, PLANS } from "./config";
import { checkRequestUsage, recordRequestUsage } from "./requests";
import { validateEnv } from "./env";
import {
  addUsageHeaders,
  corsResponse,
  fetchBackend,
  filterHeaders,
  jsonError,
  jsonResponse,
} from "./http";
import { createLogger, maskApiKey } from "./logger";
import { checkRateLimit, getRequestLimitMessage } from "./rate-limit";
import { handleUsageExport } from "./usage-export";

import type { ApiKeyData, WorkerBindings } from "./types";

async function fetch(
  request: Request,
  bindings: WorkerBindings,
): Promise<Response> {
  // Validate environment variables
  const env = validateEnv(bindings);
  const url = new URL(request.url);
  const pathname = url.pathname;

  if (pathname === "/healthz") {
    return jsonResponse({ status: "ok" });
  }

  if (pathname === "/internal/usage/export") {
    return handleUsageExport(request, bindings);
  }

  const log = createLogger(request);

  if (request.method === "OPTIONS") {
    return corsResponse;
  }

  const apiKey = request.headers.get("requiems-api-key");

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

  const requestUsage = await checkRequestUsage(
    bindings,
    keyData.userId,
    "monthly",
    plan.requestLimit,
    keyData.billingCycleStart,
  );

  const requestMultiplier = getRequestMultiplier(request.method, pathname);

  if (requestUsage.usage >= plan.requestLimit) {
    log.info("Request limit exceeded", {
      key: maskApiKey(apiKey),
      plan: keyData.plan,
      usage: requestUsage.usage,
    });

    return jsonError(429, getRequestLimitMessage(), {
      "X-Requests-Used": "0",
      "X-Requests-Remaining": "0",
      "X-Requests-Reset": requestUsage.resetAt,
      "X-Plan": keyData.plan,
    });
  }

  const backendUrl = new URL(pathname + url.search, env.BACKEND_URL);

  const backendHeaders = filterHeaders(request.headers, env.BACKEND_SECRET);

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
      requestsUsed: 0,
      requestsRemaining: requestUsage.remaining,
      requestsReset: requestUsage.resetAt,
      plan: keyData.plan,
      rateLimitLimit: plan.ratePerMinute,
      rateLimitRemaining: rateLimit.remaining,
    });

    return response;
  }

  void recordRequestUsage(bindings, apiKey, keyData.userId, pathname, requestMultiplier);

  const response = addUsageHeaders(backendResponse, {
    requestsUsed: requestMultiplier,
    requestsRemaining: Math.max(0, requestUsage.remaining - requestMultiplier),
    requestsReset: requestUsage.resetAt,
    plan: keyData.plan,
    rateLimitLimit: plan.ratePerMinute,
    rateLimitRemaining: rateLimit.remaining,
  });

  return response;
}

export default { fetch };

export type { ApiKeyData, WorkerBindings } from "./types";
