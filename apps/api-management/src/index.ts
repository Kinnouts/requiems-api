import { Hono } from "hono";
import { swaggerUI } from "@hono/swagger-ui";

import { validateEnv } from "./shared/env";
import { jsonResponse } from "./shared/http";

import { swaggerAuthMiddleware } from "./middleware/swagger-auth";
import { apiKeyAuthMiddleware } from "./middleware/api-key-auth";

import apiKeysRoute from "./routes/api-keys";
import usageRoute from "./routes/usage";
import analyticsRoute from "./routes/analytics";
import swaggerRoute from "./routes/swagger";

import type { WorkerBindings } from "./shared/types";

const app = new Hono<{ Bindings: WorkerBindings }>();

app.get("/healthz", (_c) => jsonResponse({ status: "ok", service: "api-management" }));

app.use("/docs/*", swaggerAuthMiddleware);

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

app.onError((err, _c) => {
  console.error("Unhandled error:", err);
  return jsonResponse({ error: "Internal server error" }, 500);
});

export default {
  async fetch(request: Request, env: WorkerBindings): Promise<Response> {
    validateEnv(env);

    return app.fetch(request, env);
  },
};
