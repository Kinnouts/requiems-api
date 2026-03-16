import type { BaseWorkerBindings } from "@requiem/workers-shared";
import { createEnvValidator } from "@requiem/workers-shared";

import { z } from "zod";

export interface WorkerBindings extends BaseWorkerBindings {
  API_MANAGEMENT_API_KEY: string;
  SWAGGER_USERNAME?: string;
  SWAGGER_PASSWORD?: string;
}

const envSchema = z.object({
  API_MANAGEMENT_API_KEY: z.string().min(32),
  SWAGGER_USERNAME: z.string().optional(),
  SWAGGER_PASSWORD: z.string().optional(),
  ENVIRONMENT: z.enum(["development", "staging", "production"]).default("production"),
});

export const validateEnv = createEnvValidator(envSchema);

export type ValidatedEnv = z.infer<typeof envSchema>;
