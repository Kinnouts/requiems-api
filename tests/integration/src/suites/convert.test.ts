/**
 * Integration tests — Technology conversion endpoints (formerly "convert")
 *
 * Covers:
 *   /v1/technology/base64/encode
 *   /v1/technology/base64/decode
 *   /v1/technology/base   (number base conversion)
 *   /v1/technology/color  (color format conversion)
 *   /v1/technology/markdown
 */

import { describe, it } from "vitest";
import * as client from "../client.js";
import { assertEnvelope, repeat } from "../helpers.js";

const SUITE = "convert";

describe("Technology — Conversion API", () => {
  describe("POST /v1/technology/base64/encode", () => {
    it("encodes a string to Base64", async () => {
      await repeat(async () => {
        const { response } = await client.post(
          "/v1/technology/base64/encode",
          {
            value: "Hello, World!",
          },
        );
        const { data } = await assertEnvelope(response, SUITE, "base64_encode");
        const d = data as Record<string, unknown>;
        expect(d["result"]).toBe("SGVsbG8sIFdvcmxkIQ==");
      });
    });
  });

  describe("POST /v1/technology/base64/decode", () => {
    it("decodes a Base64 string", async () => {
      await repeat(async () => {
        const { response } = await client.post(
          "/v1/technology/base64/decode",
          {
            value: "SGVsbG8sIFdvcmxkIQ==",
          },
        );
        const { data } = await assertEnvelope(response, SUITE, "base64_decode");
        const d = data as Record<string, unknown>;
        expect(d["result"]).toBe("Hello, World!");
      });
    });
  });

  describe("GET /v1/technology/base", () => {
    it("converts decimal 255 to hexadecimal", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/technology/base", {
          from: "10",
          to: "16",
          value: "255",
        });
        const { data } = await assertEnvelope(response, SUITE, "base_convert");
        const d = data as Record<string, unknown>;
        expect((d["result"] as string).toLowerCase()).toBe("ff");
      });
    });
  });

  describe("GET /v1/technology/color", () => {
    it("converts a HEX color to RGB", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/technology/color", {
          from: "hex",
          to: "rgb",
          value: "#FF5733",
        });
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "color_hex_to_rgb",
        );
        expect(data).toBeTruthy();
      });
    });
  });

  describe("POST /v1/technology/markdown", () => {
    it("converts Markdown to HTML", async () => {
      await repeat(async () => {
        const { response } = await client.post("/v1/technology/markdown", {
          markdown: "# Hello\n\nThis is **bold**.",
        });
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "markdown_to_html",
        );
        const d = data as Record<string, unknown>;
        expect(typeof d["html"]).toBe("string");
        expect(d["html"] as string).toContain("<h1");
        expect(d["html"] as string).toContain("<strong>");
      });
    });
  });
});
