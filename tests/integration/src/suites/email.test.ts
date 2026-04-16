/**
 * Integration tests — Email-related endpoints
 *
 * Covers:
 *   /v1/networking/disposable/check   (disposable domain detection)
 *   /v1/networking/disposable/stats
 *   /v1/validation/email              (full email validation)
 *   /v1/text/normalize                (email normalization)
 */

import { describe, it } from "vitest";
import * as client from "../client.js";
import { assertEnvelope, repeat } from "../helpers.js";

const SUITE = "email";

describe("Email API", () => {
  describe("POST /v1/networking/disposable/check", () => {
    it("identifies a disposable domain", async () => {
      await repeat(async () => {
        const { response } = await client.post(
          "/v1/networking/disposable/check",
          {
            email: "test@mailinator.com",
          },
        );
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "disposable_check",
        );
        const d = data as Record<string, unknown>;
        expect(d["disposable"]).toBe(true);
      });
    });

    it("identifies a legitimate domain", async () => {
      const { response } = await client.post(
        "/v1/networking/disposable/check",
        {
          email: "user@gmail.com",
        },
      );
      const { data } = await assertEnvelope(
        response,
        SUITE,
        "disposable_check_legit",
      );
      const d = data as Record<string, unknown>;
      expect(d["disposable"]).toBe(false);
    });
  });

  describe("POST /v1/validation/email", () => {
    it("validates a well-formed email", async () => {
      await repeat(async () => {
        const { response } = await client.post("/v1/validation/email", {
          email: "user@example.com",
        });
        const { data } = await assertEnvelope(response, SUITE, "validate");
        const d = data as Record<string, unknown>;
        expect(typeof d["valid"]).toBe("boolean");
        expect(typeof d["syntax_valid"]).toBe("boolean");
      });
    });
  });

  describe("POST /v1/text/normalize", () => {
    it("normalizes a gmail alias address", async () => {
      const { response } = await client.post("/v1/text/normalize", {
        email: "User+alias@Gmail.com",
      });
      const { data } = await assertEnvelope(response, SUITE, "normalize");
      const d = data as Record<string, unknown>;
      expect(typeof d["normalized"]).toBe("string");
      // The normalized form should be lowercase
      expect((d["normalized"] as string).toLowerCase()).toBe(d["normalized"]);
    });
  });

  describe("GET /v1/networking/disposable/stats", () => {
    it("returns disposable domain statistics", async () => {
      const { response } = await client.get(
        "/v1/networking/disposable/stats",
      );
      const { data } = await assertEnvelope(
        response,
        SUITE,
        "disposable_stats",
      );
      const d = data as Record<string, unknown>;
      expect(typeof d["total_domains"]).toBe("number");
      expect(d["total_domains"] as number).toBeGreaterThan(0);
    });
  });
});
