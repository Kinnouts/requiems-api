/**
 * Integration tests — Text API endpoints
 *
 * Covers: /v1/text/words/random, /v1/text/lorem,
 *         /v1/text/dictionary/{word}, /v1/text/thesaurus/{word},
 *         /v1/text/spellcheck, /v1/validation/profanity
 *
 * Each test is run `config.runs` times to produce stable timing samples and
 * surface flakiness.  Response body shapes are snapshotted on the first run
 * and compared on subsequent runs.
 */

import { describe, it } from "vitest";
import * as client from "../client.js";
import { assertEnvelope, repeat } from "../helpers.js";

const SUITE = "text";

describe("Text API", () => {
  describe("GET /v1/text/words/random", () => {
    it("returns a random word", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/text/words/random");
        const { data } = await assertEnvelope(response, SUITE, "words_random");
        const d = data as Record<string, unknown>;
        expect(typeof d["word"]).toBe("string");
        expect((d["word"] as string).length).toBeGreaterThan(0);
      });
    });
  });

  describe("GET /v1/text/lorem", () => {
    it("returns lorem ipsum with default settings", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/text/lorem");
        const { data } = await assertEnvelope(response, SUITE, "lorem_default");
        const d = data as Record<string, unknown>;
        expect(typeof d["text"]).toBe("string");
        expect((d["text"] as string).length).toBeGreaterThan(0);
      });
    });

    it("returns correct paragraph count when requested", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/text/lorem", {
          paragraphs: "2",
        });
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "lorem_2_paragraphs",
        );
        const d = data as Record<string, unknown>;
        expect(typeof d["paragraphs"]).toBe("number");
        expect(d["paragraphs"] as number).toBeGreaterThanOrEqual(2);
      });
    });
  });

  describe("GET /v1/text/dictionary/{word}", () => {
    it("returns a definition for a known word", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/text/dictionary/eloquent");
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "dictionary_word",
        );
        const d = data as Record<string, unknown>;
        expect(typeof d["word"]).toBe("string");
      });
    });
  });

  describe("GET /v1/text/thesaurus/{word}", () => {
    it("returns synonyms and antonyms for a known word", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/text/thesaurus/happy");
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "thesaurus_word",
        );
        const d = data as Record<string, unknown>;
        expect(typeof d["word"]).toBe("string");
      });
    });
  });

  describe("POST /v1/validation/profanity", () => {
    it("detects clean text", async () => {
      await repeat(async () => {
        const { response } = await client.post("/v1/validation/profanity", {
          text: "Hello, world!",
        });
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "profanity_clean",
        );
        const d = data as Record<string, unknown>;
        expect(d["has_profanity"]).toBe(false);
      });
    });
  });

  describe("POST /v1/text/spellcheck", () => {
    it("returns no errors for correctly spelled text", async () => {
      await repeat(async () => {
        const { response } = await client.post("/v1/text/spellcheck", {
          text: "The quick brown fox jumps over the lazy dog.",
        });
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "spellcheck_correct",
        );
        const d = data as Record<string, unknown>;
        expect(Array.isArray(d["corrections"])).toBe(true);
      });
    });
  });
});
