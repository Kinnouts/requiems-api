import { Hono } from "hono";
import type { WorkerBindings } from "../../env";
import exportRoute from "./export";

const usageRoute = new Hono<{ Bindings: WorkerBindings }>();

// Mount endpoint routes
usageRoute.route("/", exportRoute);

export { usageRoute };
