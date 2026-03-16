export type YamlParameterLocation = "query" | "path" | "body";

export interface YamlParameter {
  name: string;
  type: string;
  required: boolean;
  location?: YamlParameterLocation; // (defaults to body)
  description?: string;
  example?: unknown;
}

export interface YamlError {
  code?: number | string;
  status?: number | string;
  message?: string;
  description?: string;
}

export interface YamlEndpoint {
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

export interface YamlApiDoc {
  api_id: string;
  api_name: string;
  description?: string;
  endpoints?: YamlEndpoint[];
}

export interface CatalogEntry {
  id: string;
  name: string;
  description?: string;
}