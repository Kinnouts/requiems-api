/**
 * Configuration loader.
 *
 * Reads API_BASE_URL and REQUIEMS_API_KEY from the environment.
 * Fails fast with a helpful message if either is missing.
 */

export interface Config {
  baseUrl: string;
  apiKey: string;
  /** How many times each endpoint is exercised to compute timing stats */
  runs: number;
}

let _config: Config | undefined;

export function getConfig(): Config {
  if (_config) return _config;

  const baseUrl = process.env["API_BASE_URL"] ?? "https://api.requiems.xyz";
  const apiKey = process.env["REQUIEMS_API_KEY"] ?? "";
  const runs = Number(process.env["INTEGRATION_RUNS"] ?? "20");

  if (!apiKey) {
    throw new Error(
      [
        "",
        "⛔  REQUIEMS_API_KEY is not set.",
        "",
        "   Copy tests/integration/.env.example to tests/integration/.env",
        "   and fill in your production API key, then re-run the tests.",
        "",
      ].join("\n"),
    );
  }

  _config = { baseUrl, apiKey, runs };
  return _config;
}
