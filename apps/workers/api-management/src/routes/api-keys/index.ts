import { Hono } from "hono";

import createRoute from "./create";
import deleteRoute from "./delete";
import listRoute from "./list";
import patchRoute from "./patch";

import type { WorkerBindings } from "../../env";

const apiKeysRoute = new Hono<{ Bindings: WorkerBindings }>();

apiKeysRoute.route("/", listRoute);
apiKeysRoute.route("/", createRoute);
apiKeysRoute.route("/", deleteRoute);
apiKeysRoute.route("/", patchRoute);

export { apiKeysRoute };
