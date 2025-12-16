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
 */
export type PlanName = "free" | "starter" | "pro" | "business" | "enterprise";

/**
 * Plan configuration
 */
export interface PlanConfig {
  /** Credit limit for the period */
  creditLimit: number;
  /** Credit period (daily for free, monthly for paid) */
  creditPeriod: "daily" | "monthly";
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
 * Credit check result
 */
export interface CreditCheckResult {
  usage: number;
  remaining: number;
  limit: number;
  resetAt: string;
}
