/**
 * Generates an OpenAPI 3.0 spec from api_docs YAML files in the dashboard.
 * Run: pnpm generate:openapi
 * Auto-runs before dev and deploy via predev/predeploy hooks.
 */

import { readFileSync, writeFileSync, readdirSync, mkdirSync } from "fs";
import { join, resolve, dirname } from "path";
import { fileURLToPath } from "url";
import yaml from "js-yaml";

const __dirname = dirname(fileURLToPath(import.meta.url));
const monorepoRoot = resolve(__dirname, "../../../../");
const apiDocsDir = join(monorepoRoot, "apps/dashboard/config/api_docs");
const catalogPath = join(monorepoRoot, "apps/dashboard/config/api_catalog.yml");
const outputDir = join(__dirname, "../src/generated");
const outputPath = join(outputDir, "openapi.ts");

// --- Types matching the YAML structure ---

interface YamlParameter {
  name: string;
  type: string;
  required: boolean;
  location?: string; // "query" | "path" | "body" | undefined (defaults to body)
  description?: string;
  example?: unknown;
}

interface YamlError {
  code?: number | string;
  status?: number | string;
  message?: string;
  description?: string;
}

interface YamlEndpoint {
  name: string;
  method: string;
  path: string;
  description?: string;
  parameters?: YamlParameter[];
  request_example?: string;
  response_example?: string;
  response_fields?: { name: string; type: string; description?: string }[];
  errors?: YamlError[];
}

interface YamlApiDoc {
  api_id: string;
  api_name: string;
  description?: string;
  endpoints?: YamlEndpoint[];
}

interface CatalogEntry {
  id: string;
  name: string;
  description?: string;
}

// --- Helpers ---

function yamlTypeToSchema(type: string): Record<string, unknown> {
  if (type === "array") return { type: "array", items: {} };
  if (type === "object") return { type: "object" };
  if (type === "integer") return { type: "integer" };
  if (type === "number") return { type: "number" };
  if (type === "boolean") return { type: "boolean" };
  return { type: "string" };
}

function getErrorStatus(error: YamlError): number {
  if (typeof error.code === "number") return error.code;
  if (typeof error.status === "number") return error.status;
  if (typeof error.status === "string") {
    const parsed = parseInt(error.status, 10);
    if (!isNaN(parsed)) return parsed;
  }
  const knownCodes: Record<string, number> = {
    validation_failed: 422,
    not_found: 404,
    unauthorized: 401,
    forbidden: 403,
    internal_error: 500,
    bad_request: 400,
    service_unavailable: 503,
  };
  return knownCodes[String(error.code)] ?? 400;
}

function buildOperation(
  endpoint: YamlEndpoint,
  apiId: string,
): Record<string, unknown> {
  const params = endpoint.parameters ?? [];

  // Path and query parameters
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

  // Body parameters
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
    ["POST", "PUT", "PATCH"].includes(method)
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

  // Build responses
  const responses: Record<string, unknown> = {};

  // 200 success
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
                properties:
                  dataProperties,
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

// --- Main ---

const catalog = yaml.load(readFileSync(catalogPath, "utf8")) as {
  apis: CatalogEntry[];
};
const catalogMap = new Map<string, CatalogEntry>();
for (const entry of catalog.apis ?? []) {
  catalogMap.set(entry.id, entry);
}

const paths: Record<string, Record<string, unknown>> = {};
const tags: { name: string; description?: string }[] = [];

const docFiles = readdirSync(apiDocsDir)
  .filter((f) => f.endsWith(".yml") || f.endsWith(".yaml"))
  .sort();

for (const file of docFiles) {
  const doc = yaml.load(
    readFileSync(join(apiDocsDir, file), "utf8"),
  ) as YamlApiDoc;

  if (!doc?.api_id || !doc.endpoints?.length) continue;

  // Add tag
  const catalogEntry = catalogMap.get(doc.api_id);
  tags.push({
    name: doc.api_id,
    description: catalogEntry?.description ?? doc.description ?? doc.api_name,
  });

  for (const endpoint of doc.endpoints) {
    const pathKey = endpoint.path;
    const method = endpoint.method.toLowerCase();

    if (!paths[pathKey]) paths[pathKey] = {};
    paths[pathKey][method] = buildOperation(endpoint, doc.api_id);
  }
}

const spec = {
  openapi: "3.0.3",
  info: {
    title: "Requiems API",
    version: "1.0.0",
    description:
      "Unified access to enterprise-grade APIs — email validation, text utilities, and more. Authenticate with the `requiems-api-key` header.",
  },
  servers: [{ url: "https://api.requiems.xyz", description: "Production" }],
  components: {
    securitySchemes: {
      "requiems-api-key": {
        type: "apiKey",
        in: "header",
        name: "requiems-api-key",
        description: "Your Requiems API key",
      },
    },
  },
  security: [{ "requiems-api-key": [] }],
  tags,
  paths,
};

mkdirSync(outputDir, { recursive: true });
writeFileSync(
  outputPath,
  `// AUTO-GENERATED — do not edit. Run \`pnpm generate:openapi\` to regenerate.\n` +
    `export const openApiSpec = ${JSON.stringify(spec, null, 2)};\n`,
);

console.log(
  `✅ OpenAPI spec generated: ${Object.keys(paths).length} paths across ${docFiles.length} APIs`,
);
