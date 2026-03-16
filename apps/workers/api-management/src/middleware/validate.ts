import { sValidator } from "@hono/standard-validator";
import type { StandardSchemaV1 } from "@standard-schema/spec";
import { jsonError } from "@requiem/workers-shared";

export function validateQuery<S extends StandardSchemaV1>(schema: S) {
  return sValidator("query", schema, (result, _c) => {
    if (!result.success) {
      return jsonError(400, result.error[0]?.message ?? "Validation error");
    }
  });
}

export function validateJson<S extends StandardSchemaV1>(schema: S) {
  return sValidator("json", schema, (result, _c) => {
    if (!result.success) {
      return jsonError(400, result.error[0]?.message ?? "Validation error");
    }
  });
}
