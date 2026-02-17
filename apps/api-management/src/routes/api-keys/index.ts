import { Hono } from "hono";

import createRoute from "./create";
import deleteRoute from "./delete";
import patchRoute from "./patch";

import type { WorkerBindings } from "../../shared/types";


const app = new Hono<{ Bindings: WorkerBindings }>();

app.route("/", createRoute);
app.route("/", deleteRoute);
app.route("/", patchRoute);

export default app;
