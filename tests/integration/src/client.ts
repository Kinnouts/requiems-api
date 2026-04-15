/**
 * HTTP client with built-in latency tracking.
 *
 * Every call returns the response together with wall-clock timing so callers
 * can feed results into the statistics tracker.
 */

import { getConfig } from "./config.js";
import { stats } from "./stats.js";

export interface TimedResponse {
  response: Response;
  /** Total round-trip time in milliseconds */
  durationMs: number;
  status: number;
  ok: boolean;
}

/**
 * Perform a GET request against the production gateway.
 *
 * @param path   API path (e.g. "/v1/text/advice")
 * @param params Optional query-string parameters
 */
export async function get(
  path: string,
  params?: Record<string, string>,
): Promise<TimedResponse> {
  return request("GET", path, undefined, params);
}

/**
 * Perform a POST request against the production gateway.
 *
 * @param path API path (e.g. "/v1/email/validate")
 * @param body JSON-serialisable request body
 */
export async function post(
  path: string,
  body: unknown,
): Promise<TimedResponse> {
  return request("POST", path, body);
}

async function request(
  method: string,
  path: string,
  body?: unknown,
  params?: Record<string, string>,
): Promise<TimedResponse> {
  const cfg = getConfig();

  const url = new URL(path, cfg.baseUrl);
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      url.searchParams.set(k, v);
    }
  }

  const headers: Record<string, string> = {
    "requiems-api-key": cfg.apiKey,
  };

  let requestBody: string | undefined;
  if (body !== undefined) {
    headers["Content-Type"] = "application/json";
    requestBody = JSON.stringify(body);
  }

  const start = performance.now();
  const response = await fetch(url.toString(), {
    method,
    headers,
    body: requestBody,
  });
  const durationMs = Math.round(performance.now() - start);

  // Record timing for the summary report
  stats.record(path, durationMs);

  return { response, durationMs, status: response.status, ok: response.ok };
}
