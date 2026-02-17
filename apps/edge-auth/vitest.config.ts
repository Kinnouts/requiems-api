import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    // Test environment - Workers run in Node environment for testing
    environment: "node",

    // File patterns
    include: ["src/**/*.{test,spec}.ts"],
    exclude: ["node_modules", "dist", "scripts"],

    // Coverage configuration
    coverage: {
      provider: "v8",
      reporter: ["text", "json", "html"],
      exclude: [
        "node_modules/",
        "dist/",
        "scripts/",
        "**/*.config.ts",
        "**/*.d.ts",
      ],
      // Non-blocking: no thresholds initially
      // Can add thresholds later as coverage improves
    },

    // Output
    reporter: "default",

    // Performance - Workers are single-threaded
    pool: "threads",
    poolOptions: {
      threads: {
        singleThread: true,  // Match Cloudflare Workers single-threaded execution
      },
    },

    // Timeouts
    testTimeout: 10000,
    hookTimeout: 10000,

    // Globals - convenient for testing
    globals: true,
  },
});
