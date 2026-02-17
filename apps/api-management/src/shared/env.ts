// SHARED FILE - Keep in sync with auth-gateway/src/shared/env.ts

import { z } from "zod";

export interface CloudflareBindings {
  KV: KVNamespace;
  DB: D1Database;
  API_MANAGEMENT_API_KEY: string;
  SWAGGER_USERNAME?: string;
  SWAGGER_PASSWORD?: string;
  ENVIRONMENT: "development" | "staging" | "production";
}

const envSchema = z.object({
  API_MANAGEMENT_API_KEY: z.string().min(32),
  SWAGGER_USERNAME: z.string().optional(),
  SWAGGER_PASSWORD: z.string().optional(),
  ENVIRONMENT: z.enum(["development", "staging", "production"]).default("production"),
});

export function validateEnv(env: CloudflareBindings): CloudflareBindings {
  envSchema.parse({
    API_MANAGEMENT_API_KEY: env.API_MANAGEMENT_API_KEY,
    SWAGGER_USERNAME: env.SWAGGER_USERNAME,
    SWAGGER_PASSWORD: env.SWAGGER_PASSWORD,
    ENVIRONMENT: env.ENVIRONMENT,
  });

  return env;
}

export type ValidatedEnv = z.infer<typeof envSchema>;
