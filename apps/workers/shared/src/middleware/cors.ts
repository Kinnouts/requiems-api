import type { MiddlewareHandler } from "hono";
import { corsResponse } from "../http";

/**
 * CORS Middleware
 *
 * Handles OPTIONS preflight requests with CORS headers
 */
export const corsMiddleware: MiddlewareHandler = async (c, next) => {
	if (c.req.method === "OPTIONS") {
		return corsResponse;
	}

	await next();
};
