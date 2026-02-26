import { z } from "zod";
import type { BaseWorkerBindings } from "@requiem/workers-shared";

export interface WorkerBindings extends BaseWorkerBindings {
  BACKEND_URL: string;
  BACKEND_SECRET: string;
}

const envSchema = z.object({
  BACKEND_URL: z.string().url(),
  BACKEND_SECRET: z.string().min(32),
  ENVIRONMENT: z.enum(["development", "staging", "production"]).default(
    "production",
  ),
});

export function validateEnv(env: WorkerBindings): WorkerBindings {
  envSchema.parse({
    BACKEND_URL: env.BACKEND_URL,
    BACKEND_SECRET: env.BACKEND_SECRET,
    ENVIRONMENT: env.ENVIRONMENT,
  });

  return env;
}

export type ValidatedEnv = z.infer<typeof envSchema>;
