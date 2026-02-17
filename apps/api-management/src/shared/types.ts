// SHARED FILE - Keep in sync with auth-gateway/src/shared/types.ts

import type { CloudflareBindings } from "./env";

/**
 * Worker bindings passed to fetch handler
 * Env vars are accessed via process.env (nodejs_compat)
 */
export type WorkerBindings = CloudflareBindings;

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
 * Plan names
 * https://github.com/bobadilla-tech/requiems-api/docs/business.md
 */
export type PlanName = "free" | "developer" | "business" | "professional" | "enterprise";

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
