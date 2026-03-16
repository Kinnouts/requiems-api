import { z } from "zod";

/**
 * Creates a validateEnv function for the given zod schema.
 * Workers call the returned function in their fetch handler to validate
 * environment bindings at runtime and get a typed result back.
 */
export function createEnvValidator<T extends z.ZodTypeAny>(
  schema: T,
): (env: z.infer<T>) => z.infer<T> {
  return (env) => {
    schema.parse(env);
    return env;
  };
}
