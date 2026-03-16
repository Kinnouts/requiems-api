import { readFileSync, readdirSync } from "fs";
import { join, resolve } from "path";

import yaml from "js-yaml";

import type { YamlError, YamlEndpoint, CatalogEntry, YamlApiDoc } from "./types";


const TYPE_SCHEMAS: Record<string, Record<string, unknown>> = {
  array: { type: "array", items: {} },
  object: { type: "object" },
  integer: { type: "integer" },
  number: { type: "number" },
  boolean: { type: "boolean" },
  string: { type: "string" },
};

export function yamlTypeToSchema(type: string): Record<string, unknown> {
  return TYPE_SCHEMAS[type] ?? { type: "string" };
}

const KNOWN_ERROR_CODES: Record<string, number> = {
  validation_failed: 422,
  not_found: 404,
  unauthorized: 401,
  forbidden: 403,
  internal_error: 500,
  bad_request: 400,
  service_unavailable: 503,
};

export function getErrorStatus(error: YamlError): number {
  if (typeof error.code === "number") return error.code;
  if (typeof error.status === "number") return error.status;

  if (typeof error.status === "string") {
    const parsed = Number.parseInt(error.status, 10);
    if (!Number.isNaN(parsed)) return parsed;
  }

  return KNOWN_ERROR_CODES[String(error.code)] ?? 400;
}

const METHODS_WITH_BODY =  ["POST", "PUT", "PATCH"]

export function buildOperation(
  endpoint: YamlEndpoint,
  apiId: string,
): Record<string, unknown> {
  const params = endpoint.parameters ?? [];

  const pathAndQueryParams = params
    .filter((p) => p.location === "path" || p.location === "query")
    .map((p) => {
      const schema: Record<string, unknown> = {
        ...yamlTypeToSchema(p.type),
      };

      if (p.example !== undefined) schema.example = p.example;

      const param: Record<string, unknown> = {
        name: p.name,
        in: p.location,
        required: p.location === "path" ? true : (p.required ?? false),
        schema,
      };

      if (p.description) param.description = p.description;
      return param;
    });

  const bodyParams = params.filter(
    (p) => !p.location || p.location === "body",
  );

  const operation: Record<string, unknown> = {
    summary: endpoint.name,
    tags: [apiId],
    security: [{ "requiems-api-key": [] }],
  };

  if (endpoint.description) operation.description = endpoint.description;
  if (pathAndQueryParams.length > 0) operation.parameters = pathAndQueryParams;

  // Build requestBody for POST/PUT/PATCH with body params
  const method = endpoint.method.toUpperCase();

  if (
    bodyParams.length > 0 &&
    METHODS_WITH_BODY.includes(method)
  ) {
    const properties: Record<string, unknown> = {};
    const required: string[] = [];

    for (const p of bodyParams) {
      const schema: Record<string, unknown> = { ...yamlTypeToSchema(p.type) };
      if (p.description) schema.description = p.description;
      if (p.example !== undefined) schema.example = p.example;
      properties[p.name] = schema;
      if (p.required) required.push(p.name);
    }

    const bodySchema: Record<string, unknown> = {
      type: "object",
      properties,
    };

    if (required.length > 0) bodySchema.required = required;

    // Include request example if available
    if (endpoint.request_example) {
      try {
        bodySchema.example = JSON.parse(endpoint.request_example);
      } catch {
        // skip malformed examples
      }
    }

    operation.requestBody = {
      required: true,
      content: {
        "application/json": { schema: bodySchema },
      },
    };
  }

  const responses: Record<string, unknown> = {};


  const successResponse: Record<string, unknown> = {
    description: "Successful response",
  };

  if (endpoint.response_example) {
    try {
      const example = JSON.parse(endpoint.response_example);
      const responseFields = endpoint.response_fields ?? [];
      const dataProperties: Record<string, unknown> = {};

      for (const field of responseFields) {
        const fieldSchema: Record<string, unknown> = {
          ...yamlTypeToSchema(field.type),
        };

        if (field.description) fieldSchema.description = field.description;
        dataProperties[field.name] = fieldSchema;
      }

      successResponse.content = {
        "application/json": {
          schema: {
            type: "object",
            properties: {
              data: {
                type: "object",
                properties: dataProperties,
              },
              metadata: {
                type: "object",
                properties: {
                  timestamp: { type: "string", format: "date-time" },
                },
              },
            },
          },
          example,
        },
      };
    } catch {
      // skip malformed response examples
    }
  }

  responses["200"] = successResponse;

  // Error responses
  const errorsByStatus: Record<number, string[]> = {};

  for (const err of endpoint.errors ?? []) {
    const status = getErrorStatus(err);

    if (!errorsByStatus[status]) errorsByStatus[status] = [];

    const msg = err.description ?? err.message ?? String(err.code ?? "Error");
    errorsByStatus[status].push(msg);
  }

  for (const [status, messages] of Object.entries(errorsByStatus)) {
    responses[status] = { description: messages.join("; ") };
  }

  operation.responses = responses;

  return operation;
}

const monorepoRoot = resolve(import.meta.dirname, "../../../../../");


export const apiDocsDir = join(monorepoRoot, "apps/dashboard/config/api_docs");
export const catalogPath = join(monorepoRoot, "apps/dashboard/config/api_catalog.yml");


export function loadCatalog(){
  try {
    return  yaml.load(readFileSync(catalogPath, "utf8")) as {
      apis: CatalogEntry[];
    };
  } catch (err) {
    console.error(`❌ Failed to read catalog at ${catalogPath}:`, err);
    process.exit(1);
  }
}

export function loadAPIDocs(){
    try {
     return readdirSync(apiDocsDir)
      .filter((f) => f.endsWith(".yml") || f.endsWith(".yaml"))
      .sort();
  } catch (err) {
    console.error(`❌ Failed to read api_docs directory at ${apiDocsDir}:`, err);
    process.exit(1);
  }
}

export function loadAPIDoc(file:string): YamlApiDoc | null {
      try {
        return yaml.load(readFileSync(join(apiDocsDir, file), "utf8")) as YamlApiDoc;
      } catch (err) {
        console.warn(`⚠️  Skipping ${file}: failed to parse YAML —`, err);
        return null;
      }
}
