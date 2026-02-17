import { Hono } from "hono";
import { basicAuth } from "hono/basic-auth";
import { swaggerUI } from "@hono/swagger-ui";

import { validateEnv } from "./shared/env";
import { corsResponse, jsonResponse } from "./shared/http";
import type { WorkerBindings } from "./shared/types";

import apiKeysRoute from "./routes/api-keys";
import usageRoute from "./routes/usage";
import analyticsRoute from "./routes/analytics";

const app = new Hono<{ Bindings: WorkerBindings }>();

app.options("*", (_c) => corsResponse);

app.get("/healthz", (_c) => jsonResponse({ status: "ok", service: "api-management" }));

app.use("/docs/*", async (c, next) => {
  const env = c.env;

  if (env.ENVIRONMENT === "production" && env.SWAGGER_USERNAME && env.SWAGGER_PASSWORD) {
    const auth = basicAuth({
      username: env.SWAGGER_USERNAME,
      password: env.SWAGGER_PASSWORD,
    });
    
    return auth(c, next);
  }

  await next();
});

app.get("/docs", swaggerUI({ url: "/openapi.json" }));

app.get("/openapi.json", (c) => {
  return c.json({
    openapi: "3.0.0",
    info: {
      title: "Requiems API Management",
      version: "1.0.0",
      description:
        "Internal API for managing API keys, usage data, and analytics. Only accessible by Rails dashboard.",
    },
    servers: [
      { url: "https://api-management.requiems.xyz", description: "Production" },
      { url: "http://localhost:6001", description: "Local development" },
    ],
    components: {
      securitySchemes: {
        ApiManagementKey: {
          type: "apiKey",
          in: "header",
          name: "X-API-Management-Key",
          description: "API Management key (only Rails dashboard has this)",
        },
      },
    },
    security: [{ ApiManagementKey: [] }],
    paths: {
      "/healthz": {
        get: {
          summary: "Health check",
          security: [],
          responses: {
            "200": {
              description: "Service is healthy",
            },
          },
        },
      },
      "/api-keys": {
        post: {
          summary: "Create, revoke, or update API keys",
          requestBody: {
            required: true,
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  required: ["action", "key", "userId", "plan"],
                  properties: {
                    action: { type: "string", enum: ["create", "revoke", "update"] },
                    key: { type: "string" },
                    userId: { type: "string" },
                    plan: {
                      type: "string",
                      enum: ["free", "developer", "business", "professional", "enterprise"],
                    },
                    billingCycleStart: { type: "string", format: "date-time" },
                  },
                },
              },
            },
          },
          responses: {
            "200": { description: "Success" },
            "401": { description: "Unauthorized" },
          },
        },
      },
      "/usage/export": {
        get: {
          summary: "Export usage data from D1",
          parameters: [
            { name: "since", in: "query", required: true, schema: { type: "string" } },
            { name: "limit", in: "query", schema: { type: "integer" } },
            { name: "cursor", in: "query", schema: { type: "string" } },
          ],
          responses: {
            "200": { description: "Usage data with pagination" },
          },
        },
      },
      "/analytics/by-endpoint": {
        get: {
          summary: "Usage breakdown by endpoint",
          parameters: [{ name: "userId", in: "query", required: true, schema: { type: "string" } }],
          responses: {
            "200": { description: "Endpoint usage statistics" },
          },
        },
      },
      "/analytics/by-date": {
        get: {
          summary: "Usage trends over time",
          parameters: [{ name: "userId", in: "query", required: true, schema: { type: "string" } }],
          responses: {
            "200": { description: "Time series usage data" },
          },
        },
      },
      "/analytics/summary": {
        get: {
          summary: "Overall usage summary",
          parameters: [{ name: "userId", in: "query", required: true, schema: { type: "string" } }],
          responses: {
            "200": { description: "Usage summary" },
          },
        },
      },
    },
  });
});

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
