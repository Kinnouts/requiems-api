import { jsonError } from "@requiem/workers-shared";

import type { Context, Next } from "hono";
import type { WorkerBindings } from "../env";

function timingSafeEqual(a: string, b: string): boolean {
  const enc = new TextEncoder();
  const aBytes = enc.encode(a);
  const bBytes = enc.encode(b);
  if (aBytes.byteLength !== bBytes.byteLength) return false;
  return crypto.subtle.timingSafeEqual(aBytes, bBytes);
}

/**
 * API Management key authentication middleware
 * Validates X-API-Management-Key header for all private endpoints
 * Only the Rails dashboard should have this key
 */
export async function apiKeyAuthMiddleware(c: Context<{ Bindings: WorkerBindings }>, next: Next) {
  const apiKey = c.req.header("X-API-Management-Key");
  const expectedKey = c.env.API_MANAGEMENT_API_KEY;

  if (!apiKey || !timingSafeEqual(apiKey, expectedKey)) {
    return jsonError(401, "Unauthorized - Invalid or missing API management key");
  }

  await next();
}
