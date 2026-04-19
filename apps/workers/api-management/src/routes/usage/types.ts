/**
 * Usage data record from D1.
 * id is used internally for cursor-based pagination and is not exposed in API responses.
 */
export interface UsageRecord {
  id: number;
  api_key: string;
  user_id: string;
  endpoint: string;
  credits_used: number;
  request_method: string;
  status_code: number;
  response_time_ms: number;
  used_at: string;
}

/**
 * Usage export response with ID-based cursor pagination.
 * Pass nextCursor as the `cursor` query param to fetch the next page.
 */
export interface UsageExportResponse {
  usage: Omit<UsageRecord, "id">[];
  hasMore: boolean;
  nextCursor?: string;
}
