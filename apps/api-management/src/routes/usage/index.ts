import { Hono } from "hono";
import type { WorkerBindings } from "../../shared/types";
import exportRoute from "./export";

const app = new Hono<{ Bindings: WorkerBindings }>();

// Mount endpoint routes
app.route("/", exportRoute);

export default app;
