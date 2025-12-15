export interface Env {
  BACKEND_ORIGIN: string;
  API_KEY_SECRET: string;
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);

    // Simple API key auth: expect `x-api-key` header to match secret
    const clientKey = request.headers.get("x-api-key");
    if (!clientKey || clientKey !== env.API_KEY_SECRET) {
      return new Response(JSON.stringify({ error: "unauthorized" }), {
        status: 401,
        headers: { "Content-Type": "application/json" },
      });
    }

    // Basic rate limiting etc. can be added later here.

    const backendURL = new URL(url.pathname + url.search, env.BACKEND_ORIGIN);

    const init: RequestInit = {
      method: request.method,
      headers: request.headers,
      body: request.body,
    };

    return fetch(backendURL.toString(), init);
  },
};
