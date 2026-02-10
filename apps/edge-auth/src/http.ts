import { env } from "./env";

const CORS_HEADERS = {
  "Access-Control-Allow-Origin": "*",
};

export const corsResponse = new Response(null, {
  headers: {
    ...CORS_HEADERS,
    "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
    "Access-Control-Allow-Headers": "Content-Type, requiems-api-key",
    "Access-Control-Max-Age": "86400",
  },
});

/**
 * JSON response helper
 */
export function jsonResponse(
  data: unknown,
  status = 200,
  headers: Record<string, string> = {},
): Response {
  return new Response(JSON.stringify(data), {
    status,
    headers: {
      "Content-Type": "application/json",
      ...CORS_HEADERS,
      ...headers,
    },
  });
}

/**
 * JSON error response helper
 */
export function jsonError(
  status: number,
  message: string,
  headers: Record<string, string> = {},
): Response {
  return jsonResponse({ error: message }, status, headers);
}

/**
 * Filter headers before forwarding to backend
 * Removes Cloudflare headers and sensitive data
 */
export function filterHeaders(headers: Headers): Headers {
  const filtered = new Headers();

  for (const [key, value] of headers.entries()) {
    const lowerKey = key.toLowerCase();

    if (lowerKey.startsWith("cf-")) continue;
    if (lowerKey === "requiems-api-key") continue;
    if (lowerKey === "connection") continue;
    if (lowerKey === "keep-alive") continue;

    filtered.set(key, value);
  }

  filtered.set("X-Backend-Secret", env.BACKEND_SECRET);

  return filtered;
}

/**
 * Add credit and rate limit headers to response
 */
export function addUsageHeaders(
  response: Response,
  headers: {
    creditsUsed: number;
    creditsRemaining: number;
    creditsReset: string;
    plan: string;
    rateLimitLimit: number;
    rateLimitRemaining: number;
  },
): Response {
  const newResponse = new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: response.headers,
  });

  newResponse.headers.set("Access-Control-Allow-Origin", "*");
  newResponse.headers.set("X-Credits-Used", headers.creditsUsed.toString());
  newResponse.headers.set(
    "X-Credits-Remaining",
    headers.creditsRemaining.toString(),
  );
  newResponse.headers.set("X-Credits-Reset", headers.creditsReset);
  newResponse.headers.set("X-Plan", headers.plan);
  newResponse.headers.set(
    "X-RateLimit-Limit",
    headers.rateLimitLimit.toString(),
  );
  newResponse.headers.set(
    "X-RateLimit-Remaining",
    headers.rateLimitRemaining.toString(),
  );

  return newResponse;
}

export type BackendResult =
  | { ok: true; response: Response }
  | { ok: false; error: string };

/**
 * Fetch from backend with error handling
 */
export async function fetchBackend(
  url: string | URL,
  init: RequestInit,
): Promise<BackendResult> {
  try {
    const response = await fetch(url, init);
    return { ok: true, response };
  } catch (error) {
    console.error("Backend error:", error);
    return { ok: false, error: "Backend unavailable" };
  }
}
