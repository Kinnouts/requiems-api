import { describe, expect, it, vi } from "vitest";
import { recordRequestUsage } from "../requests";
import type { WorkerBindings } from "../env";

function makeBindings() {
  const kv = new Map<string, string>();
  const prepare = vi.fn((_sql: string) => ({
    bind: (...args: unknown[]) => ({
      run: vi.fn().mockResolvedValue({ success: true, meta: {} }),
      args,
    }),
  }));

  return {
    KV: {
      get: vi.fn(async (key: string) => kv.get(key) ?? null),
      put: vi.fn(async (key: string, value: string) => {
        kv.set(key, value);
      }),
    } as unknown as KVNamespace,
    DB: { prepare } as unknown as D1Database,
  } as WorkerBindings;
}

describe("recordRequestUsage", () => {
  it("writes request method and telemetry fields to D1", async () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-04-19T12:00:00Z"));

    const bindings = makeBindings();
    const prepare = bindings.DB.prepare as unknown as ReturnType<typeof vi.fn>;

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

    expect(prepare).toHaveBeenCalled();
    const bindCall = (prepare.mock.results[0]?.value as { bind: (...args: unknown[]) => { args: unknown[] } }).bind;
    const bound = bindCall.mock.calls[0] as unknown[];

    expect(bound).toEqual([
      "requiem_test_key",
      "user-1",
      "/v1/text/advice",
      2,
      "PATCH",
      503,
      128,
      new Date("2026-04-19T12:00:00Z").toISOString(),
    ]);
  });
});