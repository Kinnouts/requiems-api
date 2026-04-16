import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    environment: "node",
    include: ["src/**/*.test.ts"],
    globals: true,
    // Long timeouts — these make real HTTP requests against production
    testTimeout: 30_000,
    hookTimeout: 15_000,
    // Run all test files in a single fork so the stats singleton is shared
    // across all suites and persisted once at the end.
    pool: "forks",
    poolOptions: {
      forks: {
        singleFork: true,
      },
    },
    // Clears the stats temp file once before all tests start
    globalSetup: ["./src/globalSetup.ts"],
    // Loads .env and registers the afterAll timing-persistence hook
    setupFiles: ["./src/setup.ts"],
    // Custom reporters — default + our performance reporter
    reporters: ["default", "./src/reporter.ts"],
  },
});
