import type { Context, Next } from "hono";
import { jsonError } from "@requiem/workers-shared";
import type { WorkerBindings } from "../shared/env";

/**
 * API Management key authentication middleware
 * Validates X-API-Management-Key header for all private endpoints
 * Only the Rails dashboard should have this key
 */
export async function apiKeyAuthMiddleware(c: Context<{ Bindings: WorkerBindings }>, next: Next) {
  const apiKey = c.req.header("X-API-Management-Key");
  const expectedKey = c.env.API_MANAGEMENT_API_KEY;

  if (!apiKey || apiKey !== expectedKey) {
    return jsonError(401, "Unauthorized - Invalid or missing API management key");
  }

  await next();
}
