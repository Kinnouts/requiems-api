import { Hono } from "hono";
import { validateEnv } from "./shared/env";
import { corsResponse, jsonResponse } from "./shared/http";
import type { WorkerBindings } from "./shared/types";

import apiKeysRoute from "./routes/api-keys";
import usageRoute from "./routes/usage";
import analyticsRoute from "./routes/analytics";

const app = new Hono<{ Bindings: WorkerBindings }>();

// CORS preflight
app.options("*", (c) => corsResponse);

// Health check
app.get("/healthz", (c) => jsonResponse({ status: "ok", service: "api-management" }));

// Mount routes
app.route("/api-keys", apiKeysRoute);
app.route("/usage", usageRoute);
app.route("/analytics", analyticsRoute);

// 404 handler
app.notFound((c) => {
  return jsonResponse({ error: "Not found" }, 404);
});

// Error handler
app.onError((err, c) => {
  console.error("Unhandled error:", err);
  return jsonResponse({ error: "Internal server error" }, 500);
});

// Export fetch handler
export default {
  async fetch(request: Request, env: WorkerBindings): Promise<Response> {
    // Validate environment variables
    validateEnv(env);

    return app.fetch(request, env);
  },
};
