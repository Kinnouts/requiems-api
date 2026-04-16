/**
 * Maps the exact endpoint paths recorded by the integration test client to
 * their corresponding API documentation markdown files (relative to the
 * project root).
 *
 * Rules:
 * - Use the path exactly as passed to client.get() / client.post() in the
 *   test suites (including any concrete param values like /v1/finance/crypto/BTC).
 * - When multiple paths correspond to the same doc (e.g. both /encode and
 *   /decode share docs/apis/technology/base64.md), list them all — the
 *   injection script will pick the entry with the most samples.
 * - Paths with no matching doc file should be omitted entirely.
 */

export const endpointDocMap: Record<string, string> = {
  // ── Networking ─────────────────────────────────────────────────────────────
  "/v1/networking/disposable/check": "docs/apis/networking/disposable-email.md",
  "/v1/networking/disposable/stats": "docs/apis/networking/disposable-email.md",
  "/v1/networking/ip": "docs/apis/networking/ip-geolocation.md",
  "/v1/networking/ip/8.8.8.8": "docs/apis/networking/ip-geolocation.md",
  "/v1/networking/mx/gmail.com": "docs/apis/networking/mx-lookup.md",

  // ── Validation ─────────────────────────────────────────────────────────────
  "/v1/validation/email": "docs/apis/validation/email-validate.md",
  "/v1/validation/phone": "docs/apis/validation/phone-validation.md",
  "/v1/validation/profanity": "docs/apis/validation/profanity.md",

  // ── Text ───────────────────────────────────────────────────────────────────
  "/v1/text/normalize": "docs/apis/validation/email-validate.md",
  "/v1/text/words/random": "docs/apis/text/random-word.md",
  "/v1/text/lorem": "docs/apis/text/lorem-ipsum.md",
  "/v1/text/dictionary/eloquent": "docs/apis/text/dictionary.md",
  "/v1/text/thesaurus/happy": "docs/apis/text/thesaurus.md",
  "/v1/text/spellcheck": "docs/apis/text/spell-check.md",

  // ── Entertainment ──────────────────────────────────────────────────────────
  "/v1/entertainment/advice": "docs/apis/entertainment/advice.md",
  "/v1/entertainment/quotes/random": "docs/apis/entertainment/quotes.md",
  "/v1/entertainment/chuck-norris": "docs/apis/entertainment/chuck-norris.md",
  "/v1/entertainment/jokes/dad": "docs/apis/entertainment/dad-jokes.md",
  "/v1/entertainment/facts": "docs/apis/entertainment/facts.md",
  "/v1/entertainment/trivia": "docs/apis/entertainment/trivia.md",
  "/v1/entertainment/emoji/random": "docs/apis/entertainment/emoji.md",
  "/v1/entertainment/sudoku": "docs/apis/entertainment/sudoku.md",
  "/v1/entertainment/horoscope/aries": "docs/apis/entertainment/horoscope.md",

  // ── Finance ────────────────────────────────────────────────────────────────
  "/v1/finance/mortgage": "docs/apis/finance/mortgage.md",
  "/v1/finance/inflation": "docs/apis/finance/inflation.md",
  "/v1/finance/exchange-rate": "docs/apis/finance/exchange-rate.md",
  "/v1/finance/convert": "docs/apis/finance/exchange-rate.md",
  "/v1/finance/crypto/BTC": "docs/apis/finance/crypto-price.md",
  "/v1/finance/commodities/gold": "docs/apis/finance/commodities.md",

  // ── Technology ─────────────────────────────────────────────────────────────
  "/v1/technology/password": "docs/apis/technology/password-generator.md",
  "/v1/technology/useragent": "docs/apis/technology/user-agent.md",
  "/v1/technology/base64/encode": "docs/apis/technology/base64.md",
  "/v1/technology/base64/decode": "docs/apis/technology/base64.md",
  "/v1/technology/base": "docs/apis/technology/number-base-conversion.md",
  "/v1/technology/color": "docs/apis/technology/color-conversion.md",
  "/v1/technology/markdown": "docs/apis/technology/data-format-conversion.md",
  "/v1/technology/random-user": "docs/apis/technology/random-user.md",
  "/v1/technology/convert": "docs/apis/technology/unit-conversion.md",
};
