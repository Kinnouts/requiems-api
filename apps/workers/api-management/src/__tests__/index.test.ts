import { describe, expect, it } from "vitest";
import worker from "../index";
import { makeBindings, makeCtx } from "./helpers";

describe("worker smoke tests", () => {
  it("responds to requests without throwing", async () => {
    const req = new Request("http://localhost/healthz");
    const res = await worker.fetch(req, makeBindings(), makeCtx());

    expect(res).toBeInstanceOf(Response);
  });

  it("serves JSON responses", async () => {
    const req = new Request("http://localhost/healthz");
    const res = await worker.fetch(req, makeBindings(), makeCtx());

    expect(res.headers.get("Content-Type")).toBe("application/json");
  });
});
