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
    connectUrl: (platform: string) =>
      `${API_URL}/api/connections/${platform}/connect`,
  },
};
