import type { Context, Next } from "hono";
import { basicAuth } from "hono/basic-auth";
import type { WorkerBindings } from "../shared/types";

/**
 * Swagger documentation authentication middleware
 * Protects /docs/* routes with HTTP Basic Auth in production
 */
export async function swaggerAuthMiddleware(c: Context<{ Bindings: WorkerBindings }>, next: Next) {
	const env = c.env;

	if (env.ENVIRONMENT === "production") {
		const username = env.SWAGGER_USERNAME;
		const password = env.SWAGGER_PASSWORD;

		if (!username || !password) {
			console.error("Swagger docs access denied: authentication not configured in production");
			return c.json({ error: "Documentation access requires authentication configuration" }, 403);
		}

		const auth = basicAuth({
			username,
			password,
		});

		return await auth(c, next);
	}

	await next();
}
