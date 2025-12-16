import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

export interface CloudflareBindings {
  KV: KVNamespace;
  DB: D1Database;
}

export const env = createEnv({
  server: {
    BACKEND_URL: z.string().url(),
    BACKEND_SECRET: z.string().min(32),
    ENVIRONMENT: z
      .enum(["development", "staging", "production"])
      .default("development"),
  },
  runtimeEnv: process.env,
  emptyStringAsUndefined: true,
});

export type ValidatedEnv = typeof env;
