import { Hono } from "hono";
import type { WorkerBindings } from "../../shared/types";
import byEndpointRoute from "./by-endpoint";
import byDateRoute from "./by-date";
import summaryRoute from "./summary";

const app = new Hono<{ Bindings: WorkerBindings }>();

// Mount endpoint routes
app.route("/", byEndpointRoute);
app.route("/", byDateRoute);
app.route("/", summaryRoute);

export default app;
