import { Hono } from "hono";

import { validateEnv, type WorkerBindings } from "./env";
import { errorHandler, jsonResponse, corsMiddleware } from "@requiem/workers-shared";

import { apiKeyAuthMiddleware } from "./middleware/api-key-auth";

import proxyRoute from "./routes/proxy";

const app = new Hono<{ Bindings: WorkerBindings }>();

app.get("/healthz", (_c) => jsonResponse({ status: "ok", service: "auth-gateway" }));

app.use("*", corsMiddleware);

app.use("/*", apiKeyAuthMiddleware);

app.route("/", proxyRoute);

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

      return  jsonResponse({ error: "Configuration error" },500);
    }

    return app.fetch(request, env);
  },
};
