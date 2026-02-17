import { z } from "zod";

export interface CloudflareBindings {
  KV: KVNamespace;
  DB: D1Database;
  BACKEND_URL: string;
  BACKEND_SECRET: string;
  ENVIRONMENT: "development" | "staging" | "production";
}

const envSchema = z.object({
  BACKEND_URL: z.string().url(),
  BACKEND_SECRET: z.string().min(32),
  ENVIRONMENT: z.enum(["development", "staging", "production"]).default("production"),
});

export function validateEnv(env: CloudflareBindings): CloudflareBindings {
  envSchema.parse({
    BACKEND_URL: env.BACKEND_URL,
    BACKEND_SECRET: env.BACKEND_SECRET,
    ENVIRONMENT: env.ENVIRONMENT,
  });

  return env;
}

export type ValidatedEnv = z.infer<typeof envSchema>;
