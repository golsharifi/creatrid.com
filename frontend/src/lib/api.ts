import type { User, PublicUser, Connection } from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

type ApiResponse<T> =
  | { data: T; error?: never }
  | { data?: never; error: string };

async function request<T>(
  path: string,
  options?: RequestInit
): Promise<ApiResponse<T>> {
  try {
    const res = await fetch(`${API_URL}${path}`, {
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        ...options?.headers,
      },
      ...options,
    });

    const json = await res.json();

    if (!res.ok) {
      return { error: json.error || "Something went wrong" };
    }

    return { data: json };
  } catch {
    return { error: "Network error" };
  }
}

export const api = {
  auth: {
    me: () => request<{ user: User }>("/api/auth/me"),
    logout: () =>
      request<{ success: boolean }>("/api/auth/logout", { method: "POST" }),
    googleUrl: `${API_URL}/api/auth/google`,
  },
  users: {
    onboard: (data: { username: string; name: string }) =>
      request<{ success: boolean }>("/api/users/onboard", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    updateProfile: (data: {
      name?: string;
      bio?: string;
      username?: string;
      theme?: string;
      customLinks?: { title: string; url: string }[];
    }) =>
      request<{ success: boolean }>("/api/users/profile", {
        method: "PATCH",
        body: JSON.stringify(data),
      }),
    uploadImage: async (file: File): Promise<ApiResponse<{ image: string }>> => {
      try {
        const formData = new FormData();
        formData.append("image", file);
        const res = await fetch(`${API_URL}/api/users/profile/image`, {
          method: "POST",
          credentials: "include",
          body: formData,
        });
        const json = await res.json();
        if (!res.ok) return { error: json.error || "Upload failed" };
        return { data: json };
      } catch {
        return { error: "Network error" };
      }
    },
    publicProfile: (username: string) =>
      request<{ user: PublicUser }>(`/api/users/${username}`),
    publicConnections: (username: string) =>
      request<{ connections: Connection[] }>(
        `/api/users/${username}/connections`
      ),
  },
  connections: {
    list: () => request<{ connections: Connection[] }>("/api/connections"),
    disconnect: (platform: string) =>
      request<{ success: boolean }>(`/api/connections/${platform}`, {
        method: "DELETE",
      }),
    refresh: (platform: string) =>
      request<{ status: string }>(`/api/connections/${platform}/refresh`, {
        method: "POST",
      }),
    connectUrl: (platform: string) =>
      `${API_URL}/api/connections/${platform}/connect`,
  },
  analytics: {
    summary: () =>
      request<{
        totalViews: number;
        viewsToday: number;
        viewsThisWeek: number;
        totalClicks: number;
        clicksByType: Record<string, number>;
        viewsByDay: { date: string; count: number }[];
      }>("/api/analytics"),
    trackView: (username: string) =>
      request<{ ok: boolean }>(`/api/users/${username}/view`, {
        method: "POST",
      }),
    trackClick: (username: string, type: string, value: string) =>
      request<{ ok: boolean }>(`/api/users/${username}/click`, {
        method: "POST",
        body: JSON.stringify({ type, value }),
      }),
  },
  admin: {
    stats: () =>
      request<{
        totalUsers: number;
        onboardedUsers: number;
        verifiedUsers: number;
        totalConnections: number;
        totalViews: number;
        totalClicks: number;
      }>("/api/admin/stats"),
    users: (limit = 50, offset = 0) =>
      request<{
        users: {
          id: string;
          name: string | null;
          email: string;
          username: string | null;
          image: string | null;
          role: string;
          creatorScore: number | null;
          isVerified: boolean;
          onboarded: boolean;
          connections: number;
        }[];
        total: number;
      }>(`/api/admin/users?limit=${limit}&offset=${offset}`),
    setVerified: (userId: string, verified: boolean) =>
      request<{ success: boolean }>("/api/admin/users/verify", {
        method: "POST",
        body: JSON.stringify({ userId, verified }),
      }),
    sendDigest: () =>
      request<{ sent: number }>("/api/admin/digest", { method: "POST" }),
  },
  discover: {
    list: (params?: { limit?: number; offset?: number; minScore?: number; platform?: string }) => {
      const p = new URLSearchParams();
      if (params?.limit) p.set("limit", String(params.limit));
      if (params?.offset) p.set("offset", String(params.offset));
      if (params?.minScore) p.set("minScore", String(params.minScore));
      if (params?.platform) p.set("platform", params.platform);
      return request<{
        creators: {
          id: string;
          name: string | null;
          username: string | null;
          image: string | null;
          bio: string | null;
          creatorScore: number | null;
          isVerified: boolean;
          connections: number;
        }[];
        total: number;
      }>(`/api/discover?${p.toString()}`);
    },
  },
  collaborations: {
    send: (toUserId: string, message: string) =>
      request<{ id: string }>("/api/collaborations", {
        method: "POST",
        body: JSON.stringify({ toUserId, message }),
      }),
    inbox: () =>
      request<{
        requests: {
          id: string;
          fromUserId: string;
          toUserId: string;
          message: string;
          status: string;
          createdAt: string;
          fromName: string | null;
          fromUsername: string | null;
          fromImage: string | null;
        }[];
      }>("/api/collaborations/inbox"),
    outbox: () =>
      request<{
        requests: {
          id: string;
          fromUserId: string;
          toUserId: string;
          message: string;
          status: string;
          createdAt: string;
          toName: string | null;
          toUsername: string | null;
          toImage: string | null;
        }[];
      }>("/api/collaborations/outbox"),
    respond: (id: string, action: "accept" | "decline") =>
      request<{ status: string }>(`/api/collaborations/${id}/respond`, {
        method: "POST",
        body: JSON.stringify({ action }),
      }),
  },
};
