"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState, Suspense } from "react";
import { api } from "@/lib/api";
import { DollarSign, CreditCard, Clock, ExternalLink } from "@/components/icons";
import { useTranslation } from "react-i18next";

type Payout = {
  id: string;
  amountCents: number;
  currency: string;
  status: string;
  createdAt: string;
  completedAt: string | null;
};

function EarningsContent() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [connectStatus, setConnectStatus] = useState<{ connected: boolean; onboarded: boolean } | null>(null);
  const [dashboard, setDashboard] = useState<{ totalEarnedCents: number; totalPaidCents: number; pendingCents: number } | null>(null);
  const [payouts, setPayouts] = useState<Payout[]>([]);
  const [loadingData, setLoadingData] = useState(true);
  const [connecting, setConnecting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  useEffect(() => {
    if (!user) return;
    Promise.all([
      api.payouts.connectStatus(),
      api.payouts.dashboard(),
      api.payouts.list(),
    ]).then(([statusRes, dashRes, payoutsRes]) => {
      if (statusRes.data) setConnectStatus(statusRes.data as any);
      if (dashRes.data) setDashboard(dashRes.data as any);
      if (payoutsRes.data) setPayouts((payoutsRes.data as any).payouts ?? []);
      setLoadingData(false);
    }).catch(() => {
      setLoadingData(false);
    });
  }, [user]);

  const handleConnect = async () => {
    setConnecting(true);
    setError(null);
    const result = await api.payouts.connect();
    if (result.data && (result.data as any).url) {
      window.location.href = (result.data as any).url;
    } else {
      const errMsg = (result as any).error?.error || (result as any).error || "Failed to start Stripe Connect onboarding";
      setError(typeof errMsg === "string" ? errMsg : JSON.stringify(errMsg));
      setConnecting(false);
    }
  };

  if (loading || !user) return null;

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">{t("earnings.title")}</h1>
        <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{t("earnings.subtitle")}</p>
      </div>

      {loadingData ? (
        <div className="flex items-center justify-center py-16">
          <div className="h-6 w-6 animate-spin rounded-full border-2 border-zinc-300 border-t-zinc-900 dark:border-zinc-600 dark:border-t-zinc-100" />
        </div>
      ) : (
        <>
          {/* Stripe Connect card */}
          {!connectStatus?.onboarded && (
            <div className="mb-8 rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <div className="flex items-start gap-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-purple-100 text-purple-600 dark:bg-purple-900/50 dark:text-purple-400">
                  <CreditCard className="h-6 w-6" />
                </div>
                <div className="flex-1">
                  <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">{t("earnings.connectStripe")}</h2>
                  <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{t("earnings.connectDesc")}</p>
                  <button
                    onClick={handleConnect}
                    disabled={connecting}
                    className="mt-4 inline-flex items-center gap-2 rounded-lg bg-purple-600 px-4 py-2 text-sm font-medium text-white hover:bg-purple-700 disabled:opacity-50"
                  >
                    <ExternalLink className="h-4 w-4" />
                    {connecting ? t("common.loading") : t("earnings.connectButton")}
                  </button>
                  {error && (
                    <p className="mt-2 text-sm text-red-600 dark:text-red-400">{error}</p>
                  )}
                </div>
              </div>
            </div>
          )}

          {/* Revenue summary */}
          {dashboard && (
            <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-3">
              <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-950">
                <div className="flex items-center gap-3">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-green-100 text-green-600 dark:bg-green-900/50 dark:text-green-400">
                    <DollarSign className="h-5 w-5" />
                  </div>
                  <div>
                    <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">${(dashboard.totalEarnedCents / 100).toFixed(2)}</p>
                    <p className="text-xs text-zinc-500 dark:text-zinc-400">{t("earnings.totalEarned")}</p>
                  </div>
                </div>
              </div>
              <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-950">
                <div className="flex items-center gap-3">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100 text-blue-600 dark:bg-blue-900/50 dark:text-blue-400">
                    <CreditCard className="h-5 w-5" />
                  </div>
                  <div>
                    <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">${(dashboard.totalPaidCents / 100).toFixed(2)}</p>
                    <p className="text-xs text-zinc-500 dark:text-zinc-400">{t("earnings.totalPaid")}</p>
                  </div>
                </div>
              </div>
              <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-950">
                <div className="flex items-center gap-3">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-amber-100 text-amber-600 dark:bg-amber-900/50 dark:text-amber-400">
                    <Clock className="h-5 w-5" />
                  </div>
                  <div>
                    <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">${(dashboard.pendingCents / 100).toFixed(2)}</p>
                    <p className="text-xs text-zinc-500 dark:text-zinc-400">{t("earnings.pending")}</p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Payout history */}
          <h2 className="mb-4 text-lg font-semibold text-zinc-900 dark:text-zinc-100">{t("earnings.payoutHistory")}</h2>
          {payouts.length === 0 ? (
            <div className="rounded-xl border border-zinc-200 bg-zinc-50 py-12 text-center dark:border-zinc-800 dark:bg-zinc-900">
              <p className="text-sm text-zinc-500 dark:text-zinc-400">{t("earnings.noPayouts")}</p>
            </div>
          ) : (
            <div className="overflow-hidden rounded-xl border border-zinc-200 dark:border-zinc-800">
              <table className="w-full text-left text-sm">
                <thead className="border-b border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900">
                  <tr>
                    <th className="px-4 py-3 font-medium text-zinc-600 dark:text-zinc-400">{t("earnings.amount")}</th>
                    <th className="px-4 py-3 font-medium text-zinc-600 dark:text-zinc-400">{t("earnings.status")}</th>
                    <th className="px-4 py-3 font-medium text-zinc-600 dark:text-zinc-400">{t("earnings.date")}</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
                  {payouts.map((p) => (
                    <tr key={p.id} className="bg-white dark:bg-zinc-950">
                      <td className="px-4 py-3 font-medium text-zinc-900 dark:text-zinc-100">${(p.amountCents / 100).toFixed(2)}</td>
                      <td className="px-4 py-3">
                        <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${p.status === "completed" ? "bg-green-100 text-green-700 dark:bg-green-900/50 dark:text-green-400" : "bg-amber-100 text-amber-700 dark:bg-amber-900/50 dark:text-amber-400"}`}>
                          {p.status}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-zinc-500 dark:text-zinc-400">{new Date(p.createdAt).toLocaleDateString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}
    </div>
  );
}

export default function EarningsPage() {
  return (
    <Suspense fallback={null}>
      <EarningsContent />
    </Suspense>
  );
}
