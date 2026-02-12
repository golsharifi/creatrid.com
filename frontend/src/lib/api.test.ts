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
