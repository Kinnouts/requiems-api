/**
 * Usage data record from D1
 */
export interface UsageRecord {
	api_key: string;
	endpoint: string;
	credits_used: number;
	used_at: string;
}

/**
 * Usage export response with pagination
 */
export interface UsageExportResponse {
	usage: UsageRecord[];
	total: number;
	hasMore: boolean;
	nextCursor?: string;
}
