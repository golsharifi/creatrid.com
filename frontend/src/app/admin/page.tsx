"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { Users, BarChart3, Eye, MousePointerClick, Link2, CheckCircle, Shield, ClipboardList } from "lucide-react";
import { useTranslation } from "react-i18next";

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
  const { t } = useTranslation();
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [activeTab, setActiveTab] = useState<"users" | "audit">("users");
  const [auditEntries, setAuditEntries] = useState<any[]>([]);
  const [auditTotal, setAuditTotal] = useState(0);
  const [auditPage, setAuditPage] = useState(0);

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

  useEffect(() => {
    if (user?.role === "ADMIN" && activeTab === "audit") {
      api.adminAudit.list(50, auditPage * 50).then((r) => {
        if (r.data) {
          setAuditEntries(r.data.entries || []);
          setAuditTotal(r.data.total);
        }
      });
    }
  }, [user, activeTab, auditPage]);

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
        <h1 className="text-2xl font-bold">{t("admin.title")}</h1>
        <p className="mt-1 text-zinc-500">{t("admin.subtitle")}</p>
      </div>

      {/* Stats Grid */}
      {stats && (
        <div className="mb-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <StatCard icon={Users} label={t("admin.totalUsers")} value={stats.totalUsers} sub={t("admin.onboarded", { count: stats.onboardedUsers })} />
          <StatCard icon={Shield} label={t("admin.verifiedUsers")} value={stats.verifiedUsers} />
          <StatCard icon={Link2} label={t("admin.totalConnections")} value={stats.totalConnections} />
          <StatCard icon={Eye} label={t("admin.profileViews")} value={stats.totalViews} />
          <StatCard icon={MousePointerClick} label={t("admin.linkClicks")} value={stats.totalClicks} />
          <StatCard
            icon={BarChart3}
            label={t("admin.avgScore")}
            value={stats.totalUsers > 0 ? "\u2014" : "\u2014"}
            sub={t("admin.acrossAllUsers")}
          />
        </div>
      )}

      {/* Tabs */}
      <div className="mb-6 flex gap-4 border-b border-zinc-200 dark:border-zinc-800">
        <button
          onClick={() => setActiveTab("users")}
          className={`border-b-2 px-4 py-2 text-sm font-medium ${activeTab === "users" ? "border-zinc-900 text-zinc-900 dark:border-zinc-100 dark:text-zinc-100" : "border-transparent text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300"}`}
        >
          <Users className="mr-1.5 inline h-4 w-4" />
          {t("admin.tableUser")}s
        </button>
        <button
          onClick={() => setActiveTab("audit")}
          className={`border-b-2 px-4 py-2 text-sm font-medium ${activeTab === "audit" ? "border-zinc-900 text-zinc-900 dark:border-zinc-100 dark:text-zinc-100" : "border-transparent text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300"}`}
        >
          <ClipboardList className="mr-1.5 inline h-4 w-4" />
          {t("admin.auditLog")}
        </button>
      </div>

      {/* Audit Log */}
      {activeTab === "audit" && (
        <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
          <div className="border-b border-zinc-200 px-6 py-4 dark:border-zinc-800">
            <h2 className="font-semibold">{t("admin.auditLog")} ({auditTotal})</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
                <tr>
                  <th className="px-6 py-3 font-medium">{t("admin.auditTime")}</th>
                  <th className="px-6 py-3 font-medium">{t("admin.auditAdmin")}</th>
                  <th className="px-6 py-3 font-medium">{t("admin.auditAction")}</th>
                  <th className="px-6 py-3 font-medium">{t("admin.auditTarget")}</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                {auditEntries.map((entry: any) => (
                  <tr key={entry.id}>
                    <td className="px-6 py-3 text-xs text-zinc-500">{new Date(entry.createdAt).toLocaleString()}</td>
                    <td className="px-6 py-3">{entry.adminName || entry.adminUserId}</td>
                    <td className="px-6 py-3">
                      <span className="rounded-full bg-zinc-100 px-2 py-0.5 text-xs font-medium dark:bg-zinc-800">{entry.action}</span>
                    </td>
                    <td className="px-6 py-3 text-xs text-zinc-500">{entry.targetType}{entry.targetId ? `: ${entry.targetId}` : ""}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {auditTotal > 50 && (
            <div className="flex items-center justify-between border-t border-zinc-200 px-6 py-3 dark:border-zinc-800">
              <button onClick={() => setAuditPage((p) => Math.max(0, p - 1))} disabled={auditPage === 0} className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium disabled:opacity-50 dark:border-zinc-700">{t("common.previous")}</button>
              <span className="text-xs text-zinc-500">{t("common.page", { current: auditPage + 1, total: Math.ceil(auditTotal / 50) })}</span>
              <button onClick={() => setAuditPage((p) => p + 1)} disabled={(auditPage + 1) * 50 >= auditTotal} className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium disabled:opacity-50 dark:border-zinc-700">{t("common.next")}</button>
            </div>
          )}
        </div>
      )}

      {/* Users Table */}
      {activeTab === "users" && (
      <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
        <div className="border-b border-zinc-200 px-6 py-4 dark:border-zinc-800">
          <h2 className="font-semibold">{t("admin.usersCount", { count: total })}</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
              <tr>
                <th className="px-6 py-3 font-medium">{t("admin.tableUser")}</th>
                <th className="px-6 py-3 font-medium">{t("admin.tableUsername")}</th>
                <th className="px-6 py-3 font-medium">{t("admin.tableScore")}</th>
                <th className="px-6 py-3 font-medium">{t("admin.tableConnections")}</th>
                <th className="px-6 py-3 font-medium">{t("admin.tableStatus")}</th>
                <th className="px-6 py-3 font-medium">{t("admin.tableActions")}</th>
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
                        <p className="font-medium">{u.name || t("common.noData")}</p>
                        <p className="text-xs text-zinc-500">{u.email}</p>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-3 text-zinc-500">
                    {u.username ? `@${u.username}` : t("common.noData")}
                  </td>
                  <td className="px-6 py-3">
                    {u.creatorScore !== null ? u.creatorScore : t("common.noData")}
                  </td>
                  <td className="px-6 py-3">{u.connections}</td>
                  <td className="px-6 py-3">
                    <div className="flex items-center gap-2">
                      {u.isVerified && (
                        <span className="inline-flex items-center gap-1 rounded-full bg-blue-50 px-2 py-0.5 text-xs font-medium text-blue-600 dark:bg-blue-900/30 dark:text-blue-400">
                          <CheckCircle className="h-3 w-3" />
                          {t("admin.verified")}
                        </span>
                      )}
                      {!u.onboarded && (
                        <span className="rounded-full bg-zinc-100 px-2 py-0.5 text-xs text-zinc-500 dark:bg-zinc-800">
                          {t("admin.notOnboarded")}
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="px-6 py-3">
                    <button
                      onClick={() => toggleVerified(u.id, u.isVerified)}
                      className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                    >
                      {u.isVerified ? t("admin.unverify") : t("admin.verify")}
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
              {t("common.previous")}
            </button>
            <span className="text-xs text-zinc-500">
              {t("common.page", { current: page + 1, total: Math.ceil(total / 50) })}
            </span>
            <button
              onClick={() => setPage((p) => p + 1)}
              disabled={(page + 1) * 50 >= total}
              className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium disabled:opacity-50 dark:border-zinc-700"
            >
              {t("common.next")}
            </button>
          </div>
        )}
      </div>
      )}
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
