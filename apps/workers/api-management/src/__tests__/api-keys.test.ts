import { beforeEach, describe, expect, it } from "vitest";
import worker from "../index";
import { authedRequest, makeBindings, makeCtx, makeDB, makeKV } from "./helpers";
import type { WorkerBindings } from "../env";
import type { ApiKeyData } from "@requiem/workers-shared";

// A realistic key prefix (first 12 chars of a requiem_ key)
const KEY_PREFIX = "requiem_abcd";
const FULL_KEY = "requiem_abcdefghijklmnopqrstuvwx";

describe("GET /api-keys", () => {
  it("returns empty list when no keys exist", async () => {
    const req = authedRequest("/api-keys");
    const res = await worker.fetch(req, makeBindings(), makeCtx());

    expect(res.status).toBe(200);
    const body = (await res.json()) as { keys: unknown[]; total: number };
    expect(body.keys).toEqual([]);
    expect(body.total).toBe(0);
  });

  it("returns keys from D1", async () => {
    const row = {
      key_prefix: KEY_PREFIX,
      user_id: "user-1",
      plan: "free",
      active: 1,
      created_at: "2024-01-01T00:00:00Z",
      updated_at: null,
      revoked_at: null,
      billing_cycle_start: "2024-01-01T00:00:00Z",
    };
    const bindings = makeBindings({ DB: makeDB([row]) });
    const req = authedRequest("/api-keys");
    const res = await worker.fetch(req, bindings, makeCtx());

    expect(res.status).toBe(200);
    const body = (await res.json()) as {
      keys: Array<{ keyPrefix: string; userId: string; active: boolean }>;
      total: number;
    };
    expect(body.total).toBe(1);
    expect(body.keys[0].keyPrefix).toBe(KEY_PREFIX);
    expect(body.keys[0].userId).toBe("user-1");
    expect(body.keys[0].active).toBe(true);
  });

  it("maps active=0 rows to active: false", async () => {
    const row = {
      key_prefix: KEY_PREFIX,
      user_id: "user-2",
      plan: "developer",
      active: 0,
      created_at: "2024-01-01T00:00:00Z",
      updated_at: null,
      revoked_at: "2024-06-01T00:00:00Z",
      billing_cycle_start: "2024-01-01T00:00:00Z",
    };
    const bindings = makeBindings({ DB: makeDB([row]) });
    const req = authedRequest("/api-keys?active=false");
    const res = await worker.fetch(req, bindings, makeCtx());

    expect(res.status).toBe(200);
    const body = (await res.json()) as { keys: Array<{ active: boolean }> };
    expect(body.keys[0].active).toBe(false);
  });
});

describe("POST /api-keys", () => {
  let kvStore: Map<string, string>;
  let bindings: WorkerBindings;

  beforeEach(() => {
    kvStore = new Map();
    bindings = makeBindings({
      KV: makeKV(kvStore),
      DB: makeDB(),
    });
  });

  it("returns 400 for an empty body", async () => {
    const req = authedRequest("/api-keys", { method: "POST" });
    const res = await worker.fetch(req, bindings, makeCtx());

    expect(res.status).toBe(400);
  });

  it("returns 400 when required fields are missing", async () => {
    const req = authedRequest("/api-keys", {
      method: "POST",
      body: JSON.stringify({ userId: "u1" }), // missing plan and name
      headers: { "Content-Type": "application/json" },
    });
    const res = await worker.fetch(req, bindings, makeCtx());

    expect(res.status).toBe(400);
  });

  it("returns 400 for an invalid plan value", async () => {
    const req = authedRequest("/api-keys", {
      method: "POST",
      body: JSON.stringify({ userId: "u1", plan: "invalid-plan", name: "My Key" }),
      headers: { "Content-Type": "application/json" },
    });
    const res = await worker.fetch(req, bindings, makeCtx());

    expect(res.status).toBe(400);
  });

  it("creates a key and returns 201 with apiKey and keyPrefix", async () => {
    const req = authedRequest("/api-keys", {
      method: "POST",
      body: JSON.stringify({ userId: "u1", plan: "free", name: "Test Key" }),
      headers: { "Content-Type": "application/json" },
    });
    const res = await worker.fetch(req, bindings, makeCtx());

    expect(res.status).toBe(201);
    const body = (await res.json()) as {
      apiKey: string;
      keyPrefix: string;
      userId: string;
      plan: string;
    };
    expect(body.apiKey).toMatch(/^requiem_/);
    expect(body.keyPrefix).toHaveLength(12);
    expect(body.userId).toBe("u1");
    expect(body.plan).toBe("free");
  });

  it("writes the key and prefix index to KV after creation", async () => {
    const req = authedRequest("/api-keys", {
      method: "POST",
      body: JSON.stringify({ userId: "u2", plan: "developer", name: "Dev Key" }),
      headers: { "Content-Type": "application/json" },
    });
    const res = await worker.fetch(req, bindings, makeCtx());
    expect(res.status).toBe(201);

    const body = (await res.json()) as { apiKey: string; keyPrefix: string };
    expect(kvStore.has(`key:${body.apiKey}`)).toBe(true);
    expect(kvStore.has(`prefix:${body.keyPrefix}`)).toBe(true);
  });

  it("returns 409 when the generated key already exists in KV", async () => {
    // Pre-seed KV to simulate an extremely unlikely collision
    const keyData: ApiKeyData = {
      userId: "u3",
      plan: "free",
      createdAt: "2024-01-01T00:00:00Z",
    };

    // We cannot know the generated key in advance, so instead we inject a KV
    // that always returns an existing value for any key: lookup
    const collidingKV = {
      get: async (key: string) => {
        if (key.startsWith("key:")) return JSON.stringify(keyData);
        return null;
      },
      put: async () => {},
      delete: async () => {},
      list: async () => ({ keys: [], list_complete: true, cursor: "" }),
      getWithMetadata: async () => ({ value: null, metadata: null }),
    } as unknown as KVNamespace;

    const collisionBindings = makeBindings({ KV: collidingKV });
    const req = authedRequest("/api-keys", {
      method: "POST",
      body: JSON.stringify({ userId: "u3", plan: "free", name: "Collision Key" }),
      headers: { "Content-Type": "application/json" },
    });
    const res = await worker.fetch(req, collisionBindings, makeCtx());

    expect(res.status).toBe(409);
  });
});

describe("DELETE /api-keys/:keyPrefix", () => {
  it("returns 404 when the prefix is not in KV", async () => {
    const req = authedRequest(`/api-keys/${KEY_PREFIX}`, { method: "DELETE" });
    const res = await worker.fetch(req, makeBindings(), makeCtx());

    expect(res.status).toBe(404);
    const body = (await res.json()) as { error: string };
    expect(body.error).toMatch(/not found/i);
  });

  it("revokes the key and removes both KV entries", async () => {
    const kvStore = new Map<string, string>();
    const keyData: ApiKeyData = {
      userId: "u1",
      plan: "free",
      createdAt: "2024-01-01T00:00:00Z",
    };
    kvStore.set(`key:${FULL_KEY}`, JSON.stringify(keyData));
    kvStore.set(`prefix:${KEY_PREFIX}`, FULL_KEY);

    const bindings = makeBindings({ KV: makeKV(kvStore) });
    const req = authedRequest(`/api-keys/${KEY_PREFIX}`, { method: "DELETE" });
    const res = await worker.fetch(req, bindings, makeCtx());

    expect(res.status).toBe(200);
    const body = (await res.json()) as { success: boolean; keyPrefix: string };
    expect(body.success).toBe(true);
    expect(body.keyPrefix).toBe(KEY_PREFIX);

    // Both KV entries must be deleted
    expect(kvStore.has(`key:${FULL_KEY}`)).toBe(false);
    expect(kvStore.has(`prefix:${KEY_PREFIX}`)).toBe(false);
  });
});

describe("PATCH /api-keys/:keyPrefix", () => {
  it("returns 404 when the prefix is not in KV", async () => {
    const req = authedRequest(`/api-keys/${KEY_PREFIX}`, {
      method: "PATCH",
      body: JSON.stringify({ plan: "developer" }),
      headers: { "Content-Type": "application/json" },
    });
    const res = await worker.fetch(req, makeBindings(), makeCtx());

    expect(res.status).toBe(404);
  });

  it("returns 400 when no updatable fields are provided", async () => {
    const kvStore = new Map<string, string>();
    kvStore.set(`prefix:${KEY_PREFIX}`, FULL_KEY);
    const keyData: ApiKeyData = {
      userId: "u1",
      plan: "free",
      createdAt: "2024-01-01T00:00:00Z",
    };
    kvStore.set(`key:${FULL_KEY}`, JSON.stringify(keyData));

    const bindings = makeBindings({ KV: makeKV(kvStore) });
    const req = authedRequest(`/api-keys/${KEY_PREFIX}`, {
      method: "PATCH",
      body: JSON.stringify({}), // no plan or billingCycleStart
      headers: { "Content-Type": "application/json" },
    });
    const res = await worker.fetch(req, bindings, makeCtx());

    expect(res.status).toBe(400);
  });

  it("updates the plan and reflects new plan in the response", async () => {
    const kvStore = new Map<string, string>();
    kvStore.set(`prefix:${KEY_PREFIX}`, FULL_KEY);
    const keyData: ApiKeyData = {
      userId: "u1",
      plan: "free",
      createdAt: "2024-01-01T00:00:00Z",
    };
    kvStore.set(`key:${FULL_KEY}`, JSON.stringify(keyData));

    const bindings = makeBindings({ KV: makeKV(kvStore) });
    const req = authedRequest(`/api-keys/${KEY_PREFIX}`, {
      method: "PATCH",
      body: JSON.stringify({ plan: "developer" }),
      headers: { "Content-Type": "application/json" },
    });
    const res = await worker.fetch(req, bindings, makeCtx());

    expect(res.status).toBe(200);
    const body = (await res.json()) as { success: boolean; plan: string };
    expect(body.success).toBe(true);
    expect(body.plan).toBe("developer");
  });

  it("persists the updated plan to KV", async () => {
    const kvStore = new Map<string, string>();
    kvStore.set(`prefix:${KEY_PREFIX}`, FULL_KEY);
    const keyData: ApiKeyData = {
      userId: "u1",
      plan: "free",
      createdAt: "2024-01-01T00:00:00Z",
    };
    kvStore.set(`key:${FULL_KEY}`, JSON.stringify(keyData));

    const bindings = makeBindings({ KV: makeKV(kvStore) });
    const req = authedRequest(`/api-keys/${KEY_PREFIX}`, {
      method: "PATCH",
      body: JSON.stringify({ plan: "business" }),
      headers: { "Content-Type": "application/json" },
    });
    await worker.fetch(req, bindings, makeCtx());

    const stored = JSON.parse(kvStore.get(`key:${FULL_KEY}`) ?? "{}") as ApiKeyData;
    expect(stored.plan).toBe("business");
  });
});
