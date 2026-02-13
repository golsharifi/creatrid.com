import type { User, PublicUser, Connection, EmailPrefs } from "./types";
import { captureException } from "./sentry";

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
  } catch (err) {
    captureException(err);
    return { error: "Network error" };
  }
}

export const api = {
  auth: {
    me: () => request<{ user: User }>("/api/auth/me"),
    logout: () =>
      request<{ success: boolean }>("/api/auth/logout", { method: "POST" }),
    googleUrl: `${API_URL}/api/auth/google`,
    sendVerificationEmail: () =>
      request<{ success: boolean }>("/api/auth/verify-email/send", {
        method: "POST",
      }),
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
      emailPrefs?: EmailPrefs;
    }) =>
      request<{ success: boolean }>("/api/users/profile", {
        method: "PATCH",
        body: JSON.stringify(data),
      }),
    deleteAccount: () =>
      request<{ deleted: boolean }>("/api/users/account", { method: "DELETE" }),
    exportProfile: () =>
      request<Record<string, unknown>>("/api/users/export"),
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
        clicksByDay: { date: string; count: number }[];
        topReferrers: { referrer: string; count: number }[];
        viewsByHour: { hour: number; count: number }[];
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
    list: (params?: { limit?: number; offset?: number; minScore?: number; platform?: string; q?: string }) => {
      const p = new URLSearchParams();
      if (params?.limit) p.set("limit", String(params.limit));
      if (params?.offset) p.set("offset", String(params.offset));
      if (params?.minScore) p.set("minScore", String(params.minScore));
      if (params?.platform) p.set("platform", params.platform);
      if (params?.q) p.set("q", params.q);
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
  apiKeys: {
    create: (name: string) =>
      request<{ id: string; name: string; key: string }>("/api/keys", {
        method: "POST",
        body: JSON.stringify({ name }),
      }),
    list: () =>
      request<{
        keys: {
          id: string;
          name: string;
          keyPrefix: string;
          lastUsedAt: string | null;
          createdAt: string;
        }[];
      }>("/api/keys"),
    delete: (id: string) =>
      request<{ success: boolean }>(`/api/keys/${id}`, { method: "DELETE" }),
  },
  billing: {
    checkout: (plan: string) =>
      request<{ url: string }>("/api/billing/checkout", {
        method: "POST",
        body: JSON.stringify({ plan }),
      }),
    subscription: () =>
      request<{
        plan: string;
        status: string;
        currentPeriodEnd: string | null;
      }>("/api/billing/subscription"),
    portal: () =>
      request<{ url: string }>("/api/billing/portal", { method: "POST" }),
  },
  content: {
    upload: async (file: File, title: string, description: string, tags: string[], isPublic: boolean): Promise<ApiResponse<{ item: any }>> => {
      try {
        const formData = new FormData();
        formData.append("file", file);
        formData.append("title", title);
        formData.append("description", description);
        formData.append("tags", tags.join(","));
        formData.append("is_public", String(isPublic));
        const res = await fetch(`${API_URL}/api/content`, {
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
    list: (limit = 20, offset = 0) =>
      request<{ items: any[]; total: number }>(`/api/content?limit=${limit}&offset=${offset}`),
    get: (id: string) => request<any>(`/api/content/${id}`),
    update: (id: string, data: { title?: string; description?: string; tags?: string[]; isPublic?: boolean }) =>
      request<{ success: boolean }>(`/api/content/${id}`, {
        method: "PATCH",
        body: JSON.stringify(data),
      }),
    delete: (id: string) =>
      request<{ success: boolean }>(`/api/content/${id}`, { method: "DELETE" }),
    download: (id: string) => `${API_URL}/api/content/${id}/download`,
    proof: (id: string) => request<{ id: string; hashSha256: string; createdAt: string; title: string }>(`/api/content/${id}/proof`),
    publicList: (username: string) =>
      request<{ items: any[] }>(`/api/users/${username}/content`),
  },
  licenses: {
    create: (contentId: string, licenseType: string, priceCents: number, termsText?: string) =>
      request<any>(`/api/content/${contentId}/licenses`, {
        method: "POST",
        body: JSON.stringify({ licenseType, priceCents, termsText }),
      }),
    list: (contentId: string) =>
      request<{ offerings: any[] }>(`/api/content/${contentId}/licenses`),
    update: (id: string, data: { priceCents?: number; isActive?: boolean; termsText?: string }) =>
      request<{ success: boolean }>(`/api/licenses/${id}`, {
        method: "PATCH",
        body: JSON.stringify(data),
      }),
    delete: (id: string) =>
      request<{ success: boolean }>(`/api/licenses/${id}`, { method: "DELETE" }),
    checkout: (id: string) =>
      request<{ url: string }>(`/api/licenses/${id}/checkout`, { method: "POST" }),
    purchases: () => request<{ purchases: any[] }>("/api/licenses/purchases"),
    sales: () => request<{ sales: any[] }>("/api/licenses/sales"),
  },
  marketplace: {
    browse: (params?: { type?: string; q?: string; sort?: string; limit?: number; offset?: number }) => {
      const p = new URLSearchParams();
      if (params?.type) p.set("type", params.type);
      if (params?.q) p.set("q", params.q);
      if (params?.sort) p.set("sort", params.sort);
      if (params?.limit) p.set("limit", String(params.limit));
      if (params?.offset) p.set("offset", String(params.offset));
      return request<{ items: any[]; total: number }>(`/api/marketplace?${p.toString()}`);
    },
    detail: (id: string) => request<any>(`/api/marketplace/${id}`),
  },
  dmca: {
    report: (contentId: string, data: { reporterEmail: string; reporterName: string; reason: string; evidenceUrl?: string }) =>
      request<{ id: string }>(`/api/content/${contentId}/report`, {
        method: "POST",
        body: JSON.stringify(data),
      }),
  },
  notifications: {
    list: (limit = 20, offset = 0) =>
      request<{ notifications: any[]; total: number }>(`/api/notifications?limit=${limit}&offset=${offset}`),
    unreadCount: () =>
      request<{ count: number }>("/api/notifications/unread-count"),
    markRead: (id: string) =>
      request<{ success: boolean }>(`/api/notifications/${id}/read`, { method: "POST" }),
    markAllRead: () =>
      request<{ success: boolean }>("/api/notifications/read-all", { method: "POST" }),
  },
  contentAnalytics: {
    trackView: (contentId: string) =>
      request<{ ok: boolean }>(`/api/content/${contentId}/view`, { method: "POST" }),
    item: (contentId: string) =>
      request<any>(`/api/content/${contentId}/analytics`),
    summary: () =>
      request<{ items: any[]; totalRevenue: number }>("/api/content-analytics"),
  },
  payouts: {
    connect: () =>
      request<{ url: string }>("/api/payouts/connect", { method: "POST" }),
    connectStatus: () =>
      request<{ connected: boolean; onboarded: boolean }>("/api/payouts/connect/status"),
    dashboard: () =>
      request<{ totalEarnedCents: number; totalPaidCents: number; pendingCents: number }>("/api/payouts/dashboard"),
    list: (limit = 20, offset = 0) =>
      request<{ payouts: any[]; total: number }>(`/api/payouts?limit=${limit}&offset=${offset}`),
  },
  collections: {
    create: (title: string, description?: string, isPublic = true) =>
      request<any>("/api/collections", {
        method: "POST",
        body: JSON.stringify({ title, description, isPublic }),
      }),
    list: (limit = 20, offset = 0) =>
      request<{ collections: any[]; total: number }>(`/api/collections?limit=${limit}&offset=${offset}`),
    get: (id: string) =>
      request<{ collection: any; items: any[] }>(`/api/collections/${id}`),
    update: (id: string, data: { title?: string; description?: string; isPublic?: boolean }) =>
      request<{ success: boolean }>(`/api/collections/${id}`, {
        method: "PATCH",
        body: JSON.stringify(data),
      }),
    delete: (id: string) =>
      request<{ success: boolean }>(`/api/collections/${id}`, { method: "DELETE" }),
    addItem: (collectionId: string, contentId: string, position = 0) =>
      request<{ success: boolean }>(`/api/collections/${collectionId}/items`, {
        method: "POST",
        body: JSON.stringify({ contentId, position }),
      }),
    removeItem: (collectionId: string, contentId: string) =>
      request<{ success: boolean }>(`/api/collections/${collectionId}/items/${contentId}`, {
        method: "DELETE",
      }),
    items: (collectionId: string) =>
      request<{ items: any[] }>(`/api/collections/${collectionId}/items`),
    publicList: (username: string, limit = 20, offset = 0) =>
      request<{ collections: any[]; total: number }>(`/api/users/${username}/collections?limit=${limit}&offset=${offset}`),
  },
  search: {
    query: (q: string, type?: string, limit = 20, offset = 0) => {
      const p = new URLSearchParams();
      p.set("q", q);
      if (type) p.set("type", type);
      p.set("limit", String(limit));
      p.set("offset", String(offset));
      return request<any>(`/api/search?${p.toString()}`);
    },
    suggestions: (q: string) =>
      request<{ suggestions: string[] }>(`/api/search/suggestions?q=${encodeURIComponent(q)}`),
  },
};
