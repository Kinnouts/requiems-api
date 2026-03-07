import { Hono } from "hono";

import type { WorkerBindings } from "../../env";

import byEndpointRoute from "./by-endpoint";
import byDateRoute from "./by-date";
import summaryRoute from "./summary";

const analyticsRoute = new Hono<{ Bindings: WorkerBindings }>();

analyticsRoute.route("/", byEndpointRoute);
analyticsRoute.route("/", byDateRoute);
analyticsRoute.route("/", summaryRoute);

export { analyticsRoute };
