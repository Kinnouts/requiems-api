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
 * Filter headers before forwarding to backend
 * Removes Cloudflare headers and sensitive data
 */
export function filterHeaders(headers: Headers): Headers {
	const filtered = new Headers();

	for (const [key, value] of headers.entries()) {
		const lowerKey = key.toLowerCase();

		// Skip Cloudflare-specific headers
		if (lowerKey.startsWith("cf-")) continue;

		// Skip API key (backend trusts gateway)
		if (lowerKey === "x-api-key") continue;

		// Skip hop-by-hop headers
		if (lowerKey === "connection") continue;
		if (lowerKey === "keep-alive") continue;

		filtered.set(key, value);
	}

	return filtered;
}

/**
 * Add credit and rate limit headers to response
 */
export function addUsageHeaders(
	response: Response,
	headers: {
		creditsUsed: number;
		creditsRemaining: number;
		creditsReset: string;
		plan: string;
		rateLimitLimit: number;
		rateLimitRemaining: number;
	},
): Response {
	// Clone response to modify headers
	const newResponse = new Response(response.body, {
		status: response.status,
		statusText: response.statusText,
		headers: response.headers,
	});

	newResponse.headers.set("X-Credits-Used", headers.creditsUsed.toString());
	newResponse.headers.set(
		"X-Credits-Remaining",
		headers.creditsRemaining.toString(),
	);
	newResponse.headers.set("X-Credits-Reset", headers.creditsReset);
	newResponse.headers.set("X-Plan", headers.plan);
	newResponse.headers.set(
		"X-RateLimit-Limit",
		headers.rateLimitLimit.toString(),
	);
	newResponse.headers.set(
		"X-RateLimit-Remaining",
		headers.rateLimitRemaining.toString(),
	);

	return newResponse;
}
