/**
 * Integration tests — Finance API endpoints
 *
 * Covers: /v1/finance/mortgage, /v1/finance/inflation, /v1/finance/exchange-rate,
 *         /v1/finance/convert, /v1/finance/crypto/{symbol},
 *         /v1/finance/commodities/{commodity}
 */

import { describe, it } from "vitest";
import * as client from "../client.js";
import { assertEnvelope, repeat } from "../helpers.js";

const SUITE = "finance";

describe("Finance API", () => {
  describe("GET /v1/finance/mortgage", () => {
    it("calculates a 30-year mortgage", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/finance/mortgage", {
          principal: "300000",
          rate: "6.5",
          years: "30",
        });
        const { data } = await assertEnvelope(response, SUITE, "mortgage");
        const d = data as Record<string, unknown>;
        expect(typeof d["monthly_payment"]).toBe("number");
        expect(typeof d["total_payment"]).toBe("number");
        expect(typeof d["total_interest"]).toBe("number");
        expect(Array.isArray(d["schedule"])).toBe(true);
        expect((d["schedule"] as unknown[]).length).toBe(30 * 12);
      });
    });
  });

  describe("GET /v1/finance/inflation", () => {
    it("returns inflation data for US", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/finance/inflation", {
          country: "US",
        });
        const { data } = await assertEnvelope(response, SUITE, "inflation_us");
        const d = data as Record<string, unknown>;
        expect(typeof d["rate"]).toBe("number");
        expect(typeof d["period"]).toBe("string");
      });
    });
  });

  describe("GET /v1/finance/exchange-rate", () => {
    it("returns USD → EUR exchange rate", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/finance/exchange-rate", {
          from: "USD",
          to: "EUR",
        });
        const { data } = await assertEnvelope(response, SUITE, "exchange_rate");
        const d = data as Record<string, unknown>;
        expect(d["from"]).toBe("USD");
        expect(d["to"]).toBe("EUR");
        expect(typeof d["rate"]).toBe("number");
        expect(d["rate"] as number).toBeGreaterThan(0);
      });
    });
  });

  describe("GET /v1/finance/convert", () => {
    it("converts 100 USD to GBP", async () => {
      const { response } = await client.get("/v1/finance/convert", {
        from: "USD",
        to: "GBP",
        amount: "100",
      });
      const { data } = await assertEnvelope(response, SUITE, "convert");
      const d = data as Record<string, unknown>;
      expect(d["from"]).toBe("USD");
      expect(d["to"]).toBe("GBP");
      expect(typeof d["converted"]).toBe("number");
      expect(d["converted"] as number).toBeGreaterThan(0);
    });
  });

  describe("GET /v1/finance/crypto/{symbol}", () => {
    it("returns price for BTC", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/finance/crypto/BTC");
        const { data } = await assertEnvelope(response, SUITE, "crypto_btc");
        const d = data as Record<string, unknown>;
        expect(typeof d["price"]).toBe("number");
        expect(d["price"] as number).toBeGreaterThan(0);
      });
    });
  });

  describe("GET /v1/finance/commodities/{commodity}", () => {
    it("returns price for gold", async () => {
      await repeat(async () => {
        const { response } = await client.get("/v1/finance/commodities/gold");
        const { data } = await assertEnvelope(
          response,
          SUITE,
          "commodities_gold",
        );
        const d = data as Record<string, unknown>;
        expect(typeof d["price"]).toBe("number");
        expect(d["price"] as number).toBeGreaterThan(0);
      });
    });
  });
});
