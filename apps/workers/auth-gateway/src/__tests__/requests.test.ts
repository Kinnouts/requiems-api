import { describe, expect, it, vi } from "vitest";
import { recordRequestUsage } from "../requests";
import type { WorkerBindings } from "../env";

function makeBindings() {
  const kv = new Map<string, string>();
  const bind = vi.fn((...args: unknown[]) => ({
    run: vi.fn().mockResolvedValue({ success: true, meta: {} }),
  }));
  const prepare = vi.fn((_sql: string) => ({
    bind,
  }));

  return {
    bindings: {
      KV: {
        get: vi.fn(async (key: string) => kv.get(key) ?? null),
        put: vi.fn(async (key: string, value: string) => {
          kv.set(key, value);
        }),
      } as unknown as KVNamespace,
      DB: { prepare } as unknown as D1Database,
    } as WorkerBindings,
    mocks: { prepare, bind },
  };
}

describe("recordRequestUsage", () => {
  it("writes request method and telemetry fields to D1", async () => {
    const { bindings, mocks } = makeBindings();

    await recordRequestUsage(
      bindings,
      "requiem_test_key",
      "user-1",
      "/v1/text/advice",
      2,
      503,
      128,
      "PATCH",
      "2026-04-01T00:00:00.000Z",
    );

    expect(mocks.prepare).toHaveBeenCalled();
    expect(mocks.bind).toHaveBeenCalled();
    const bindArgs = mocks.bind.mock.calls[0] as unknown[];
    
    // Verify all arguments except the timestamp (which is generated at runtime)
    expect(bindArgs[0]).toBe("requiem_test_key");
    expect(bindArgs[1]).toBe("user-1");
    expect(bindArgs[2]).toBe("/v1/text/advice");
    expect(bindArgs[3]).toBe(2);
    expect(bindArgs[4]).toBe("PATCH");
    expect(bindArgs[5]).toBe(503);
    expect(bindArgs[6]).toBe(128);
    // bindArgs[7] is the timestamp, which is generated at call time
    expect(typeof bindArgs[7]).toBe("string");
    expect((bindArgs[7] as string).length).toBeGreaterThan(0);
  });
});