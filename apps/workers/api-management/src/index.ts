import { Hono } from "hono";
import { swaggerUI } from "@hono/swagger-ui";

import { validateEnv, type WorkerBindings } from "./env";
import { basicAuthMiddleware, errorHandler, jsonResponse } from "@requiem/workers-shared";

import { apiKeyAuthMiddleware } from "./middleware/";

import { apiKeysRoute, usageRoute, analyticsRoute, swaggerRoute } from "./routes";

const app = new Hono<{ Bindings: WorkerBindings }>();

app.get("/healthz", (_c) => jsonResponse({ status: "ok", service: "api-management" }));

app.use("/docs", basicAuthMiddleware);

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

app.onError(errorHandler);

export default {
  async fetch(request: Request, env: WorkerBindings): Promise<Response> {
    try {
      validateEnv(env);
    } catch (error) {
      console.error("Environment validation failed:", error);

      return new Response(
        JSON.stringify({
          error: "Configuration error",
          details: error instanceof Error ? error.message : String(error),
        }),
        {
          status: 500,
          headers: { "Content-Type": "application/json" },
        },
      );
    }

    return app.fetch(request, env);
  },
};
