import { describe, it, expect, vi, beforeEach } from "vitest";
import { api } from "./api";

// Mock the sentry module so captureException doesn't try to init Sentry
vi.mock("./sentry", () => ({
  captureException: vi.fn(),
}));

const API_URL = "http://localhost:8080";

describe("api.auth.me", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("calls the correct URL with credentials: include", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ user: { id: "1", name: "Test" } }),
    });
    global.fetch = mockFetch;

    await api.auth.me();

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/auth/me`,
      expect.objectContaining({
        credentials: "include",
      })
    );
  });

  it("returns { data: ... } on successful response", async () => {
    const user = { id: "1", name: "Test User", email: "test@example.com" };
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ user }),
    });

    const result = await api.auth.me();

    expect(result).toEqual({ data: { user } });
  });

  it("returns { error: string } on error response", async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({ error: "Unauthorized" }),
    });

    const result = await api.auth.me();

    expect(result).toEqual({ error: "Unauthorized" });
  });

  it("returns { error: 'Something went wrong' } when error response has no error field", async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({}),
    });

    const result = await api.auth.me();

    expect(result).toEqual({ error: "Something went wrong" });
  });

  it("returns { error: 'Network error' } on network failure", async () => {
    global.fetch = vi.fn().mockRejectedValue(new Error("Failed to fetch"));

    const result = await api.auth.me();

    expect(result).toEqual({ error: "Network error" });
  });
});

describe("api.auth.logout", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("sends POST request to /api/auth/logout", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ success: true }),
    });
    global.fetch = mockFetch;

    await api.auth.logout();

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/auth/logout`,
      expect.objectContaining({
        method: "POST",
        credentials: "include",
      })
    );
  });
});

describe("api.users.onboard", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("sends POST with username and name", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ success: true }),
    });
    global.fetch = mockFetch;

    await api.users.onboard({ username: "testuser", name: "Test" });

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/users/onboard`,
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ username: "testuser", name: "Test" }),
      })
    );
  });
});

describe("api.connections", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("list calls correct endpoint", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ connections: [] }),
    });
    global.fetch = mockFetch;

    await api.connections.list();

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/connections`,
      expect.objectContaining({ credentials: "include" })
    );
  });

  it("disconnect sends DELETE", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ success: true }),
    });
    global.fetch = mockFetch;

    await api.connections.disconnect("github");

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/connections/github`,
      expect.objectContaining({ method: "DELETE" })
    );
  });

  it("connectUrl returns correct URL", () => {
    expect(api.connections.connectUrl("youtube")).toBe(
      `${API_URL}/api/connections/youtube/connect`
    );
  });
});

describe("api.analytics", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("summary calls correct endpoint", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ totalViews: 100 }),
    });
    global.fetch = mockFetch;

    const result = await api.analytics.summary();

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/analytics`,
      expect.objectContaining({ credentials: "include" })
    );
    expect(result).toEqual({ data: { totalViews: 100 } });
  });

  it("trackView sends POST", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ ok: true }),
    });
    global.fetch = mockFetch;

    await api.analytics.trackView("testuser");

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/users/testuser/view`,
      expect.objectContaining({ method: "POST" })
    );
  });
});

describe("api.collaborations", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("send sends POST with toUserId and message", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ id: "c1" }),
    });
    global.fetch = mockFetch;

    await api.collaborations.send("user123", "Let's collab!");

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/collaborations`,
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ toUserId: "user123", message: "Let's collab!" }),
      })
    );
  });

  it("inbox calls correct endpoint", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ requests: [] }),
    });
    global.fetch = mockFetch;

    await api.collaborations.inbox();

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/collaborations/inbox`,
      expect.objectContaining({ credentials: "include" })
    );
  });

  it("respond sends POST with action", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ status: "accepted" }),
    });
    global.fetch = mockFetch;

    await api.collaborations.respond("c1", "accept");

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/collaborations/c1/respond`,
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ action: "accept" }),
      })
    );
  });
});

describe("api.adminModeration", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("list calls correct endpoint with params", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ flags: [], total: 0 }),
    });
    global.fetch = mockFetch;

    await api.adminModeration.list("pending", 20, 0);

    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/admin/moderation"),
      expect.objectContaining({ credentials: "include" })
    );
    // Verify the URL contains the status parameter
    const calledUrl = mockFetch.mock.calls[0][0];
    expect(calledUrl).toContain("status=pending");
    expect(calledUrl).toContain("limit=20");
  });

  it("resolve sends POST with status and notes", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ success: true }),
    });
    global.fetch = mockFetch;

    await api.adminModeration.resolve("flag1", "resolved", "Looks fine");

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/admin/moderation/flag1/resolve`,
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ status: "resolved", notes: "Looks fine" }),
      })
    );
  });
});

describe("api.discover", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("list builds correct query params", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ creators: [], total: 0 }),
    });
    global.fetch = mockFetch;

    await api.discover.list({ minScore: 50, platform: "youtube", q: "test" });

    const calledUrl = mockFetch.mock.calls[0][0];
    expect(calledUrl).toContain("minScore=50");
    expect(calledUrl).toContain("platform=youtube");
    expect(calledUrl).toContain("q=test");
  });
});

describe("api.notifications", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("list calls correct endpoint with params", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ notifications: [], total: 0 }),
    });
    global.fetch = mockFetch;

    await api.notifications.list(10, 5);

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/notifications?limit=10&offset=5`,
      expect.objectContaining({ credentials: "include" })
    );
  });

  it("markAllRead sends POST", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ success: true }),
    });
    global.fetch = mockFetch;

    await api.notifications.markAllRead();

    expect(mockFetch).toHaveBeenCalledWith(
      `${API_URL}/api/notifications/read-all`,
      expect.objectContaining({ method: "POST" })
    );
  });
});
