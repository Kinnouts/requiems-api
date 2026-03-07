import * as z from "zod";

export const planSchema = z.enum([
  "free",
  "developer",
  "business",
  "professional",
  "enterprise",
]);
