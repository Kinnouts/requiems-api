import type { PlanName } from "./config";

// Re-export PlanName for convenience
export type { PlanName };

/**
 * Base Cloudflare bindings shared by all workers
 * Workers should extend this with their specific env vars
 */
export interface BaseWorkerBindings {
  KV: KVNamespace;
  DB: D1Database;
  ENVIRONMENT: "development" | "staging" | "production";
}

/**
 * API key data stored in KV
 */
export interface ApiKeyData {
  userId: string;
  plan: PlanName;
  createdAt: string;
  /** Optional: billing cycle start for paid plans */
  billingCycleStart?: string;
}

/**
 * Plan configuration
 * All plans use monthly billing cycles
 */
export interface PlanConfig {
  /** Monthly request limit */
  requestLimit: number;
  /** Max requests per minute */
  ratePerMinute: number;
}

/**
 * Rate limit check result
 */
export interface RateLimitResult {
  allowed: boolean;
  remaining: number;
  resetAt: number;
}

/**
 * Request usage check result
 */
export interface RequestCheckResult {
  usage: number;
  remaining: number;
  limit: number;
  resetAt: string;
}

/**
 * API key management request from Rails
 */
export interface ApiKeyManagementRequest {
  action: "create" | "revoke" | "update";
  key: string;
  userId: string;
  plan: PlanName;
  billingCycleStart?: string;
}

/**
 * API key management response to Rails
 */
export interface ApiKeyManagementResponse {
  success: boolean;
  message?: string;
  error?: string;
}
