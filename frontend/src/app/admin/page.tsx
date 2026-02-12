"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { Users, BarChart3, Eye, MousePointerClick, Link2, CheckCircle, Shield } from "lucide-react";

type AdminStats = {
  totalUsers: number;
  onboardedUsers: number;
  verifiedUsers: number;
  totalConnections: number;
  totalViews: number;
  totalClicks: number;
};

type AdminUser = {
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
};

export default function AdminPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
    if (!loading && user && user.role !== "ADMIN") router.push("/dashboard");
  }, [user, loading, router]);

  useEffect(() => {
    if (user?.role === "ADMIN") {
      api.admin.stats().then((r) => {
        if (r.data) setStats(r.data);
      });
    }
  }, [user]);

  useEffect(() => {
    if (user?.role === "ADMIN") {
      api.admin.users(50, page * 50).then((r) => {
        if (r.data) {
          setUsers(r.data.users);
          setTotal(r.data.total);
        }
      });
    }
  }, [user, page]);

  async function toggleVerified(userId: string, current: boolean) {
    const result = await api.admin.setVerified(userId, !current);
    if (result.data) {
      setUsers((prev) =>
        prev.map((u) => (u.id === userId ? { ...u, isVerified: !current } : u))
      );
    }
  }

  if (loading || !user || user.role !== "ADMIN") return null;

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">Admin Dashboard</h1>
        <p className="mt-1 text-zinc-500">Platform overview and user management.</p>
      </div>

      {/* Stats Grid */}
      {stats && (
        <div className="mb-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <StatCard icon={Users} label="Total Users" value={stats.totalUsers} sub={`${stats.onboardedUsers} onboarded`} />
          <StatCard icon={Shield} label="Verified Users" value={stats.verifiedUsers} />
          <StatCard icon={Link2} label="Total Connections" value={stats.totalConnections} />
          <StatCard icon={Eye} label="Profile Views" value={stats.totalViews} />
          <StatCard icon={MousePointerClick} label="Link Clicks" value={stats.totalClicks} />
          <StatCard
            icon={BarChart3}
            label="Avg Score"
            value={stats.totalUsers > 0 ? "—" : "—"}
            sub="Across all users"
          />
        </div>
      )}

      {/* Users Table */}
      <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
        <div className="border-b border-zinc-200 px-6 py-4 dark:border-zinc-800">
          <h2 className="font-semibold">Users ({total})</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
              <tr>
                <th className="px-6 py-3 font-medium">User</th>
                <th className="px-6 py-3 font-medium">Username</th>
                <th className="px-6 py-3 font-medium">Score</th>
                <th className="px-6 py-3 font-medium">Connections</th>
                <th className="px-6 py-3 font-medium">Status</th>
                <th className="px-6 py-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
              {users.map((u) => (
                <tr key={u.id}>
                  <td className="px-6 py-3">
                    <div className="flex items-center gap-3">
                      {u.image ? (
                        <img src={u.image} alt="" className="h-8 w-8 rounded-full object-cover" />
                      ) : (
                        <div className="flex h-8 w-8 items-center justify-center rounded-full bg-zinc-100 text-xs font-bold dark:bg-zinc-800">
                          {(u.name || u.email)[0].toUpperCase()}
                        </div>
                      )}
                      <div>
                        <p className="font-medium">{u.name || "—"}</p>
                        <p className="text-xs text-zinc-500">{u.email}</p>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-3 text-zinc-500">
                    {u.username ? `@${u.username}` : "—"}
                  </td>
                  <td className="px-6 py-3">
                    {u.creatorScore !== null ? u.creatorScore : "—"}
                  </td>
                  <td className="px-6 py-3">{u.connections}</td>
                  <td className="px-6 py-3">
                    <div className="flex items-center gap-2">
                      {u.isVerified && (
                        <span className="inline-flex items-center gap-1 rounded-full bg-blue-50 px-2 py-0.5 text-xs font-medium text-blue-600 dark:bg-blue-900/30 dark:text-blue-400">
                          <CheckCircle className="h-3 w-3" />
                          Verified
                        </span>
                      )}
                      {!u.onboarded && (
                        <span className="rounded-full bg-zinc-100 px-2 py-0.5 text-xs text-zinc-500 dark:bg-zinc-800">
                          Not onboarded
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="px-6 py-3">
                    <button
                      onClick={() => toggleVerified(u.id, u.isVerified)}
                      className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                    >
                      {u.isVerified ? "Unverify" : "Verify"}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {total > 50 && (
          <div className="flex items-center justify-between border-t border-zinc-200 px-6 py-3 dark:border-zinc-800">
            <button
              onClick={() => setPage((p) => Math.max(0, p - 1))}
              disabled={page === 0}
              className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium disabled:opacity-50 dark:border-zinc-700"
            >
              Previous
            </button>
            <span className="text-xs text-zinc-500">
              Page {page + 1} of {Math.ceil(total / 50)}
            </span>
            <button
              onClick={() => setPage((p) => p + 1)}
              disabled={(page + 1) * 50 >= total}
              className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium disabled:opacity-50 dark:border-zinc-700"
            >
              Next
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

function StatCard({
  icon: Icon,
  label,
  value,
  sub,
}: {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: number | string;
  sub?: string;
}) {
  return (
    <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
      <div className="mb-2 flex h-9 w-9 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
        <Icon className="h-4 w-4" />
      </div>
      <p className="text-sm text-zinc-500">{label}</p>
      <p className="mt-1 text-2xl font-bold">{typeof value === "number" ? value.toLocaleString() : value}</p>
      {sub && <p className="mt-0.5 text-xs text-zinc-400">{sub}</p>}
    </div>
  );
}
