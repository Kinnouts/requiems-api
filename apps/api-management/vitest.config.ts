import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    // Test environment - Workers run in Node environment for testing
    environment: "node",

    // File patterns
    include: ["src/**/*.{test,spec}.ts"],
    exclude: ["node_modules", "dist"],

    // Coverage configuration
    coverage: {
      provider: "v8",
      reporter: ["text", "json", "html"],
      exclude: [
        "node_modules/",
        "dist/",
        "**/*.config.ts",
        "**/*.d.ts",
      ],
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
