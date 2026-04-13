/**
 * Integration tests — Entertainment API endpoints
 *
 * Covers: /v1/entertainment/chuck-norris, /v1/entertainment/jokes/dad,
 *         /v1/entertainment/facts, /v1/entertainment/trivia,
 *         /v1/entertainment/emoji/random, /v1/entertainment/sudoku,
 *         /v1/entertainment/horoscope/{sign}
 */

import { describe, it } from "vitest";
import * as client from "../client.js";
import { assertEnvelope, repeat } from "../helpers.js";

const SUITE = "entertainment";

describe("Entertainment API", () => {
  describe("GET /v1/entertainment/chuck-norris", () => {
    it("returns a Chuck Norris joke", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/entertainment/chuck-norris");
        const { data } = await assertEnvelope(response, SUITE, "chuck_norris");
        const d = data as Record<string, unknown>;
        expect(typeof d["joke"]).toBe("string");
        expect((d["joke"] as string).length).toBeGreaterThan(0);
      });
    });
  });

  describe("GET /v1/entertainment/jokes/dad", () => {
    it("returns a dad joke", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/entertainment/jokes/dad");
        const { data } = await assertEnvelope(response, SUITE, "dad_joke");
        const d = data as Record<string, unknown>;
        expect(typeof d["joke"]).toBe("string");
        expect((d["joke"] as string).length).toBeGreaterThan(0);
      });
    });
  });

  describe("GET /v1/entertainment/facts", () => {
    it("returns a random fact", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/entertainment/facts");
        const { data } = await assertEnvelope(response, SUITE, "fact");
        const d = data as Record<string, unknown>;
        expect(typeof d["fact"]).toBe("string");
        expect((d["fact"] as string).length).toBeGreaterThan(0);
      });
    });
  });

  describe("GET /v1/entertainment/trivia", () => {
    it("returns a trivia question", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/entertainment/trivia");
        const { data } = await assertEnvelope(response, SUITE, "trivia");
        const d = data as Record<string, unknown>;
        expect(typeof d["question"]).toBe("string");
      });
    });
  });

  describe("GET /v1/entertainment/emoji/random", () => {
    it("returns a random emoji", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/entertainment/emoji/random");
        const { data } = await assertEnvelope(response, SUITE, "emoji_random");
        const d = data as Record<string, unknown>;
        expect(typeof d["emoji"]).toBe("string");
      });
    });
  });

  describe("GET /v1/entertainment/sudoku", () => {
    it("returns a sudoku puzzle", async () => {
      const { response } = await client.get("/v1/entertainment/sudoku");
      const { data } = await assertEnvelope(response, SUITE, "sudoku");
      const d = data as Record<string, unknown>;
      expect(Array.isArray(d["puzzle"])).toBe(true);
      // A standard sudoku board has 9 rows
      expect((d["puzzle"] as unknown[]).length).toBe(9);
    });
  });

  describe("GET /v1/entertainment/horoscope/{sign}", () => {
    it("returns a horoscope for aries", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/entertainment/horoscope/aries");
        const { data } = await assertEnvelope(response, SUITE, "horoscope_aries");
        const d = data as Record<string, unknown>;
        expect(typeof d["horoscope"]).toBe("string");
      });
    });
  });
});
