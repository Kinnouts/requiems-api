import { describe, expect, it } from "vitest";
import worker from "../index";
import { authedRequest, makeBindings, makeCtx } from "./helpers";

describe("healthz route", () => {
  it("returns 200 with service name", async () => {
    const req = new Request("http://localhost/healthz");
    const res = await worker.fetch(req, makeBindings(), makeCtx());

    expect(res.status).toBe(200);
    const body = (await res.json()) as { status: string; service: string };
    expect(body.status).toBe("ok");
    expect(body.service).toBe("api-management");
  });

  it("does not require the API management key", async () => {
    // No X-API-Management-Key header — route is public
    const req = new Request("http://localhost/healthz");
    const res = await worker.fetch(req, makeBindings(), makeCtx());

    expect(res.status).toBe(200);
  });
});

describe("worker — unknown route", () => {
  it("returns 404 for unregistered paths", async () => {
    const req = authedRequest("/does-not-exist");
    const res = await worker.fetch(req, makeBindings(), makeCtx());

    expect(res.status).toBe(404);
  });
});

describe("worker — env validation", () => {
  it("returns 500 when API_MANAGEMENT_API_KEY is missing", async () => {
    const bindings = makeBindings({
      API_MANAGEMENT_API_KEY: "" as unknown as string,
    });

    const req = new Request("http://localhost/healthz");
    const res = await worker.fetch(req, bindings, makeCtx());

    // createWorkerFetch catches validation errors and returns 500
    expect(res.status).toBe(500);
  });
});
