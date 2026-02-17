import type { Context, Next } from "hono";
import { basicAuth } from "hono/basic-auth";
import type { WorkerBindings } from "../shared/types";

/**
 * Swagger documentation authentication middleware
 * Protects /docs/* routes with HTTP Basic Auth in production
 */
export async function swaggerAuthMiddleware(c: Context<{ Bindings: WorkerBindings }>, next: Next) {
	const env = c.env;

	if (env.ENVIRONMENT === "production" && env.SWAGGER_USERNAME && env.SWAGGER_PASSWORD) {
		const auth = basicAuth({
			username: env.SWAGGER_USERNAME,
			password: env.SWAGGER_PASSWORD,
		});

		return auth(c, next);
	}

	await next();
}
