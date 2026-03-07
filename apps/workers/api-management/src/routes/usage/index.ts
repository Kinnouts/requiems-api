import { Hono } from "hono";
import type { WorkerBindings } from "../../env";
import exportRoute from "./export";

const usageRoute = new Hono<{ Bindings: WorkerBindings }>();

usageRoute.route("/", exportRoute);

export { usageRoute };
