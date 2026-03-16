/**
 * Generates an OpenAPI 3.0 spec from api_docs YAML files in the dashboard.
 * Run: pnpm generate:openapi
 * Auto-runs before dev and deploy via predev/predeploy hooks.
 */

import { writeFileSync, mkdirSync } from "fs";
import { join} from "path";

import { buildOperation, loadAPIDoc, loadAPIDocs, loadCatalog } from "./helpers";

import type { CatalogEntry } from "./types";

const outputDir = join(import.meta.dirname, "../../src/generated");
const outputPath = join(outputDir, "openapi.ts");


function main() {
  const catalog = loadCatalog();

  const catalogMap = new Map<string, CatalogEntry>();
  
  for (const entry of catalog.apis ?? []) {
    catalogMap.set(entry.id, entry);
  }

  const docFiles= loadAPIDocs();


  const paths: Record<string, Record<string, unknown>> = {};
  const tags: { name: string; description?: string }[] = [];

  for (const file of docFiles) {
    const doc = loadAPIDoc(file);

    if (!doc || !doc?.api_id || !doc.endpoints?.length) continue;

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

  // Assemble spec
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

  // Write output
  try {
    mkdirSync(outputDir, { recursive: true });
    writeFileSync(
      outputPath,
      `// AUTO-GENERATED — do not edit. Run \`pnpm generate:openapi\` to regenerate.\n` +
        `export const openApiSpec = ${JSON.stringify(spec, null, 2)};\n`,
    );
  } catch (err) {
    console.error(`❌ Failed to write output to ${outputPath}:`, err);
    process.exit(1);
  }

  console.log(
    `✅ OpenAPI spec generated: ${Object.keys(paths).length} paths across ${docFiles.length} APIs`,
  );
}

main();
