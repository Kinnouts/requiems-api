/**
 * Endpoint usage stats
 */
export interface EndpointStats {
	endpoint: string;
	requests: number;
	credits: number;
}

/**
 * Date-based usage stats
 */
export interface DateStats {
	date: string;
	requests: number;
	credits: number;
}

/**
 * Usage summary for a user
 */
export interface UsageSummary {
	userId: string;
	totalRequests: number;
	totalCredits: number;
	dateRange: {
		since: string;
		until: string;
	};
	topEndpoints: EndpointStats[];
}
