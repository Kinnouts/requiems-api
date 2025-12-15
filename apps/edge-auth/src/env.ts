import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

/**
 * Cloudflare Worker bindings (KV, D1, etc.)
 * These are passed by Cloudflare, not process.env
 */
export interface CloudflareBindings {
  /** Cloudflare KV namespace for API keys, config, rate limits */
  KV: KVNamespace;
  /** Cloudflare D1 database for usage tracking */
  DB: D1Database;
}

/**
 * Validated environment variables using t3-env
 *
 * With nodejs_compat flag, we can use process.env directly.
 * Cloudflare vars/secrets are automatically available in process.env.
 */
export const env = createEnv({
  server: {
    /**
     * Internal backend URL (kept secret)
     * Example: https://internal-backend.fly.dev
     */
    BACKEND_URL: z.string().url(),

    /**
     * Secret key to authenticate gateway -> backend requests
     * Backend rejects requests without this header
     */
    BACKEND_SECRET: z.string().min(32),

    /**
     * Environment name
     */
    ENVIRONMENT: z
      .enum(["development", "staging", "production"])
      .default("development"),
  },

  /**
   * Use process.env directly with nodejs_compat
   */
  runtimeEnv: process.env,

  /**
   * Treat empty strings as undefined
   */
  emptyStringAsUndefined: true,
});

export type ValidatedEnv = typeof env;
