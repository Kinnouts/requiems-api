/**
 * Integration tests — Networking & validation endpoints (formerly "tech")
 *
 * Covers:
 *   /v1/networking/ip           (IP geolocation)
 *   /v1/networking/mx/{domain}  (MX record lookup)
 *   /v1/technology/password     (password generator)
 *   /v1/technology/useragent    (user agent parser)
 *   /v1/validation/phone        (phone validation)
 */

import { describe, it } from "vitest";
import * as client from "../client.js";
import { assertEnvelope, repeat } from "../helpers.js";

const SUITE = "tech";

describe("Networking & Technology API", () => {
  describe("GET /v1/networking/ip", () => {
    it("returns IP geolocation info for the caller's IP", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/networking/ip");
        const { data } = await assertEnvelope(response, SUITE, "ip_lookup");
        const d = data as Record<string, unknown>;
        expect(typeof d["ip"]).toBe("string");
        expect(typeof d["country"]).toBe("string");
      });
    });

    it("returns IP info for a specific public IP", async () => {
      const { response } = await client.get("/v1/networking/ip/8.8.8.8");
      const { data } = await assertEnvelope(
        response,
        SUITE,
        "ip_lookup_specific",
      );
      const d = data as Record<string, unknown>;
      expect(d["ip"]).toBe("8.8.8.8");
    });
  });

  describe("GET /v1/technology/password", () => {
    it("generates a password of default length", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/technology/password");
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "password_default",
        );
        const d = data as Record<string, unknown>;
        expect(typeof d["password"]).toBe("string");
        // Default length is 16
        expect((d["password"] as string).length).toBe(16);
      });
    });

    it("generates a password of a custom length", async () => {
      const { response } = await client.get("/v1/technology/password", {
        length: "24",
      });
      const { data } = await assertEnvelope(response, SUITE, "password_custom");
      const d = data as Record<string, unknown>;
      expect((d["password"] as string).length).toBe(24);
    });
  });

  describe("GET /v1/technology/useragent", () => {
    it("parses a Chrome user agent string", async () => {
      const ua =
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36";
      await repeat(async () => {
        const { response } = await client.get("/v1/technology/useragent", {
          ua,
        });
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "useragent_chrome",
        );
        const d = data as Record<string, unknown>;
        expect(typeof d["browser"]).toBe("string");
        expect(typeof d["os"]).toBe("string");
        expect(d["is_bot"]).toBe(false);
      });
    });
  });

  describe("GET /v1/validation/phone", () => {
    it("validates a well-formed US phone number", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/validation/phone", {
          number: "+14155552671",
        });
        const { data } = await assertEnvelope(response, SUITE, "phone_us");
        const d = data as Record<string, unknown>;
        expect(typeof d["valid"]).toBe("boolean");
        expect(typeof d["number"]).toBe("string");
      });
    });
  });

  describe("GET /v1/networking/mx/{domain}", () => {
    it("returns MX records for gmail.com", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/networking/mx/gmail.com");
        const { data } = await assertEnvelope(response, SUITE, "mx_gmail");
        const d = data as Record<string, unknown>;
        expect(Array.isArray(d["records"])).toBe(true);
        expect((d["records"] as unknown[]).length).toBeGreaterThan(0);
      });
    });
  });
});
