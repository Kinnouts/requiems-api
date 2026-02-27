import { Hono } from "hono";
import { jsonResponse } from "@requiem/workers-shared";
import type { WorkerBindings } from "../env";

const healthzRoute = new Hono<{ Bindings: WorkerBindings }>();

healthzRoute.get("/healthz", (c) =>
  jsonResponse({ status: "ok", service: "api-management" }),
);

export { healthzRoute };
