import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    environment: "edge-runtime",
    include: ["src/**/*.{test,spec}.ts"],
    exclude: ["node_modules", "dist", "scripts"],
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
    },
    reporters: "default",
    pool: "threads",
    poolOptions: {
      threads: {
        singleThread: true,
      },
    },
    testTimeout: 10000,
    hookTimeout: 10000,
    globals: true,
  },
});
