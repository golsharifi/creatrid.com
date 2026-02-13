"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { Users, BarChart3, Eye, Settings, X } from "@/components/icons";
import { useTranslation } from "react-i18next";

type AgencyData = {
  id: string;
  name: string;
  website: string | null;
  description: string | null;
  isVerified: boolean;
  maxCreators: number;
};

type AgencyStats = {
  totalCreators: number;
  totalViews: number;
  avgScore: number;
  apiCalls30d: number;
};

type Creator = {
  id: string;
  agencyId: string;
  creatorId: string;
  status: string;
  invitedAt: string;
  joinedAt: string | null;
  creatorName: string | null;
  creatorUsername: string | null;
  creatorImage: string | null;
  creatorScore: number | null;
};

export default function AgencyPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<"dashboard" | "creators" | "apiUsage" | "settings">("dashboard");
  const [agency, setAgency] = useState<AgencyData | null>(null);
  const [stats, setStats] = useState<AgencyStats | null>(null);
  const [creators, setCreators] = useState<Creator[]>([]);
  const [creatorsTotal, setCreatorsTotal] = useState(0);
  const [creatorsPage, setCreatorsPage] = useState(0);
  const [inviteUsername, setInviteUsername] = useState("");
  const [inviting, setInviting] = useState(false);
  const [inviteError, setInviteError] = useState("");
  const [inviteSuccess, setInviteSuccess] = useState("");
  const [apiUsage, setApiUsage] = useState<any>(null);
  const [hasAgency, setHasAgency] = useState<boolean | null>(null);

  // Create agency form
  const [createName, setCreateName] = useState("");
  const [createWebsite, setCreateWebsite] = useState("");
  const [createDescription, setCreateDescription] = useState("");
  const [creating, setCreating] = useState(false);

  // Settings form
  const [settingsName, setSettingsName] = useState("");
  const [settingsWebsite, setSettingsWebsite] = useState("");
  const [settingsDescription, setSettingsDescription] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  // Load agency data
  useEffect(() => {
    if (!user) return;
    api.agency.get().then((r) => {
      if (r.data) {
        setAgency(r.data.agency);
        setStats(r.data.stats);
        setHasAgency(true);
        setSettingsName(r.data.agency.name);
        setSettingsWebsite(r.data.agency.website || "");
        setSettingsDescription(r.data.agency.description || "");
      } else {
        setHasAgency(false);
      }
    });
  }, [user]);

  // Load creators when tab changes
  useEffect(() => {
    if (activeTab === "creators" && agency) {
      api.agency.creators(20, creatorsPage * 20).then((r) => {
        if (r.data) {
          setCreators(r.data.creators);
          setCreatorsTotal(r.data.total);
        }
      });
    }
  }, [activeTab, agency, creatorsPage]);

  // Load API usage when tab changes
  useEffect(() => {
    if (activeTab === "apiUsage" && agency) {
      api.agency.apiUsage(30).then((r) => {
        if (r.data) setApiUsage(r.data);
      });
    }
  }, [activeTab, agency]);

  async function handleCreate() {
    if (!createName.trim()) return;
    setCreating(true);
    const result = await api.agency.create(createName, createWebsite || undefined, createDescription || undefined);
    if (result.data) {
      setAgency(result.data.agency);
      setHasAgency(true);
      setSettingsName(result.data.agency.name);
      setSettingsWebsite(result.data.agency.website || "");
      setSettingsDescription(result.data.agency.description || "");
      // Reload stats
      const statsResult = await api.agency.get();
      if (statsResult.data) setStats(statsResult.data.stats);
    }
    setCreating(false);
  }

  async function handleInvite() {
    if (!inviteUsername.trim()) return;
    setInviting(true);
    setInviteError("");
    setInviteSuccess("");
    const result = await api.agency.invite(inviteUsername.trim());
    if (result.data) {
      setInviteSuccess(`Invite sent to @${inviteUsername}`);
      setInviteUsername("");
    } else {
      setInviteError(result.error || "Failed to invite creator");
    }
    setInviting(false);
  }

  async function handleRemoveCreator(creatorId: string) {
    if (!confirm(t("agency.removeConfirm"))) return;
    const result = await api.agency.removeCreator(creatorId);
    if (result.data) {
      setCreators((prev) => prev.filter((c) => c.creatorId !== creatorId));
      setCreatorsTotal((prev) => prev - 1);
    }
  }

  async function handleSaveSettings() {
    setSaving(true);
    const result = await api.agency.update({
      name: settingsName,
      website: settingsWebsite || undefined,
      description: settingsDescription || undefined,
    });
    if (result.data) {
      setAgency((prev) =>
        prev
          ? { ...prev, name: settingsName, website: settingsWebsite || null, description: settingsDescription || null }
          : prev
      );
    }
    setSaving(false);
  }

  if (loading || !user) return null;

  // Show create agency form if no agency
  if (hasAgency === false) {
    return (
      <div className="mx-auto max-w-xl px-4 py-12 sm:px-6">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold">{t("agency.createTitle")}</h1>
          <p className="mt-1 text-zinc-500">{t("agency.createSubtitle")}</p>
        </div>
        <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
          <div className="space-y-4">
            <div>
              <label className="mb-1 block text-sm font-medium">{t("agency.agencyName")}</label>
              <input
                type="text"
                value={createName}
                onChange={(e) => setCreateName(e.target.value)}
                className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
                placeholder="Acme Talent Agency"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium">{t("agency.website")}</label>
              <input
                type="text"
                value={createWebsite}
                onChange={(e) => setCreateWebsite(e.target.value)}
                className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
                placeholder="https://example.com"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium">{t("agency.description")}</label>
              <textarea
                value={createDescription}
                onChange={(e) => setCreateDescription(e.target.value)}
                className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
                rows={3}
                placeholder="Tell us about your agency..."
              />
            </div>
            <button
              onClick={handleCreate}
              disabled={creating || !createName.trim()}
              className="w-full rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {creating ? t("common.loading") : t("agency.createButton")}
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (hasAgency === null) return null;

  const tabs = [
    { key: "dashboard" as const, label: t("agency.dashboard") },
    { key: "creators" as const, label: t("agency.creators") },
    { key: "apiUsage" as const, label: t("agency.apiUsage") },
    { key: "settings" as const, label: t("agency.settings") },
  ];

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">{t("agency.title")}</h1>
        <p className="mt-1 text-zinc-500">{t("agency.subtitle")}</p>
      </div>

      {/* Tabs */}
      <div className="mb-6 flex gap-4 border-b border-zinc-200 dark:border-zinc-800">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`border-b-2 px-4 py-2 text-sm font-medium ${
              activeTab === tab.key
                ? "border-zinc-900 text-zinc-900 dark:border-zinc-100 dark:text-zinc-100"
                : "border-transparent text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300"
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Dashboard Tab */}
      {activeTab === "dashboard" && stats && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatCard icon={Users} label={t("agency.managedCreators")} value={stats.totalCreators} />
          <StatCard icon={Eye} label={t("agency.totalViews")} value={stats.totalViews} />
          <StatCard icon={BarChart3} label={t("agency.avgScore")} value={Math.round(stats.avgScore)} />
          <StatCard icon={Settings} label={t("agency.apiCalls")} value={stats.apiCalls30d} />
        </div>
      )}

      {/* Creators Tab */}
      {activeTab === "creators" && (
        <div className="space-y-6">
          {/* Invite form */}
          <div className="flex items-end gap-3">
            <div className="flex-1">
              <label className="mb-1 block text-sm font-medium">{t("agency.inviteCreator")}</label>
              <input
                type="text"
                value={inviteUsername}
                onChange={(e) => setInviteUsername(e.target.value)}
                placeholder={t("agency.invitePlaceholder")}
                className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
                onKeyDown={(e) => {
                  if (e.key === "Enter") handleInvite();
                }}
              />
            </div>
            <button
              onClick={handleInvite}
              disabled={inviting || !inviteUsername.trim()}
              className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {inviting ? "..." : t("agency.invite")}
            </button>
          </div>
          {inviteError && (
            <div className="flex items-center gap-2 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400">
              <X className="h-4 w-4 shrink-0" />
              {inviteError}
            </div>
          )}
          {inviteSuccess && (
            <div className="rounded-lg bg-green-50 px-3 py-2 text-sm text-green-600 dark:bg-green-900/20 dark:text-green-400">
              {inviteSuccess}
            </div>
          )}

          {/* Creators table */}
          {creators.length === 0 ? (
            <div className="rounded-xl border border-zinc-200 px-6 py-12 text-center dark:border-zinc-800">
              <p className="text-zinc-500">{t("agency.noCreators")}</p>
            </div>
          ) : (
            <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
                    <tr>
                      <th className="px-6 py-3 font-medium">{t("common.creator")}</th>
                      <th className="px-6 py-3 font-medium">{t("agency.score")}</th>
                      <th className="px-6 py-3 font-medium">{t("admin.tableActions")}</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                    {creators.map((c) => (
                      <tr key={c.id}>
                        <td className="px-6 py-3">
                          <div className="flex items-center gap-3">
                            {c.creatorImage ? (
                              <img src={c.creatorImage} alt="" className="h-8 w-8 rounded-full object-cover" />
                            ) : (
                              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-zinc-100 text-xs font-bold dark:bg-zinc-800">
                                {(c.creatorName || "?")[0].toUpperCase()}
                              </div>
                            )}
                            <div>
                              <p className="font-medium">{c.creatorName || t("common.noData")}</p>
                              <p className="text-xs text-zinc-500">
                                {c.creatorUsername ? `@${c.creatorUsername}` : t("common.noData")}
                              </p>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-3">
                          {c.creatorScore !== null ? c.creatorScore : t("common.noData")}
                        </td>
                        <td className="px-6 py-3">
                          <button
                            onClick={() => handleRemoveCreator(c.creatorId)}
                            className="rounded-lg border border-red-200 px-3 py-1 text-xs font-medium text-red-600 transition-colors hover:bg-red-50 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-900/30"
                          >
                            {t("agency.removeCreator")}
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              {creatorsTotal > 20 && (
                <div className="flex items-center justify-between border-t border-zinc-200 px-6 py-3 dark:border-zinc-800">
                  <button
                    onClick={() => setCreatorsPage((p) => Math.max(0, p - 1))}
                    disabled={creatorsPage === 0}
                    className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium disabled:opacity-50 dark:border-zinc-700"
                  >
                    {t("common.previous")}
                  </button>
                  <span className="text-xs text-zinc-500">
                    {t("common.page", { current: creatorsPage + 1, total: Math.ceil(creatorsTotal / 20) })}
                  </span>
                  <button
                    onClick={() => setCreatorsPage((p) => p + 1)}
                    disabled={(creatorsPage + 1) * 20 >= creatorsTotal}
                    className="rounded-lg border border-zinc-200 px-3 py-1 text-xs font-medium disabled:opacity-50 dark:border-zinc-700"
                  >
                    {t("common.next")}
                  </button>
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* API Usage Tab */}
      {activeTab === "apiUsage" && (
        <div className="space-y-6">
          {apiUsage ? (
            <>
              {/* Total calls card */}
              <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
                <p className="text-sm text-zinc-500">{t("agency.usageThisMonth")}</p>
                <p className="mt-1 text-3xl font-bold">{(apiUsage.totalCalls || 0).toLocaleString()}</p>
              </div>

              {/* Usage by day chart */}
              {apiUsage.byDay && apiUsage.byDay.length > 0 && (
                <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
                  <h3 className="mb-4 font-semibold">{t("agency.apiUsage")} (30d)</h3>
                  <div className="flex items-end gap-1" style={{ height: 120 }}>
                    {apiUsage.byDay.map((d: any, i: number) => {
                      const maxCount = Math.max(...apiUsage.byDay.map((x: any) => x.count), 1);
                      const height = (d.count / maxCount) * 100;
                      return (
                        <div
                          key={i}
                          className="flex-1 rounded-t bg-zinc-900 dark:bg-zinc-100"
                          style={{ height: `${Math.max(height, 2)}%` }}
                          title={`${d.date}: ${d.count}`}
                        />
                      );
                    })}
                  </div>
                </div>
              )}

              {/* Usage by endpoint */}
              {apiUsage.byEndpoint && apiUsage.byEndpoint.length > 0 ? (
                <div className="rounded-xl border border-zinc-200 dark:border-zinc-800">
                  <div className="border-b border-zinc-200 px-6 py-4 dark:border-zinc-800">
                    <h3 className="font-semibold">{t("agency.usageByEndpoint")}</h3>
                  </div>
                  <div className="overflow-x-auto">
                    <table className="w-full text-left text-sm">
                      <thead className="border-b border-zinc-200 text-xs text-zinc-500 dark:border-zinc-800">
                        <tr>
                          <th className="px-6 py-3 font-medium">{t("agency.endpoint")}</th>
                          <th className="px-6 py-3 font-medium">{t("agency.calls")}</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                        {apiUsage.byEndpoint.map((ep: any, i: number) => (
                          <tr key={i}>
                            <td className="px-6 py-3 font-mono text-xs">{ep.endpoint}</td>
                            <td className="px-6 py-3">{ep.count.toLocaleString()}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              ) : (
                <div className="rounded-xl border border-zinc-200 px-6 py-12 text-center dark:border-zinc-800">
                  <p className="text-zinc-500">{t("agency.noUsage")}</p>
                </div>
              )}
            </>
          ) : (
            <div className="rounded-xl border border-zinc-200 px-6 py-12 text-center dark:border-zinc-800">
              <p className="text-zinc-500">{t("agency.noUsage")}</p>
            </div>
          )}
        </div>
      )}

      {/* Settings Tab */}
      {activeTab === "settings" && agency && (
        <div className="max-w-xl space-y-4">
          <div>
            <label className="mb-1 block text-sm font-medium">{t("agency.agencyName")}</label>
            <input
              type="text"
              value={settingsName}
              onChange={(e) => setSettingsName(e.target.value)}
              className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium">{t("agency.website")}</label>
            <input
              type="text"
              value={settingsWebsite}
              onChange={(e) => setSettingsWebsite(e.target.value)}
              className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
              placeholder="https://example.com"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium">{t("agency.description")}</label>
            <textarea
              value={settingsDescription}
              onChange={(e) => setSettingsDescription(e.target.value)}
              className="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
              rows={4}
            />
          </div>
          <button
            onClick={handleSaveSettings}
            disabled={saving}
            className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
          >
            {saving ? t("agency.saving") : t("agency.saveChanges")}
          </button>
        </div>
      )}
    </div>
  );
}

function StatCard({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: number | string;
}) {
  return (
    <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
      <div className="mb-2 flex h-9 w-9 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
        <Icon className="h-4 w-4" />
      </div>
      <p className="text-sm text-zinc-500">{label}</p>
      <p className="mt-1 text-2xl font-bold">{typeof value === "number" ? value.toLocaleString() : value}</p>
    </div>
  );
}
