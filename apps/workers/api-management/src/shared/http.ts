// SHARED FILE - Keep in sync with auth-gateway/src/shared/http.ts

/**
 * JSON response helper
 */
export function jsonResponse(
  data: unknown,
  status = 200,
  headers: Record<string, string> = {},
): Response {
  return new Response(JSON.stringify(data), {
    status,
    headers: {
      "Content-Type": "application/json",
      ...headers,
    },
  });
}

/**
 * JSON error response helper
 */
export function jsonError(
  status: number,
  message: string,
  headers: Record<string, string> = {},
): Response {
  return jsonResponse({ error: message }, status, headers);
}

/**
 * Middleware to check X-API-Management-Key header
 * Only the Rails dashboard should have this key
 */
export function requireApiManagementKey(request: Request, expectedKey: string): Response | null {
  const apiKey = request.headers.get("X-API-Management-Key");

  if (!apiKey || apiKey !== expectedKey) {
    return jsonError(401, "Unauthorized - Invalid or missing API management key");
  }

  return null;
}
