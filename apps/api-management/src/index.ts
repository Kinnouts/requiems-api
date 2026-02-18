import { Hono } from "hono";
import { swaggerUI } from "@hono/swagger-ui";

import { validateEnv } from "./shared/env";
import { jsonResponse } from "./shared/http";

import { apiKeyAuthMiddleware } from "./middleware/api-key-auth";

import apiKeysRoute from "./routes/api-keys";
import usageRoute from "./routes/usage";
import analyticsRoute from "./routes/analytics";
import swaggerRoute from "./routes/swagger";

import type { WorkerBindings } from "./shared/types";
import { basicAuth } from "hono/basic-auth";

const app = new Hono<{ Bindings: WorkerBindings }>();

app.get("/healthz", (_c) => jsonResponse({ status: "ok", service: "api-management" }));

app.use("/docs", async (c, next) => {
  console.log("Swagger UI access attempt:", {
    ip: c.req.header("CF-Connecting-IP") || c.req.header("X-Forwarded-For") || "unknown",
    userAgent: c.req.header("User-Agent") || "unknown",
  });

  await next();
});

app.use("/docs", basicAuth({
  verifyUser: (username, password, c) => {
    const adminUser = c.env.SWAGGER_USERNAME!
    const adminPass = c.env.SWAGGER_PASSWORD!

    return (
      username === adminUser && password === adminPass
    )
  },
}));

app.get("/docs", swaggerUI({ url: "/openapi.json" }));

app.route("/", swaggerRoute);

app.use("/api-keys/*", apiKeyAuthMiddleware);
app.use("/usage/*", apiKeyAuthMiddleware);
app.use("/analytics/*", apiKeyAuthMiddleware);

app.route("/api-keys", apiKeysRoute);
app.route("/usage", usageRoute);
app.route("/analytics", analyticsRoute);

app.notFound((_c) => {
  return jsonResponse({ error: "Not found" }, 404);
});

app.onError((err, c) => {
  console.error("Unhandled error:", {
    message: err.message,
    name: err.name,
    stack: err.stack,
  });

  if (c.env?.ENVIRONMENT === "development") {
    return jsonResponse({
      error: "Internal server error",
      details: err.message,
      name: err.name,
      stack: err.stack,
    }, 500);
  }

  return jsonResponse({
    error: "Internal server error",
    message: err.message
  }, 500);
});

export default {
  async fetch(request: Request, env: WorkerBindings): Promise<Response> {
    try {
      validateEnv(env);
    } catch (error) {
      console.error("Environment validation failed:", error);

      return new Response(JSON.stringify({
        error: "Configuration error",
        details: error instanceof Error ? error.message : String(error),
      }), {
        status: 500,
        headers: { "Content-Type": "application/json" },
      });
    }

    return app.fetch(request, env);
  },
};
