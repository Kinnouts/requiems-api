import type { Context, ErrorHandler } from "hono";
import { jsonResponse } from "../http";

/**
 * Global Error Handler for Hono Application
 *
 * Handles all unhandled errors in the application.
 * Logs errors for debugging and returns appropriate error responses.
 */
export const errorHandler: ErrorHandler = (err, _c: Context) => {
  console.error("Unhandled error:", {
    message: err.message,
    name: err.name,
    stack: err.stack,
  });

  return jsonResponse(
    {
      error: "Internal server error",
      message: err.message,
    },
    500,
  );
};
