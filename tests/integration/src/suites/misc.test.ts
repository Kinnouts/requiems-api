/**
 * Integration tests — Misc API endpoints
 *
 * Covers: /v1/misc/random-user, /v1/misc/convert/units
 */

import { describe, it } from "vitest";
import * as client from "../client.js";
import { assertEnvelope, repeat } from "../helpers.js";

const SUITE = "misc";

describe("Misc API", () => {
  describe("GET /v1/misc/random-user", () => {
    it("returns a randomly generated user", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/misc/random-user");
        const { data } = await assertEnvelope(response, SUITE, "random_user");
        const d = data as Record<string, unknown>;
        expect(typeof d["name"]).toBe("string");
        expect(typeof d["email"]).toBe("string");
        expect(typeof d["phone"]).toBe("string");
        expect(typeof d["avatar"]).toBe("string");
        expect(typeof d["address"]).toBe("object");

        const address = d["address"] as Record<string, unknown>;
        expect(typeof address["city"]).toBe("string");
        expect(typeof address["country"]).toBe("string");
      });
    });
  });

  describe("GET /v1/misc/convert", () => {
    it("converts 5 kilometers to miles", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/misc/convert", {
          value: "5",
          from: "km",
          to: "mi",
        });
        const { data } = await assertEnvelope(response, SUITE, "convert_units");
        const d = data as Record<string, unknown>;
        expect(typeof d["result"]).toBe("number");
        // 5 km ≈ 3.107 miles
        expect(d["result"] as number).toBeGreaterThan(3);
        expect(d["result"] as number).toBeLessThan(4);
      });
    });
  });
});
