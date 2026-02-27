import { Hono } from "hono";

import { validateEnv, type WorkerBindings } from "./env";
import { createWorkerFetch, errorHandler, notFoundHandler } from "@requiem/workers-shared";

import { apiKeyAuthMiddleware, docsMiddleware } from "./middleware/";

import { analyticsRoute, apiKeysRoute, healthzRoute, swaggerRoute, usageRoute } from "./routes";

const app = new Hono<{ Bindings: WorkerBindings }>();

app.route("/", healthzRoute);

app.use("/docs", docsMiddleware);

app.route("/", swaggerRoute);

app.use("/api-keys/*", apiKeyAuthMiddleware);
app.use("/usage/*", apiKeyAuthMiddleware);
app.use("/analytics/*", apiKeyAuthMiddleware);

app.route("/api-keys", apiKeysRoute);
app.route("/usage", usageRoute);
app.route("/analytics", analyticsRoute);

app.notFound(notFoundHandler);
app.onError(errorHandler);

export default createWorkerFetch(app, validateEnv);
